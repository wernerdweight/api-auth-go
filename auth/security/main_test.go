package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v3/auth/cache"
	"github.com/wernerdweight/api-auth-go/v3/auth/config"
	"github.com/wernerdweight/api-auth-go/v3/auth/constants"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
	"github.com/wernerdweight/api-auth-go/v3/auth/fup"
)

// unknownApiClientProvider resolves no client, like the provider does for the traffic that the
// anonymous FUP limits are meant to stop
type unknownApiClientProvider struct{}

func (p unknownApiClientProvider) ProvideByIdAndSecret(_ string, _ string) (contract.ApiClientInterface, *contract.AuthError) {
	return nil, contract.NewAuthError(contract.ClientNotFound, nil)
}

func (p unknownApiClientProvider) ProvideByApiKey(_ string) (contract.ApiClientInterface, *contract.AuthError) {
	return nil, contract.NewAuthError(contract.ClientNotFound, nil)
}

func (p unknownApiClientProvider) Save(_ contract.ApiClientInterface) *contract.AuthError {
	return nil
}

// brokenApiClientProvider fails for a reason that is not the caller's fault, e.g. an unreachable database
type brokenApiClientProvider struct {
	unknownApiClientProvider
}

func (p brokenApiClientProvider) ProvideByApiKey(_ string) (contract.ApiClientInterface, *contract.AuthError) {
	return nil, contract.NewInternalError(contract.DatabaseError, nil)
}

// brokenFUPCacheDriver can't count requests, e.g. because the cache is unreachable. It overrides
// both increments - the embedded driver implements the optional contract.FUPTTLCacheDriverInterface
// too, and overriding only IncrementFUPEntry would leave the embedded one in use
type brokenFUPCacheDriver struct {
	*cache.MemoryCacheDriver
}

func (d brokenFUPCacheDriver) IncrementFUPEntry(_ string) (*contract.FUPCacheEntry, *contract.AuthError) {
	return nil, contract.NewInternalError(contract.CacheError, nil)
}

func (d brokenFUPCacheDriver) IncrementFUPEntryWithTTL(_ string, _ time.Duration) (*contract.FUPCacheEntry, *contract.AuthError) {
	return nil, contract.NewInternalError(contract.CacheError, nil)
}

// initAnonymousFUP configures the provider with a fresh cache, so that each test starts with
// empty counters
func initAnonymousFUP(scope contract.FUPScope) {
	memoryDriver := cache.NewMemoryCacheDriver()
	initAnonymousFUPWith(scope, unknownApiClientProvider{}, memoryDriver)
}

func initAnonymousFUPWith(
	scope contract.FUPScope,
	provider contract.ApiClientProviderInterface[contract.ApiClientInterface],
	driver contract.CacheDriverInterface,
) {
	useScopeAccessModel := true
	apiKeyMode := true
	prefix := "test:"
	driver.Init(prefix, time.Hour)
	config.ProviderInstance.Init(contract.Config{
		Client: contract.ClientConfig{
			Provider:            provider,
			UseScopeAccessModel: &useScopeAccessModel,
			AnonymousFUPScope:   &scope,
			// auth.Middleware defaults the checker, the config provider can't (see auth.Middleware)
			AnonymousFUPChecker: fup.IPFUPChecker{},
		},
		Mode:  &contract.ModesConfig{ApiKey: &apiKeyMode},
		Cache: &contract.CacheConfig{Driver: driver, Prefix: &prefix},
	})
}

func requestFrom(ip string, apiKey string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	// gin trusts all proxies unless the application configures otherwise, so this is what
	// c.ClientIP() (and with it the per-IP FUP key) resolves to
	c.Request.Header.Set("X-Forwarded-For", ip)
	if "" != apiKey {
		c.Request.Header.Set(constants.ApiKeyHeader, apiKey)
	}
	return c
}

func TestAuthenticate_AnonymousFUP(t *testing.T) {
	const limit = 3
	scope := contract.FUPScope{
		constants.FUPIPKey: map[string]any{string(constants.PeriodMinutely): limit},
	}
	tests := []struct {
		name   string
		apiKey string
		// the error the request is expected to fail with while below the limit
		wantCode contract.AuthErrorCode
	}{
		{name: "no credentials", wantCode: contract.NoCredentialsProvided},
		{name: "unknown api key", apiKey: "does-not-exist", wantCode: contract.ClientNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initAnonymousFUP(scope)

			for i := 1; i <= limit; i++ {
				err := Authenticate(requestFrom("10.0.0.1", tt.apiKey))
				if nil == err {
					t.Fatalf("request %d: Authenticate() = nil, want an authentication error", i)
				}
				if err.Code != tt.wantCode {
					t.Fatalf("request %d: Authenticate() code = %v, want %v", i, err.Code, tt.wantCode)
				}
				if err.Status != http.StatusUnauthorized {
					t.Fatalf("request %d: Authenticate() status = %d, want %d", i, err.Status, http.StatusUnauthorized)
				}
			}

			c := requestFrom("10.0.0.1", tt.apiKey)
			err := Authenticate(c)
			if nil == err {
				t.Fatal("request over the limit: Authenticate() = nil, want a FUP error")
			}
			if err.Code != contract.RequestLimitDepleted {
				t.Errorf("request over the limit: Authenticate() code = %v, want %v", err.Code, contract.RequestLimitDepleted)
			}
			if err.Status != http.StatusTooManyRequests {
				t.Errorf("request over the limit: Authenticate() status = %d, want %d", err.Status, http.StatusTooManyRequests)
			}
			if "" == c.Writer.Header().Get(constants.RetryAfterHeader) {
				t.Errorf("request over the limit: %s header is empty", constants.RetryAfterHeader)
			}
		})
	}
}

// TestAuthenticate_AnonymousFUP_PerIP makes sure one depleted source does not lock out the others
func TestAuthenticate_AnonymousFUP_PerIP(t *testing.T) {
	const limit = 2
	initAnonymousFUP(contract.FUPScope{
		constants.FUPIPKey: map[string]any{string(constants.PeriodMinutely): limit},
	})

	for i := 0; i <= limit; i++ {
		Authenticate(requestFrom("10.0.0.1", ""))
	}

	err := Authenticate(requestFrom("10.0.0.2", ""))
	if nil == err {
		t.Fatal("Authenticate() = nil, want an authentication error")
	}
	if err.Code != contract.NoCredentialsProvided {
		t.Errorf("Authenticate() code = %v, want %v", err.Code, contract.NoCredentialsProvided)
	}
}

// TestAuthenticate_AnonymousFUP_CacheError pins the fail-open behaviour: if the limits can't be
// checked, unauthenticated requests must keep getting their authentication error - neither a 500
// for all of them, nor a 429 that was never counted
func TestAuthenticate_AnonymousFUP_CacheError(t *testing.T) {
	initAnonymousFUPWith(
		contract.FUPScope{
			constants.FUPIPKey: map[string]any{string(constants.PeriodMinutely): 1},
		},
		unknownApiClientProvider{},
		brokenFUPCacheDriver{cache.NewMemoryCacheDriver()},
	)

	for i := 0; i < 5; i++ {
		err := Authenticate(requestFrom("10.0.0.1", ""))
		if nil == err {
			t.Fatalf("request %d: Authenticate() = nil, want an authentication error", i)
		}
		if err.Code != contract.NoCredentialsProvided {
			t.Fatalf("request %d: Authenticate() code = %v, want %v", i, err.Code, contract.NoCredentialsProvided)
		}
		if err.Status != http.StatusUnauthorized {
			t.Fatalf("request %d: Authenticate() status = %d, want %d", i, err.Status, http.StatusUnauthorized)
		}
	}
}

// TestAuthenticate_AnonymousFUP_InternalError makes sure requests that fail to authenticate for a
// reason on our side are not limited - an outage would otherwise become a 429 for everyone
func TestAuthenticate_AnonymousFUP_InternalError(t *testing.T) {
	const limit = 2
	initAnonymousFUPWith(
		contract.FUPScope{
			constants.FUPIPKey: map[string]any{string(constants.PeriodMinutely): limit},
		},
		brokenApiClientProvider{},
		cache.NewMemoryCacheDriver(),
	)

	for i := 0; i <= limit+1; i++ {
		err := Authenticate(requestFrom("10.0.0.1", "any-api-key"))
		if nil == err {
			t.Fatalf("request %d: Authenticate() = nil, want the provider error", i)
		}
		if err.Code != contract.DatabaseError {
			t.Fatalf("request %d: Authenticate() code = %v, want %v", i, err.Code, contract.DatabaseError)
		}
		if err.Status != http.StatusInternalServerError {
			t.Fatalf("request %d: Authenticate() status = %d, want %d", i, err.Status, http.StatusInternalServerError)
		}
	}
}

// TestAuthenticate_AnonymousFUP_NoLimit covers a configured scope that does not limit by IP -
// unauthenticated traffic must keep getting the plain authentication error
func TestAuthenticate_AnonymousFUP_NoLimit(t *testing.T) {
	initAnonymousFUP(contract.FUPScope{})

	for i := 0; i < 10; i++ {
		err := Authenticate(requestFrom("10.0.0.1", ""))
		if nil == err {
			t.Fatalf("request %d: Authenticate() = nil, want an authentication error", i)
		}
		if err.Code != contract.NoCredentialsProvided {
			t.Fatalf("request %d: Authenticate() code = %v, want %v", i, err.Code, contract.NoCredentialsProvided)
		}
	}
}

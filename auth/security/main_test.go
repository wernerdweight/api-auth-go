package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/cache"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
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

// initAnonymousFUP configures the provider with a fresh cache, so that each test starts with
// empty counters
func initAnonymousFUP(scope contract.FUPScope) {
	useScopeAccessModel := true
	apiKeyMode := true
	prefix := "test:"
	driver := cache.NewMemoryCacheDriver()
	driver.Init(prefix, time.Hour)
	config.ProviderInstance.Init(contract.Config{
		Client: contract.ClientConfig{
			Provider:            unknownApiClientProvider{},
			UseScopeAccessModel: &useScopeAccessModel,
			AnonymousFUPScope:   &scope,
		},
		Mode:  &contract.ModesConfig{ApiKey: &apiKeyMode},
		Cache: &contract.CacheConfig{Driver: driver, Prefix: &prefix},
	})
}

func requestFrom(ip string, apiKey string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
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

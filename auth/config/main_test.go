package config

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/wernerdweight/api-auth-go/auth/cache"
	"github.com/wernerdweight/api-auth-go/auth/checker"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"testing"
	"time"
)

type mockApiClientProvider struct{}

func (m mockApiClientProvider) ProvideByIdAndSecret(id string, secret string) (contract.ApiClientInterface, *contract.AuthError) {
	return nil, nil
}

func (m mockApiClientProvider) ProvideByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	return nil, nil
}

func (m mockApiClientProvider) Save(client contract.ApiClientInterface) *contract.AuthError {
	return nil
}

type mockApiUserProvider struct{}

func (m mockApiUserProvider) ProvideByLoginAndPassword(login string, password string) (contract.ApiUserInterface, *contract.AuthError) {
	return nil, nil
}

func (m mockApiUserProvider) ProvideByLogin(login string) (contract.ApiUserInterface, *contract.AuthError) {
	return nil, nil
}

func (m mockApiUserProvider) ProvideByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	return nil, nil
}

func (m mockApiUserProvider) ProvideByConfirmationToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	return nil, nil
}

func (m mockApiUserProvider) ProvideByResetToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	return nil, nil
}

func (m mockApiUserProvider) ProvideNew(login string, encryptedPassword string) contract.ApiUserInterface {
	return nil
}

func (m mockApiUserProvider) Save(user contract.ApiUserInterface) *contract.AuthError {
	return nil
}

type mockFUPChecker struct{}

func (m mockFUPChecker) Check(fup *contract.FUPScope, c *gin.Context, key string) contract.FUPScopeLimits {
	return contract.FUPScopeLimits{}
}

type TestSuite struct {
	suite.Suite
	provider *Provider
}

func (s *TestSuite) SetupTest() {
	s.provider = &Provider{
		config: contract.Config{
			Client: contract.ClientConfig{
				Provider:                      nil,
				UseScopeAccessModel:           &defaultClientUseScopeAccessModel,
				AccessScopeChecker:            checker.PathAccessScopeChecker{},
				FUPChecker:                    nil,
				OneOffTokenExpirationInterval: &defaultOneOffTokenExpirationInterval,
			},
			User: &contract.UserConfig{
				Provider:                            nil,
				TokenFactory:                        nil,
				ApiTokenExpirationInterval:          &defaultExpirationInterval,
				UseScopeAccessModel:                 &defaultUserUseScopeAccessModel,
				AccessScopeChecker:                  checker.PathAccessScopeChecker{},
				WithRegistration:                    &defaultWithRegistration,
				ConfirmationTokenExpirationInterval: &defaultConfirmationExpirationInterval,
			},
			Mode: &contract.ModesConfig{
				ApiKey:            &defaultApiKeyMode,
				AdditionalApiKeys: &defaultAdditionalApiKeys,
				ClientIdAndSecret: &defaultClientIdAndSecretMode,
				OneOffToken:       &defaultOneOffTokenMode,
			},
			TargetHandlers:         nil,
			ExcludeHandlers:        nil,
			ExcludeOptionsRequests: &defaultExcludeOptionsRequests,
			Cache: &contract.CacheConfig{
				Driver: nil,
				Prefix: &defaultCachePrefix,
				TTL:    &defaultCacheTTL,
			},
			TargetOneOffTokenHandlers: nil,
		},
	}
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestProvider_Init() {
	s.Nil(s.provider.config.TargetHandlers)

	handlers := []string{"*"}
	s.provider.Init(contract.Config{
		TargetHandlers: &handlers,
	})

	s.NotNil(s.provider.config.TargetHandlers)
}

func (s *TestSuite) TestProvider_IsApiKeyModeEnabled() {
	s.False(s.provider.IsApiKeyModeEnabled())
	enabled := true
	s.provider.Init(contract.Config{
		Mode: &contract.ModesConfig{
			ApiKey: &enabled,
		},
	})
	s.True(s.provider.IsApiKeyModeEnabled())
}

func (s *TestSuite) TestProvider_IsAdditionalApiKeysEnabled() {
	s.False(s.provider.IsAdditionalApiKeysEnabled())
	enabled := true
	s.provider.Init(contract.Config{
		Mode: &contract.ModesConfig{
			AdditionalApiKeys: &enabled,
		},
	})
	s.True(s.provider.IsAdditionalApiKeysEnabled())
}

func (s *TestSuite) TestProvider_IsClientIdAndSecretModeEnabled() {
	s.True(s.provider.IsClientIdAndSecretModeEnabled())
	disabled := false
	s.provider.Init(contract.Config{
		Mode: &contract.ModesConfig{
			ClientIdAndSecret: &disabled,
		},
	})
	s.False(s.provider.IsClientIdAndSecretModeEnabled())
}

func (s *TestSuite) TestProvider_GetTargetHandlers() {
	s.Nil(s.provider.GetTargetHandlers())
	handlers := []string{".*"}
	s.provider.Init(contract.Config{
		TargetHandlers: &handlers,
	})
	s.Equal(&handlers, s.provider.GetTargetHandlers())
}

func (s *TestSuite) TestProvider_GetExcludeHandlers() {
	s.Nil(s.provider.GetExcludeHandlers())
	handlers := []string{".*"}
	s.provider.Init(contract.Config{
		ExcludeHandlers: &handlers,
	})
	s.Equal(&handlers, s.provider.GetExcludeHandlers())
}

func (s *TestSuite) TestProvider_GetTargetOneOffTokenHandlers() {
	s.Nil(s.provider.GetTargetOneOffTokenHandlers())
	handlers := []string{".*"}
	s.provider.Init(contract.Config{
		TargetOneOffTokenHandlers: &handlers,
	})
	s.Equal(&handlers, s.provider.GetTargetOneOffTokenHandlers())
}

func (s *TestSuite) TestProvider_GetClientProvider() {
	s.Nil(s.provider.GetClientProvider())
	s.provider.Init(contract.Config{
		Client: contract.ClientConfig{
			Provider: mockApiClientProvider{},
		},
	})
	s.NotNil(s.provider.GetClientProvider())
}

func (s *TestSuite) TestProvider_IsClientScopeAccessModelEnabled() {
	s.False(s.provider.IsClientScopeAccessModelEnabled())
	enabled := true
	s.provider.Init(contract.Config{
		Client: contract.ClientConfig{
			UseScopeAccessModel: &enabled,
		},
	})
	s.True(s.provider.IsClientScopeAccessModelEnabled())
}

func (s *TestSuite) TestProvider_ShouldExcludeOptionsRequests() {
	s.False(s.provider.ShouldExcludeOptionsRequests())
	enabled := true
	s.provider.Init(contract.Config{
		ExcludeOptionsRequests: &enabled,
	})
	s.True(s.provider.ShouldExcludeOptionsRequests())
}

func (s *TestSuite) TestProvider_GetClientScopeAccessChecker() {
	s.NotNil(s.provider.GetClientScopeAccessChecker())
	s.provider.Init(contract.Config{
		Client: contract.ClientConfig{
			AccessScopeChecker: checker.PathAccessScopeChecker{},
		},
	})
	s.NotNil(s.provider.GetClientScopeAccessChecker())
}

func (s *TestSuite) TestProvider_GetClientFUPChecker() {
	s.Nil(s.provider.GetClientFUPChecker())
	s.provider.Init(contract.Config{
		Client: contract.ClientConfig{
			FUPChecker: mockFUPChecker{},
		},
	})
	s.NotNil(s.provider.GetClientFUPChecker())
}

func (s *TestSuite) TestProvider_GetUserProvider() {
	s.Nil(s.provider.GetUserProvider())
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			Provider: mockApiUserProvider{},
		},
	})
	s.NotNil(s.provider.GetUserProvider())
}

func (s *TestSuite) TestProvider_IsUserScopeAccessModelEnabled() {
	s.False(s.provider.IsUserScopeAccessModelEnabled())
	enabled := true
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			UseScopeAccessModel: &enabled,
		},
	})
	s.True(s.provider.IsUserScopeAccessModelEnabled())
}

func (s *TestSuite) TestProvider_GetUserScopeAccessChecker() {
	s.NotNil(s.provider.GetUserScopeAccessChecker())
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			AccessScopeChecker: checker.PathAccessScopeChecker{},
		},
	})
	s.NotNil(s.provider.GetUserScopeAccessChecker())
}

func (s *TestSuite) TestProvider_GetUserFUPChecker() {
	s.Nil(s.provider.GetUserFUPChecker())
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			FUPChecker: mockFUPChecker{},
		},
	})
	s.NotNil(s.provider.GetUserFUPChecker())
}

func (s *TestSuite) TestProvider_GetApiTokenExpirationInterval() {
	s.Equal(defaultExpirationInterval, s.provider.GetApiTokenExpirationInterval())
	interval := time.Hour
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			ApiTokenExpirationInterval: &interval,
		},
	})
	s.Equal(interval, s.provider.GetApiTokenExpirationInterval())
}

func (s *TestSuite) TestProvider_GetTokenFactory() {
	s.Nil(s.provider.GetTokenFactory())
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			TokenFactory: func() contract.ApiUserTokenInterface {
				return nil
			},
		},
	})
	s.NotNil(s.provider.GetTokenFactory())
}

func (s *TestSuite) TestProvider_IsUserRegistrationEnabled() {
	s.False(s.provider.IsUserRegistrationEnabled())
	enabled := true
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			WithRegistration: &enabled,
		},
	})
	s.True(s.provider.IsUserRegistrationEnabled())
}

func (s *TestSuite) TestProvider_GetConfirmationTokenExpirationInterval() {
	s.Equal(defaultConfirmationExpirationInterval, s.provider.GetConfirmationTokenExpirationInterval())
	interval := time.Hour
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			ConfirmationTokenExpirationInterval: &interval,
		},
	})
	s.Equal(interval, s.provider.GetConfirmationTokenExpirationInterval())
}

func (s *TestSuite) TestProvider_GetCacheDriver() {
	s.Nil(s.provider.GetCacheDriver())
	s.provider.Init(contract.Config{
		Cache: &contract.CacheConfig{
			Driver: &cache.MemoryCacheDriver{},
		},
	})
	s.NotNil(s.provider.GetCacheDriver())
}

func (s *TestSuite) TestProvider_GetCachePrefix() {
	s.Equal(defaultCachePrefix, s.provider.GetCachePrefix())
	prefix := "prefix"
	s.provider.Init(contract.Config{
		Cache: &contract.CacheConfig{
			Prefix: &prefix,
		},
	})
	s.Equal(prefix, s.provider.GetCachePrefix())
}

func (s *TestSuite) TestProvider_GetCacheTTL() {
	s.Equal(defaultCacheTTL, s.provider.GetCacheTTL())
	ttl := time.Hour
	s.provider.Init(contract.Config{
		Cache: &contract.CacheConfig{
			TTL: &ttl,
		},
	})
	s.Equal(ttl, s.provider.GetCacheTTL())
}

func (s *TestSuite) TestProvider_IsCacheEnabled() {
	s.False(s.provider.IsCacheEnabled())
	s.provider.Init(contract.Config{
		Cache: &contract.CacheConfig{
			Driver: &cache.MemoryCacheDriver{},
		},
	})
	s.True(s.provider.IsCacheEnabled())
}

func (s *TestSuite) TestProvider_IsClientFUPEnabled() {
	s.False(s.provider.IsClientFUPEnabled())
	s.provider.Init(contract.Config{
		Client: contract.ClientConfig{
			FUPChecker: mockFUPChecker{},
		},
	})
	s.True(s.provider.IsClientFUPEnabled())
}

func (s *TestSuite) TestProvider_IsUserFUPEnabled() {
	s.False(s.provider.IsUserFUPEnabled())
	s.provider.Init(contract.Config{
		User: &contract.UserConfig{
			FUPChecker: mockFUPChecker{},
		},
	})
	s.True(s.provider.IsUserFUPEnabled())
}

func (s *TestSuite) TestProvider_IsOneOffTokenModeEnabled() {
	s.False(s.provider.IsOneOffTokenModeEnabled())
	enabled := true
	s.provider.Init(contract.Config{
		Mode: &contract.ModesConfig{
			OneOffToken: &enabled,
		},
	})
	s.True(s.provider.IsOneOffTokenModeEnabled())
}

func (s *TestSuite) TestProvider_GetOneOffTokenExpirationInterval() {
	s.Equal(defaultOneOffTokenExpirationInterval, s.provider.GetOneOffTokenExpirationInterval())
	interval := time.Hour * 2
	s.provider.Init(contract.Config{
		Client: contract.ClientConfig{
			OneOffTokenExpirationInterval: &interval,
		},
	})
	s.Equal(interval, s.provider.GetOneOffTokenExpirationInterval())
}

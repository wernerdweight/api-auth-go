package config

import (
	"github.com/wernerdweight/api-auth-go/auth/checker"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"time"
)

type Provider struct {
	config contract.Config
}

func (p *Provider) IsApiKeyModeEnabled() bool {
	return *p.config.Mode.ApiKey
}

func (p *Provider) IsClientIdAndSecretModeEnabled() bool {
	return *p.config.Mode.ClientIdAndSecret
}

func (p *Provider) GetTargetHandlers() *[]string {
	return p.config.TargetHandlers
}

func (p *Provider) GetClientProvider() contract.ApiClientProviderInterface[contract.ApiClientInterface] {
	return p.config.Client.Provider
}

func (p *Provider) IsClientScopeAccessModelEnabled() bool {
	return *p.config.Client.UseScopeAccessModel
}

func (p *Provider) ShouldExcludeOptionsRequests() bool {
	return *p.config.ExcludeOptionsRequests
}

func (p *Provider) GetClientScopeAccessChecker() contract.AccessScopeCheckerInterface {
	return p.config.Client.AccessScopeChecker
}

func (p *Provider) GetUserProvider() contract.ApiUserProviderInterface[contract.ApiUserInterface] {
	return p.config.User.Provider
}

func (p *Provider) IsUserScopeAccessModelEnabled() bool {
	return *p.config.User.UseScopeAccessModel
}

func (p *Provider) GetUserScopeAccessChecker() contract.AccessScopeCheckerInterface {
	return p.config.User.AccessScopeChecker
}

func (p *Provider) GetApiTokenExpirationInterval() time.Duration {
	return *p.config.User.ApiTokenExpirationInterval
}

func (p *Provider) GetTokenFactory() func() contract.ApiUserTokenInterface {
	return p.config.User.TokenFactory
}

func (p *Provider) IsUserRegistrationEnabled() bool {
	return *p.config.User.WithRegistration
}

func (p *Provider) GetConfirmationTokenExpirationInterval() time.Duration {
	return *p.config.User.ConfirmationTokenExpirationInterval
}

func (p *Provider) GetCacheDriver() contract.CacheDriverInterface {
	return p.config.Cache.Driver
}

func (p *Provider) GetCachePrefix() string {
	return *p.config.Cache.Prefix
}

func (p *Provider) GetCacheTTL() time.Duration {
	return *p.config.Cache.TTL
}

func (p *Provider) IsCacheEnabled() bool {
	return nil != p.config.Cache.Driver
}

func (p *Provider) initUser(config contract.Config) {
	if nil != config.User.Provider {
		p.config.User.Provider = config.User.Provider
	}
	if nil != config.User.TokenFactory {
		p.config.User.TokenFactory = config.User.TokenFactory
	}
	if nil != config.User.ApiTokenExpirationInterval {
		p.config.User.ApiTokenExpirationInterval = config.User.ApiTokenExpirationInterval
	}
	if nil != config.User.UseScopeAccessModel {
		p.config.User.UseScopeAccessModel = config.User.UseScopeAccessModel
	}
	if nil != config.User.AccessScopeChecker {
		p.config.User.AccessScopeChecker = config.User.AccessScopeChecker
	}
	if nil != config.User.WithRegistration {
		p.config.User.WithRegistration = config.User.WithRegistration
	}
	if nil != config.User.ConfirmationTokenExpirationInterval {
		p.config.User.ConfirmationTokenExpirationInterval = config.User.ConfirmationTokenExpirationInterval
	}
}

func (p *Provider) initMode(config contract.Config) {
	if nil != config.Mode.ApiKey {
		p.config.Mode.ApiKey = config.Mode.ApiKey
	}
	if nil != config.Mode.ClientIdAndSecret {
		p.config.Mode.ClientIdAndSecret = config.Mode.ClientIdAndSecret
	}
}

func (p *Provider) initCache(config contract.Config) {
	if nil != config.Cache.Driver {
		p.config.Cache.Driver = config.Cache.Driver
	}
	if nil != config.Cache.Prefix && "" != *config.Cache.Prefix {
		p.config.Cache.Prefix = config.Cache.Prefix
	}
	if nil != config.Cache.TTL {
		p.config.Cache.TTL = config.Cache.TTL
	}
}

func (p *Provider) Init(config contract.Config) {
	if nil != config.Client.Provider {
		p.config.Client.Provider = config.Client.Provider
	}
	if nil != config.Client.UseScopeAccessModel {
		p.config.Client.UseScopeAccessModel = config.Client.UseScopeAccessModel
	}
	if nil != config.Client.AccessScopeChecker {
		p.config.Client.AccessScopeChecker = config.Client.AccessScopeChecker
	}

	if nil != config.User {
		p.initUser(config)
	}

	if nil != config.Mode {
		p.initMode(config)
	}

	if nil != config.TargetHandlers {
		p.config.TargetHandlers = config.TargetHandlers
	}

	if nil != config.ExcludeOptionsRequests {
		p.config.ExcludeOptionsRequests = config.ExcludeOptionsRequests
	}

	if nil != config.Cache {
		p.initCache(config)
	}
}

var (
	defaultApiKeyMode                     = false
	defaultClientIdAndSecretMode          = true
	defaultExcludeOptionsRequests         = false
	defaultClientUseScopeAccessModel      = false
	defaultUserUseScopeAccessModel        = false
	defaultWithRegistration               = false
	defaultExpirationInterval             = time.Hour * 24 * 30
	defaultConfirmationExpirationInterval = time.Hour * 12
	defaultCacheTTL                       = time.Hour
	defaultCachePrefix                    = "api-auth-go:"
)

var ProviderInstance = &Provider{
	config: contract.Config{
		Client: contract.ClientConfig{
			Provider:            nil,
			UseScopeAccessModel: &defaultClientUseScopeAccessModel,
			AccessScopeChecker:  checker.PathAccessScopeChecker{},
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
			ClientIdAndSecret: &defaultClientIdAndSecretMode,
		},
		TargetHandlers:         nil,
		ExcludeOptionsRequests: &defaultExcludeOptionsRequests,
		Cache: &contract.CacheConfig{
			Driver: nil,
			Prefix: &defaultCachePrefix,
			TTL:    &defaultCacheTTL,
		},
	},
}

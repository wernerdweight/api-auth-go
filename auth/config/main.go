package config

import (
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
		if nil != config.User.Provider {
			p.config.User.Provider = config.User.Provider
		}
		if nil != config.User.TokenProvider {
			p.config.User.TokenProvider = config.User.TokenProvider
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
	}

	if nil != config.Mode {
		if nil != config.Mode.ApiKey {
			p.config.Mode.ApiKey = config.Mode.ApiKey
		}
		if nil != config.Mode.ClientIdAndSecret {
			p.config.Mode.ClientIdAndSecret = config.Mode.ClientIdAndSecret
		}
	}

	if nil != config.TargetHandlers {
		p.config.TargetHandlers = config.TargetHandlers
	}

	if nil != config.ExcludeOptionsRequests {
		p.config.ExcludeOptionsRequests = config.ExcludeOptionsRequests
	}
}

var (
	defaultApiKeyMode                = false
	defaultClientIdAndSecretMode     = true
	defaultExcludeOptionsRequests    = false
	defaultClientUseScopeAccessModel = false
	defaultUserUseScopeAccessModel   = false
	defaultExpirationInterval        = time.Hour * 24 * 30
)

var ProviderInstance = &Provider{
	config: contract.Config{
		Client: contract.ClientConfig{
			Provider:            nil,
			UseScopeAccessModel: &defaultClientUseScopeAccessModel,
			AccessScopeChecker:  nil,
		},
		User: &contract.UserConfig{
			Provider:                   nil,
			TokenProvider:              nil,
			ApiTokenExpirationInterval: &defaultExpirationInterval,
			UseScopeAccessModel:        &defaultUserUseScopeAccessModel,
			AccessScopeChecker:         nil,
		},
		Mode: &contract.ModesConfig{
			ApiKey:            &defaultApiKeyMode,
			ClientIdAndSecret: &defaultClientIdAndSecretMode,
		},
		TargetHandlers:         nil,
		ExcludeOptionsRequests: &defaultExcludeOptionsRequests,
	},
}

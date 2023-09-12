package contract

import "time"

type ClientConfig struct {
	// Provider: your provider that implements ApiClientProviderInterface
	Provider ApiClientProviderInterface[ApiClientInterface]
	// UseScopeAccessModel: if set to true, client scope will be checked before granting access (see `scope access` below) - default false
	UseScopeAccessModel *bool
	// AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to PathAccessScopeChecker
	AccessScopeChecker AccessScopeCheckerInterface
}

type UserConfig struct {
	// Provider: your provider that implements ApiUserProviderInterface
	Provider ApiUserProviderInterface[ApiUserInterface]
	// TokenFactory: your token type that implements ApiUserTokenInterface
	TokenFactory func() ApiUserTokenInterface
	// ApiTokenExpirationInterval: token expiration in seconds - defaults to 2,592,000 (30 days)
	ApiTokenExpirationInterval *time.Duration
	// UseScopeAccessModel: if set to true, user scope will be checked before granting access (see `scope access` below) - default false
	UseScopeAccessModel *bool
	// AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to PathAccessScopeChecker
	AccessScopeChecker AccessScopeCheckerInterface
}

type ModesConfig struct {
	// ApiKey: api key authentication mode (optional; default false)
	ApiKey *bool
	// ClientIdAndSecret: client id and secret authentication mode (optional; default true)
	ClientIdAndSecret *bool
}

type Config struct {
	// Client: api client configuration (mandatory)
	Client ClientConfig

	// User: api user configuration (optional; if you omit user configuration, you will not be able to use `on-behalf` access mode (see below))
	User *UserConfig

	// Mode: modes of authentication (client id + secret and user token vs. api key)
	Mode *ModesConfig

	// TargetHandlers: list of handlers to target (optional)
	TargetHandlers *[]string
	// '*'   # all handlers
	// TODO: 'My\Controller\SomeInterface'
	// TODO: 'Vendor\Bundle\Controller\SomeOtherInterface'

	// ExcludeOptionsRequests: if true, requests using the OPTIONS method will be ignored (authentication will be skipped) - default false
	ExcludeOptionsRequests *bool
}

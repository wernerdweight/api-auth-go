package contract

import "time"

type ClientConfig struct {
	// Provider: your provider that implements ApiClientProviderInterface
	Provider ApiClientProviderInterface[ApiClientInterface]
	// UseScopeAccessModel: if set to true, client scope will be checked before granting access (see `scope access` below) - default false
	UseScopeAccessModel *bool
	// AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to PathAccessScopeChecker
	AccessScopeChecker AccessScopeCheckerInterface
	// FUPChecker: the checker used to check FUP limits that implements FUPCheckerInterface (optional; if you omit FUP checker, FUP limits will not be checked)
	// NOTE: if you want to use FUP limits, you must also enable Cache (see below)
	FUPChecker FUPCheckerInterface
}

type UserConfig struct {
	// Provider: your provider that implements ApiUserProviderInterface
	Provider ApiUserProviderInterface[ApiUserInterface]
	// TokenFactory: generates your token type that implements ApiUserTokenInterface
	TokenFactory func() ApiUserTokenInterface
	// ApiTokenExpirationInterval: token expiration in seconds - defaults to 2,592,000 (30 days)
	ApiTokenExpirationInterval *time.Duration
	// UseScopeAccessModel: if set to true, user scope will be checked before granting access (see `scope access` below) - default false
	UseScopeAccessModel *bool
	// AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to PathAccessScopeChecker
	AccessScopeChecker AccessScopeCheckerInterface
	// WithRegistration: if set to true, user registration will be enabled - default false
	WithRegistration *bool
	// ConfirmationTokenExpirationInterval: confirmation token expiration in seconds - defaults to 43200 (12 hours)
	ConfirmationTokenExpirationInterval *time.Duration
	// FUPChecker: the checker used to check FUP limits that implements FUPCheckerInterface (optional; if you omit FUP checker, FUP limits will not be checked)
	// NOTE: if you want to use FUP limits, you must also enable Cache (see below)
	FUPChecker FUPCheckerInterface
}

type ModesConfig struct {
	// ApiKey: api key authentication mode (optional; default false)
	ApiKey *bool
	// ClientIdAndSecret: client id and secret authentication mode (optional; default true)
	ClientIdAndSecret *bool
}

type CacheConfig struct {
	// Driver: your cache driver that implements CacheDriverInterface
	Driver CacheDriverInterface
	// Prefix: prefix to use for cache keys - defaults to `api-auth-go:`
	Prefix *string
	// TTL: cache TTL in seconds - defaults to 3600 (1 hour)
	TTL *time.Duration
}

type Config struct {
	// Client: api client configuration (mandatory)
	Client ClientConfig

	// User: api user configuration (optional; if you omit user configuration, you will not be able to use `on-behalf` access mode (see below))
	User *UserConfig

	// Mode: modes of authentication (client id + secret and user token vs. api key)
	Mode *ModesConfig

	// TargetHandlers: list of handlers to target (optional; if you omit target handlers, all handlers will be targeted)
	TargetHandlers *[]string
	// '.*'            	# all handlers
	// '/v1/*'   		# all handlers starting with '/v1/'
	// '/v1/some/path'  # only '/v1/some/path' handler

	// ExcludeHandlers: list of handlers to exclude (optional; if you omit exclude handlers, no handlers will be excluded)
	ExcludeHandlers *[]string

	// ExcludeOptionsRequests: if true, requests using the OPTIONS method will be ignored (authentication will be skipped) - default false
	ExcludeOptionsRequests *bool

	// Cache: cache configuration (optional)
	Cache *CacheConfig
}

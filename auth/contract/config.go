package contract

type clientConfig struct {
	// Provider: your provider that implements ApiClientProviderInterface
	Provider ApiClientProviderInterface[ApiClientInterface]
	// UseScopeAccessModel: if set to true, client scope will be checked before granting access (see `scope access` below) - default false
	UseScopeAccessModel bool
	// AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to RouteChecker (see below)
	AccessScopeChecker AccessScopeCheckerInterface
}

type userConfig struct {
	// Provider: your provider that implements ApiUserProviderInterface
	Provider ApiUserProviderInterface[ApiUserInterface]
	// TokenProvider: your provider that implements ApiUserTokenProviderInterface
	// TODO: maybe we need factory instead of provider?
	TokenProvider ApiUserTokenProviderInterface[ApiUserTokenInterface]
	// ApiTokenExpirationInterval: token expiration in seconds - defaults to 2,592,000 (30 days)
	ApiTokenExpirationInterval int
	// UseScopeAccessModel: if set to true, user scope will be checked before granting access (see `scope access` below) - default false
	UseScopeAccessModel bool
	// AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to RouteChecker (see below)
	AccessScopeChecker AccessScopeCheckerInterface
}

type Config struct {
	// Client: api client configuration (mandatory)
	Client clientConfig

	// User: api user configuration (optional; if you omit user configuration, you will not be able to use `on-behalf` access mode (see below))
	User userConfig

	// TargetHandlers: list of handlers to target (optional)
	TargetHandlers []string
	// '*'   # all handlers
	// TODO: 'My\Controller\SomeInterface'
	// TODO: 'Vendor\Bundle\Controller\SomeOtherInterface'

	// ExcludeOptionsRequests: if true, requests using the OPTIONS method will be ignored (authentication will be skipped) - default false
	ExcludeOptionsRequests bool
}

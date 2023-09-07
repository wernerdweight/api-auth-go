package constants

const (
	ClientIdHeader               = "X-Client-Id"
	ClientSecretHeader           = "X-Client-Secret"
	ApiUserTokenHeader           = "X-Api-User-Token"
	ApiKeyHeader                 = "Authorization"
	ScopeAccessibilityAccessible = "true"
	ScopeAccessibilityForbidden  = "false"
	ScopeAccessibilityOnBehalf   = "on-behalf"
	RouteKey                     = "_route"
	RouteOverrideKey             = "_route_override"
	Realm                        = "Basic realm=\"API\""

	ApiClient = "api-client"
	ApiUser   = "api-user"
)

var ScopeAccessibilityOptions = []string{
	ScopeAccessibilityAccessible,
	ScopeAccessibilityForbidden,
	ScopeAccessibilityOnBehalf,
}

package constants

type ScopeAccessibility string

const (
	ClientIdHeader                                  = "X-Client-Id"
	ClientSecretHeader                              = "X-Client-Secret"
	ApiUserTokenHeader                              = "X-Api-User-Token"
	ApiKeyHeader                                    = "Authorization"
	ScopeAccessibilityAccessible ScopeAccessibility = "true"
	ScopeAccessibilityForbidden  ScopeAccessibility = "false"
	ScopeAccessibilityOnBehalf   ScopeAccessibility = "on-behalf"
	RouteKey                                        = "_route"
	RouteOverrideKey                                = "_route_override"
	Realm                                           = "Basic realm=\"API\""

	ApiClient = "api-client"
	ApiUser   = "api-user"
)

var ScopeAccessibilityOptions = []ScopeAccessibility{
	ScopeAccessibilityAccessible,
	ScopeAccessibilityForbidden,
	ScopeAccessibilityOnBehalf,
}

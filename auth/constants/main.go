package constants

const (
	ClientIdHeader               = "X-Client-Id"
	ClientSecretHeader           = "X-Client-Secret"
	ApiUserTokenHeader           = "X-Api-User-Token"
	ScopeAccessibilityAccessible = "true"
	ScopeAccessibilityForbidden  = "false"
	ScopeAccessibilityOnBehalf   = "on-behalf"
	RouteKey                     = "_route"
	RouteOverrideKey             = "_route_override"
	Realm                        = "Basic realm=\"API\""
)

var ScopeAccessibilityOptions = []string{
	ScopeAccessibilityAccessible,
	ScopeAccessibilityForbidden,
	ScopeAccessibilityOnBehalf,
}

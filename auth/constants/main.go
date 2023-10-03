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
	DefaultTokenLength                              = 32

	ApiClient = "api-client"
	ApiUser   = "api-user"
)

var ScopeAccessibilityOptions = []ScopeAccessibility{
	ScopeAccessibilityAccessible,
	ScopeAccessibilityForbidden,
	ScopeAccessibilityOnBehalf,
}

var FUPScopePeriods = []string{
	"minutely",
	"hourly",
	"daily",
	"weekly",
	"monthly",
}

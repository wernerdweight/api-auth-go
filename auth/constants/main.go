package constants

type ScopeAccessibility string

const (
	ClientIdHeader                                  = "X-Client-Id"
	ClientSecretHeader                              = "X-Client-Secret"
	ApiUserTokenHeader                              = "X-Api-User-Token"
	ApiKeyHeader                                    = "Authorization"
	ClientFUPLimitsHeader                           = "X-Client-FUP-Limits"
	UserFUPLimitsHeader                             = "X-User-FUP-Limits"
	ScopeAccessibilityAccessible ScopeAccessibility = "true"
	ScopeAccessibilityForbidden  ScopeAccessibility = "false"
	ScopeAccessibilityOnBehalf   ScopeAccessibility = "on-behalf"
	ScopeAccessibilityUnlimited  ScopeAccessibility = "unlimited"
	DefaultTokenLength                              = 32
	PeriodMinutely               Period             = "minutely"
	PeriodHourly                 Period             = "hourly"
	PeriodDaily                  Period             = "daily"
	PeriodWeekly                 Period             = "weekly"
	PeriodMonthly                Period             = "monthly"

	ApiClient = "api-client"
	ApiUser   = "api-user"
)

var ScopeAccessibilityOptions = []ScopeAccessibility{
	ScopeAccessibilityAccessible,
	ScopeAccessibilityForbidden,
	ScopeAccessibilityOnBehalf,
}

var FUPScopeAccessibilityOptions = []ScopeAccessibility{
	ScopeAccessibilityAccessible,
	ScopeAccessibilityForbidden,
	ScopeAccessibilityUnlimited,
}

type Period string

var FUPScopePeriods = []Period{
	PeriodMinutely,
	PeriodHourly,
	PeriodDaily,
	PeriodWeekly,
	PeriodMonthly,
}

package constants

import (
	"fmt"
	"github.com/jinzhu/now"
	"time"
)

type ScopeAccessibility string

const (
	ClientIdHeader                                  = "X-Client-Id"
	ClientSecretHeader                              = "X-Client-Secret"
	ApiUserTokenHeader                              = "X-Api-User-Token"
	ApiKeyHeader                                    = "Authorization"
	OneOffTokenHeader                               = "X-Token"
	ClientFUPLimitsHeader                           = "X-Client-FUP-Limits"
	UserFUPLimitsHeader                             = "X-User-FUP-Limits"
	RetryAfterHeader                                = "Retry-After"
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
	FUPIPKey                                        = "per-ip"
	FUPCookieKey                                    = "per-cookie"

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

func (p Period) GetFormatToCompare(t time.Time) string {
	switch p {
	case PeriodMinutely:
		return t.Format("2006-01-02 15:04")
	case PeriodHourly:
		return t.Format("2006-01-02 15")
	case PeriodDaily:
		return t.Format("2006-01-02")
	case PeriodWeekly:
		year, week := t.ISOWeek()
		return fmt.Sprintf("%d-%d", year, week)
	case PeriodMonthly:
		return t.Format("2006-01")
	}
	return ""
}

func (p Period) GetResetTime() time.Time {
	switch p {
	case PeriodMinutely:
		return now.EndOfMinute()
	case PeriodHourly:
		return now.EndOfHour()
	case PeriodDaily:
		return now.EndOfDay()
	case PeriodWeekly:
		return now.EndOfWeek()
	case PeriodMonthly:
		return now.EndOfMonth()
	}
	return time.Now()
}

var FUPScopePeriods = []Period{
	PeriodMinutely,
	PeriodHourly,
	PeriodDaily,
	PeriodWeekly,
	PeriodMonthly,
}

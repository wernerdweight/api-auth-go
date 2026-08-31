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
	OneOffTokenLength                               = 64
	PeriodMinutely               Period             = "minutely"
	PeriodHourly                 Period             = "hourly"
	PeriodDaily                  Period             = "daily"
	PeriodWeekly                 Period             = "weekly"
	PeriodMonthly                Period             = "monthly"
	FUPIPKey                                        = "per-ip"
	FUPCookieKey                                    = "per-cookie"
	AnonymousFUPKey                                 = "anonymous"

	ApiClient = "api-client"
	ApiUser   = "api-user"
)

// FUPEntryTTL is the expiration of FUP cache entries whose limited periods are not known to the
// caller. It is slightly longer than the longest FUP period (monthly), so that counters that are
// no longer used (e.g. per-IP counters of one-off visitors) are released instead of growing
// unbounded.
const FUPEntryTTL = time.Hour * 24 * 35

// GetEntryTTL returns how long a FUP cache entry has to survive inactivity to keep counting this
// period correctly - the period itself plus a margin, since the expiration is refreshed on every
// increment and the entry only has to outlive the gap between two requests that share a period.
// A scope that limits by nothing longer than a day therefore keeps its entries for two days
// instead of the 35 the monthly period needs, which matters for keys with an unbounded key space
// (per IP, per cookie).
func (p Period) GetEntryTTL() time.Duration {
	switch p {
	case PeriodMinutely:
		return time.Hour
	case PeriodHourly:
		return time.Hour * 25
	case PeriodDaily:
		return time.Hour * 24 * 2
	case PeriodWeekly:
		return time.Hour * 24 * 8
	case PeriodMonthly:
		return FUPEntryTTL
	}
	return FUPEntryTTL
}

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

// GetTimestampBounds returns the inclusive lexicographical bounds of the period that contains t.
// A timestamp belongs to the same period as t if its RFC 3339 representation truncated to the
// length of the returned bounds is not lower than the first and not higher than the second one.
// Both bounds always have the same length, and periods other than weekly are a single value.
//
// This is the same comparison as GetFormatToCompare, expressed so that it can be evaluated by a
// cache backend (which only sees the serialized timestamp) instead of in Go.
//
// Both the bounds and the timestamps they are compared against are expressed in the local time of
// the application instance (as GetFormatToCompare and GetResetTime are), so instances that run in
// different time zones and share a cache do not agree on where a period starts and ends.
func (p Period) GetTimestampBounds(t time.Time) (string, string) {
	switch p {
	case PeriodMinutely:
		return t.Format("2006-01-02T15:04"), t.Format("2006-01-02T15:04")
	case PeriodHourly:
		return t.Format("2006-01-02T15"), t.Format("2006-01-02T15")
	case PeriodDaily:
		return t.Format("2006-01-02"), t.Format("2006-01-02")
	case PeriodWeekly:
		// two timestamps share an ISO week exactly if they fall within the same Monday to Sunday span
		monday := t.AddDate(0, 0, -((int(t.Weekday()) + 6) % 7))
		return monday.Format("2006-01-02"), monday.AddDate(0, 0, 6).Format("2006-01-02")
	case PeriodMonthly:
		return t.Format("2006-01"), t.Format("2006-01")
	}
	return "", ""
}

var FUPScopePeriods = []Period{
	PeriodMinutely,
	PeriodHourly,
	PeriodDaily,
	PeriodWeekly,
	PeriodMonthly,
}

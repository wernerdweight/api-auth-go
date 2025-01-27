package contract

import (
	"encoding/json"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"log"
	"time"
)

type FUPLimits struct {
	Limit  int              `json:"limit"`
	Used   int              `json:"used"`
	Period constants.Period `json:"-"`
}

type FUPScopeLimits struct {
	Accessible constants.ScopeAccessibility
	Limits     map[constants.Period]FUPLimits
	Error      *AuthError
}

func (l *FUPScopeLimits) GetLimitsHeader() string {
	if nil == l.Limits {
		return ""
	}
	header, err := json.Marshal(l.Limits)
	if nil != err {
		log.Printf("can't serialize FUP limits header: %+v", err)
		return ""
	}
	return string(header)
}

func (l *FUPScopeLimits) GetRetryAfter() int {
	if l.Accessible != constants.ScopeAccessibilityForbidden {
		return -1
	}
	for period := range l.Limits {
		// there will always be exactly one limit
		return int(time.Until(period.GetResetTime()).Seconds())
	}
	return -1
}

type FUPCacheEntry struct {
	UpdatedAt time.Time                `json:"updatedAt"`
	Used      map[constants.Period]int `json:"used"`
}

func (e *FUPCacheEntry) GetUsed(period constants.Period) int {
	if nil == e.Used {
		return 0
	}
	return e.Used[period]
}

func (e *FUPCacheEntry) Increment() {
	if nil == e.Used {
		e.Used = make(map[constants.Period]int)
		e.UpdatedAt = time.Now()
	}
	now := time.Now()
	for _, period := range constants.FUPScopePeriods {
		if period.GetFormatToCompare(e.UpdatedAt) == period.GetFormatToCompare(now) {
			e.Used[period]++
			continue
		}
		e.Used[period] = 1
	}
	e.UpdatedAt = time.Now()
}

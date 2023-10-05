package contract

import (
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"time"
)

type FUPLimits struct {
	Limit  int
	Used   int
	Period constants.Period
}

type FUPScopeLimits struct {
	Accessible constants.ScopeAccessibility
	Limits     map[constants.Period]FUPLimits
	Error      *AuthError
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
	}
	for _, period := range constants.FUPScopePeriods {
		e.Used[period]++
	}
}

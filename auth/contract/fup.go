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
	Limits     []FUPLimits
	Error      error
}

type FUPCacheEntry struct {
	UpdatedAt time.Time
	Used      map[constants.Period]int
}

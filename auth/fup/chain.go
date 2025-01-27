package fup

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
)

// ChainFUPChecker allows to chain multiple FUPCheckerInterface implementations
type ChainFUPChecker struct {
	Checkers []contract.FUPCheckerInterface
}

func (ch ChainFUPChecker) Check(scope *contract.FUPScope, c *gin.Context, key string) contract.FUPScopeLimits {
	limits := contract.FUPScopeLimits{Accessible: constants.ScopeAccessibilityUnlimited}
	for _, checker := range ch.Checkers {
		checkerLimits := checker.Check(scope, c, key)
		if constants.ScopeAccessibilityForbidden == checkerLimits.Accessible {
			return checkerLimits
		}
		if checkerLimits.Limits != nil {
			limits.Limits = mergeLimits(limits.Limits, checkerLimits.Limits)
			limits.Accessible = checkerLimits.Accessible
		}
	}
	return limits
}

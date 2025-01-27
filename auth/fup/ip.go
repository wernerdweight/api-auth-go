package fup

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
)

// IPFUPChecker is an implementation of the FUPCheckerInterface for the URL path-based access model
type IPFUPChecker struct {
}

func (ch IPFUPChecker) Check(scope *contract.FUPScope, c *gin.Context, key string) contract.FUPScopeLimits {
	ip := c.ClientIP()
	if nil == scope || "" == ip {
		// no limitations by default
		return contract.FUPScopeLimits{Accessible: constants.ScopeAccessibilityUnlimited}
	}
	if !scope.HasLimit(constants.FUPIPKey) {
		return contract.FUPScopeLimits{Accessible: constants.ScopeAccessibilityUnlimited}
	}
	if !config.ProviderInstance.IsCacheEnabled() {
		return contract.FUPScopeLimits{
			Error: contract.NewInternalError(contract.FUPCacheDisabled, nil),
		}
	}
	cacheDriver := config.ProviderInstance.GetCacheDriver()
	ipLimits, scopeLimits := checkLimits(scope, key, ip, constants.FUPIPKey, cacheDriver)
	if nil != scopeLimits {
		return *scopeLimits
	}
	return contract.FUPScopeLimits{
		Accessible: constants.ScopeAccessibilityAccessible,
		Limits:     ipLimits,
		Error:      nil,
	}
}

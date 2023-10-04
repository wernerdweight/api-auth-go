package fup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"strings"
)

// PathFUPChecker is an implementation of the FUPCheckerInterface for the URL path-based access model
type PathFUPChecker struct {
}

func (ch PathFUPChecker) Check(scope *contract.FUPScope, c *gin.Context, key string) contract.FUPScopeLimits {
	limits := make([]contract.FUPLimits, 0)
	if nil == scope || nil == c || nil == c.Request || nil == c.Request.URL {
		// no limitations by default
		return contract.FUPScopeLimits{Accessible: constants.ScopeAccessibilityUnlimited}
	}
	// TODO: first check if there are any limitations for "*" path (all paths - i.e. limit for all paths/per-account limits)

	path := strings.ToLower(c.Request.URL.Path)
	cachePrefix := config.ProviderInstance.GetCachePrefix()
	cacheKey := fmt.Sprintf("%s_fup_%s_%s", cachePrefix, key, strings.Replace(path, "/", "-", -1))
	// TODO: fetch FUP cache for current path (memory, redis, etc. - needs to be configurable); also handle errors
	// TODO: use the same driver as for the access scope cache
	// TODO: if no driver is configured, return error
	var cacheEntry contract.FUPCacheEntry
	// TODO: increment FUP cache for current path and save it back to cache

	for _, period := range constants.FUPScopePeriods {
		limit := scope.GetLimit(fmt.Sprintf("%s.%s", path, period))
		if nil == limit || *limit < 0 {
			// no limitations by default
			return contract.FUPScopeLimits{Accessible: constants.ScopeAccessibilityUnlimited}
		}
		used := cacheEntry.Used[period]
		if *limit < used {
			return contract.FUPScopeLimits{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits: []contract.FUPLimits{{
					Limit:  *limit,
					Used:   used,
					Period: period,
				}},
				Error: nil,
			}
		}
		limits = append(limits, contract.FUPLimits{
			Limit:  *limit,
			Used:   used,
			Period: period,
		})
	}
	return contract.FUPScopeLimits{
		Accessible: constants.ScopeAccessibilityAccessible,
		Limits:     limits,
		Error:      nil,
	}
}

func (ch PathFUPChecker) Log(c *gin.Context) error {
	// TODO: only log if there are any limitations for given path (but decide upstream based on limit prwsent in scope, not here)
	// TODO: we need to get client/user ID here (not from the context)
	// TODO: implement
	return nil
}

package fup

import (
	"fmt"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"strings"
)

func checkLimits(scope *contract.FUPScope, key string, cacheId string, path string, cacheDriver contract.CacheDriverInterface) (map[constants.Period]contract.FUPLimits, *contract.FUPScopeLimits) {
	limits := make(map[constants.Period]contract.FUPLimits)
	cacheKey := fmt.Sprintf("fup_%s_%s", key, strings.Replace(cacheId, "/", "-", -1))
	cacheEntry, err := cacheDriver.GetFUPEntry(cacheKey)
	if nil != err {
		return nil, &contract.FUPScopeLimits{
			Error: err,
		}
	}
	cacheEntry.Increment()
	err = cacheDriver.SetFUPEntry(cacheKey, cacheEntry)
	if nil != err {
		return nil, &contract.FUPScopeLimits{
			Error: err,
		}
	}

	for _, period := range constants.FUPScopePeriods {
		limit := scope.GetLimit(fmt.Sprintf("%s.%s", path, period))
		if nil == limit || *limit < 0 {
			// no limitations by default
			continue
		}
		used := cacheEntry.GetUsed(period)
		if *limit < used {
			return nil, &contract.FUPScopeLimits{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits: map[constants.Period]contract.FUPLimits{
					period: {
						Limit:  *limit,
						Used:   used,
						Period: period,
					},
				},
				Error: nil,
			}
		}
		limits[period] = contract.FUPLimits{
			Limit:  *limit,
			Used:   used,
			Period: period,
		}
	}
	return limits, nil
}

func mergeLimits(limits map[constants.Period]contract.FUPLimits, pathLimits map[constants.Period]contract.FUPLimits) map[constants.Period]contract.FUPLimits {
	if nil == limits {
		return pathLimits
	}
	for period, pathLimit := range pathLimits {
		if limit, ok := limits[period]; ok {
			remainingPathLimit := pathLimit.Limit - pathLimit.Used
			remainingLimit := limit.Limit - limit.Used
			if remainingPathLimit < remainingLimit {
				limits[period] = pathLimit
			}
			continue
		}
		limits[period] = pathLimit
	}
	return limits
}

func check(path string, scope *contract.FUPScope, key string) contract.FUPScopeLimits {
	hasRootLimit := scope.HasLimit("*")
	hasPathLimit := scope.HasLimit(path)
	if !hasRootLimit && !hasPathLimit {
		// no limitations by default
		return contract.FUPScopeLimits{Accessible: constants.ScopeAccessibilityUnlimited}
	}

	if !config.ProviderInstance.IsCacheEnabled() {
		return contract.FUPScopeLimits{
			Error: contract.NewInternalError(contract.FUPCacheDisabled, nil),
		}
	}
	cacheDriver := config.ProviderInstance.GetCacheDriver()

	var limits map[constants.Period]contract.FUPLimits
	if hasRootLimit {
		rootLimits, scopeLimits := checkLimits(scope, key, "*", "*", cacheDriver)
		if nil != scopeLimits {
			return *scopeLimits
		}
		limits = rootLimits
	}

	if hasPathLimit {
		pathLimits, scopeLimits := checkLimits(scope, key, path, path, cacheDriver)
		if nil != scopeLimits {
			return *scopeLimits
		}
		limits = mergeLimits(limits, pathLimits)
	}

	return contract.FUPScopeLimits{
		Accessible: constants.ScopeAccessibilityAccessible,
		Limits:     limits,
		Error:      nil,
	}
}

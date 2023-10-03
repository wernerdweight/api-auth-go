package fup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"strings"
)

// PathFUPChecker is an implementation of the FUPCheckerInterface for the URL path-based access model
type PathFUPChecker struct {
}

// TODO: we will need to return used/limit values here (for response headers)
func (ch PathFUPChecker) Check(scope *contract.FUPScope, c *gin.Context) (bool, error) {
	// TODO: we need to get client/user ID here (not from the context)
	if nil == scope || nil == c || nil == c.Request || nil == c.Request.URL {
		// no limitations by default
		return true, nil
	}
	// TODO: first check if there are any limitations for "*" path (all paths - i.e. limit for all paths/per-account limits)

	path := strings.ToLower(c.Request.URL.Path)
	for _, period := range constants.FUPScopePeriods {
		limit := scope.GetLimit(fmt.Sprintf("%s.%s", path, period))
		if nil == limit || *limit < 0 {
			// no limitations by default
			return true, nil
		}
		// TODO: fetch current state (used) from cache (memory, redis, etc. - needs to be configurable); also handle errors
		used := 0
		if *limit < used {
			return false, nil
		}
	}
	return true, nil
}

func (ch PathFUPChecker) Log(c *gin.Context) error {
	// TODO: only log if there are any limitations for given path (but decide upstream based on limit prwsent in scope, not here)
	// TODO: we need to get client/user ID here (not from the context)
	// TODO: implement
	return nil
}

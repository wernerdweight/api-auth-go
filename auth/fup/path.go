package fup

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"strings"
)

// PathFUPChecker is an implementation of the FUPCheckerInterface for the URL path-based access model
type PathFUPChecker struct {
}

func (ch PathFUPChecker) Check(scope *contract.FUPScope, c *gin.Context, key string) contract.FUPScopeLimits {
	path := strings.ToLower(c.Request.URL.Path)
	return check(path, scope, c, key)
}

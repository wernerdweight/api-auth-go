package fup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"strings"
)

// PathAndMethodFUPChecker is an implementation of the FUPCheckerInterface for the URL path and method-based access model
type PathAndMethodFUPChecker struct {
}

func (ch PathAndMethodFUPChecker) Check(scope *contract.FUPScope, c *gin.Context, key string) contract.FUPScopeLimits {
	path := strings.ToLower(c.Request.URL.Path)
	method := strings.ToLower(c.Request.Method)
	combinedPath := fmt.Sprintf("%s:%s", method, path)
	return check(combinedPath, scope, c, key)
}

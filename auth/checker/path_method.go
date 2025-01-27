package checker

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"strings"
)

// PathAndMethodAccessScopeChecker is an implementation of the AccessScopeCheckerInterface for the URL path and method-based access model
type PathAndMethodAccessScopeChecker struct {
	hierarchySeparator string
}

func (ch PathAndMethodAccessScopeChecker) Check(scope *contract.AccessScope, c *gin.Context) constants.ScopeAccessibility {
	if nil == scope || nil == c || nil == c.Request || nil == c.Request.URL {
		return constants.ScopeAccessibilityForbidden
	}
	path := strings.ToLower(c.Request.URL.Path)
	method := strings.ToLower(c.Request.Method)
	return scope.GetAccessibility(fmt.Sprintf("%s:%s", method, path), ch.hierarchySeparator)
}

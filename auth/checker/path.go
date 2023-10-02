package checker

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"strings"
)

// PathAccessScopeChecker is an implementation of the AccessScopeCheckerInterface for the URL path-based access model
type PathAccessScopeChecker struct {
}

func (ch PathAccessScopeChecker) Check(scope *contract.AccessScope, c *gin.Context) constants.ScopeAccessibility {
	if nil == scope || nil == c || nil == c.Request || nil == c.Request.URL {
		return constants.ScopeAccessibilityForbidden
	}
	path := strings.ToLower(c.Request.URL.Path)
	return scope.GetAccessibility(path)
}

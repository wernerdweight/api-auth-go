package checker

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
)

// PathAccessScopeChecker is an implementation of the AccessScopeCheckerInterface for the URL path-based access model
type PathAccessScopeChecker struct {
}

func (ch PathAccessScopeChecker) Check(scope *contract.AccessScope, c *gin.Context) constants.ScopeAccessibility {
	if nil == scope {
		return constants.ScopeAccessibilityForbidden
	}
	path := c.Request.URL.Path
	return scope.GetAccessibility(path)
}

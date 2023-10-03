package fup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"strings"
)

// PathAndMethodFUPChecker is an implementation of the FUPCheckerInterface for the URL path and method-based access model
type PathAndMethodFUPChecker struct {
}

func (ch PathAndMethodFUPChecker) Check(scope *contract.FUPScope, c *gin.Context) (bool, error) {
	if nil == scope || nil == c || nil == c.Request || nil == c.Request.URL {
		return constants.ScopeAccessibilityForbidden
	}
	path := strings.ToLower(c.Request.URL.Path)
	method := strings.ToLower(c.Request.Method)
	return scope.GetAccessibility(fmt.Sprintf("%s:%s", method, path))
}

func (ch PathAndMethodFUPChecker) Log(c *gin.Context) error {
	// TODO: implement
	return nil
}

package contract

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
)

type AccessScopeCheckerInterface interface {
	Check(scope *AccessScope, c *gin.Context) constants.ScopeAccessibility
}

type FUPCheckerInterface interface {
	Check(fup *FUPScope, c *gin.Context) (bool, error)
	Log(c *gin.Context) error
}

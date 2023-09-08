package contract

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
)

type AccessScopeCheckerInterface interface {
	Check(scope *AccessScope, c *gin.Context) constants.ScopeAccessibility
}

// TODO: FUP checker interface

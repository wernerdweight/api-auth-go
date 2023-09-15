package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
)

func Register(r *gin.Engine) {
	r.POST("/authenticate", authenticateHandler)
	if config.ProviderInstance.IsUserRegistrationEnabled() {
		r.POST("/registration/request", registrationRequestHandler)
		r.POST("/registration/confirm/:token", registrationConfirmHandler)
		r.POST("/resetting/request", resettingRequestHandler)
		r.POST("/resetting/reset/:token", resettingResetHandler)
	}
}

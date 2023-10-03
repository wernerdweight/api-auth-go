package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/routes"
	"github.com/wernerdweight/api-auth-go/auth/security"
	"log"
	"net/http"
)

func Middleware(r *gin.Engine, c contract.Config) gin.HandlerFunc {
	log.Println("setting up api-auth middleware...")
	config.ProviderInstance.Init(c)
	routes.Register(r)

	if config.ProviderInstance.IsCacheEnabled() {
		log.Println("initializing cache driver...")
		config.ProviderInstance.GetCacheDriver().Init(
			config.ProviderInstance.GetCachePrefix(),
			config.ProviderInstance.GetCacheTTL(),
		)
	}

	if !config.ProviderInstance.IsApiKeyModeEnabled() && !config.ProviderInstance.IsClientIdAndSecretModeEnabled() {
		log.Println("api-auth is disabled")
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if config.ProviderInstance.ShouldExcludeOptionsRequests() && http.MethodOptions == c.Request.Method {
			c.Next()
			return
		}

		err := security.Authenticate(c)
		if nil != err {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    err.Code,
				"error":   err.Err.Error(),
				"payload": err.Payload,
			})
			return
		}

		c.Next()
	}
}

package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"github.com/wernerdweight/api-auth-go/v2/auth/routes"
	"github.com/wernerdweight/api-auth-go/v2/auth/security"
	"github.com/wernerdweight/events-go"
	"log"
	"net/http"
)

func Middleware(r *gin.Engine, c contract.Config) gin.HandlerFunc {
	log.Println("setting up api-auth middleware...")
	config.ProviderInstance.Init(c)
	defer func() {
		// routes need to be registered asynchronously after the middleware is applied, so that is gets applied to the routes as well
		go routes.Register(r)
	}()

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
			errorResponse := gin.H{
				"code":    err.Code,
				"message": err.Err.Error(),
				"payload": err.Payload,
			}
			c.AbortWithStatusJSON(err.Status, errorResponse)
			events.GetEventHub().DispatchAsync(&contract.AuthenticationFailedEvent{
				Error:    *err,
				Context:  c,
				Response: errorResponse,
			})
			return
		}

		c.Next()
	}
}

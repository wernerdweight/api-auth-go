package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"github.com/wernerdweight/api-auth-go/v2/auth/security"
	"github.com/wernerdweight/events-go"
	"log"
	"net/http"
)

// Middleware returns a gin.HandlerFunc that authenticates requests based on the provided config.
// Auth routes (e.g. /authenticate, /registration/*) are NOT registered automatically.
// You must call routes.Register(r) after r.Use(auth.Middleware(...)) to register them:
//
//	r.Use(auth.Middleware(r, cfg))
//	routes.Register(r)
func Middleware(r *gin.Engine, c contract.Config) gin.HandlerFunc {
	log.Println("setting up api-auth middleware...")
	config.ProviderInstance.Init(c)
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

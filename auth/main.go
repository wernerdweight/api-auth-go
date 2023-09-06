package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/routes"
	"github.com/wernerdweight/api-auth-go/auth/security"
	"log"
	"net/http"
	"time"
)

func Middleware(r *gin.Engine, config contract.Config) gin.HandlerFunc {
	log.Println("setting up api-auth middleware...")
	routes.Register(r)
	// TODO: init config provider

	return func(c *gin.Context) {
		t := time.Now()

		// before request
		err := security.Authenticate(c)
		if nil != err {
			log.Printf("AUTH: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Next()

		// after request
		latency := time.Since(t)
		log.Printf("AUTH: latency: %d", latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Printf("AUTH: status: %d", status)
		log.Printf("AUTH: value: %d", c.GetInt("example"))
	}
}

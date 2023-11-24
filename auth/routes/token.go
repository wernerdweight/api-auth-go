package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	generator "github.com/wernerdweight/token-generator-go"
	"net/http"
	"time"
)

func generateTokenHandler(c *gin.Context) {
	if !config.ProviderInstance.IsCacheEnabled() {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    contract.CacheDisabled,
			"error":   contract.AuthErrorCodes[contract.CacheDisabled],
			"payload": nil,
		})
		return
	}
	cacheDriver := config.ProviderInstance.GetCacheDriver()

	apiClient, ok := c.Get(constants.ApiClient)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    contract.Unauthorized,
			"error":   contract.AuthErrorCodes[contract.Unauthorized],
			"payload": nil,
		})
		return
	}

	tokenGenerator := generator.NewTokenGenerator("")
	token := contract.OneOffToken{
		Value:   tokenGenerator.Generate(constants.DefaultTokenLength),
		Expires: time.Now().Add(config.ProviderInstance.GetOneOffTokenExpirationInterval()),
	}

	cacheDriver.SetApiClientByOneOffToken(token, apiClient.(contract.ApiClientInterface))

	// do not return milliseconds in response
	rfc3339Output := map[string]string{
		"token": token.Value,
		"expires": token.Expires.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, rfc3339Output)
}

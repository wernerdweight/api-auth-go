package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	generator "github.com/wernerdweight/token-generator-go"
	"net/http"
	"time"
)

func generateTokenHandler(c *gin.Context) {
	if !config.ProviderInstance.IsCacheEnabled() {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    contract.CacheDisabled,
			"message": contract.AuthErrorCodes[contract.CacheDisabled],
			"payload": nil,
		})
		return
	}
	cacheDriver := config.ProviderInstance.GetCacheDriver()

	apiClient, ok := c.Get(constants.ApiClient)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    contract.Unauthorized,
			"message": contract.AuthErrorCodes[contract.Unauthorized],
			"payload": nil,
		})
		return
	}

	tokenGenerator := generator.NewTokenGenerator("")
	token := contract.OneOffToken{
		Value:   tokenGenerator.Generate(constants.OneOffTokenLength),
		Expires: time.Now().Add(config.ProviderInstance.GetOneOffTokenExpirationInterval()),
	}

	cacheDriver.SetApiClientByOneOffToken(token, apiClient.(contract.ApiClientInterface))

	c.JSON(http.StatusOK, token)
}

package security

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"log"
	"regexp"
)

func shouldAuthenticate(c *gin.Context) bool {
	log.Printf("handler: %s", c.Request.URL.String())
	if targetHandlers := config.ProviderInstance.GetTargetHandlers(); targetHandlers != nil {
		for _, targetHandler := range *targetHandlers {
			matched, err := regexp.MatchString(targetHandler, c.Request.URL.String())
			if nil != err {
				log.Printf("can't match target handler pattern '%s': %v", targetHandler, err)
			}
			if matched {
				return true
			}
		}
	}
	return false
}

func shouldAuthenticateByApiClientAndUser(c *gin.Context) bool {
	return c.Request.Header.Get(constants.ClientIdHeader) != "" && c.Request.Header.Get(constants.ClientSecretHeader) != ""
}

func shouldAuthenticateByApiKey(c *gin.Context) bool {
	return c.Request.Header.Get(constants.ApiKeyHeader) != ""
}

func Authenticate(c *gin.Context) *contract.AuthError {
	if !shouldAuthenticate(c) {
		return nil
	}

	if shouldAuthenticateByApiClientAndUser(c) {
		apiClientProvider := config.ProviderInstance.GetClientProvider()
		apiClient, err := apiClientProvider.ProvideByIdAndSecret(
			c.Request.Header.Get(constants.ClientIdHeader),
			c.Request.Header.Get(constants.ClientSecretHeader),
		)
		if nil != err {
			return err
		}
		c.Set(constants.ApiClient, apiClient)
		// TODO: if scope access model is enabled, check if client (and user) is allowed to access this endpoint
		if config.ProviderInstance.IsClientScopeAccessModelEnabled() {
			// TODO: check client scope (and FUP)
			// TODO: - if client is not allowed to access this endpoint (omitted or false), return respective error
			// TODO: - if client is allowed to access this endpoint freely (true), return nil (done)
			// TODO: - if client is allowed to access this endpoint on-behalf ("on-behalf"), continue
			// TODO: authenticate api user (user token)
			// TODO: check user scope (and FUP)
			// TODO: - if user is not allowed to access this endpoint (omitted or false), return respective error
			// TODO: - if user is allowed to access this endpoint freely (true), return nil (done)
		}
	}

	if shouldAuthenticateByApiKey(c) {
		// TODO: authenticate user by api key (use some default api client)
		// TODO: check scope access (and FUP)
	}
	return contract.NewAuthError(contract.NoCredentialsProvided, nil)
}

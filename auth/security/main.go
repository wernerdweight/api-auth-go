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

func authenticateApiClient(c *gin.Context) (contract.ApiClientInterface, *contract.AuthError) {
	apiClientProvider := config.ProviderInstance.GetClientProvider()
	if shouldAuthenticateByApiClientAndUser(c) {
		apiClient, err := apiClientProvider.ProvideByIdAndSecret(
			c.Request.Header.Get(constants.ClientIdHeader),
			c.Request.Header.Get(constants.ClientSecretHeader),
		)
		if nil != err {
			return nil, err
		}
		return apiClient, nil
	}
	if shouldAuthenticateByApiKey(c) {
		apiClient, err := apiClientProvider.ProvideByApiKey(c.Request.Header.Get(constants.ApiKeyHeader))
		if nil != err {
			return nil, err
		}
		return apiClient, nil
	}
	return nil, contract.NewAuthError(contract.NoCredentialsProvided, nil)
}

func authenticateApiUser(c *gin.Context) (contract.ApiUserInterface, *contract.AuthError) {
	if c.Request.Header.Get(constants.ApiUserTokenHeader) == "" {
		return nil, contract.NewAuthError(contract.UserTokenRequired, nil)
	}
	apiUserProvider := config.ProviderInstance.GetUserProvider()
	if nil == apiUserProvider {
		return nil, contract.NewAuthError(contract.UserProviderNotConfigured, nil)
	}
	apiUser, err := apiUserProvider.ProvideByToken(c.Request.Header.Get(constants.ApiUserTokenHeader))
	if nil != err {
		return nil, err
	}
	return apiUser, nil
}

func authenticateOnBehalf(c *gin.Context) *contract.AuthError {
	apiUser, err := authenticateApiUser(c)
	if nil != err {
		return err
	}

	c.Set(constants.ApiUser, apiUser)

	if !config.ProviderInstance.IsUserScopeAccessModelEnabled() {
		return nil
	}

	// TODO: check user FUP
	userAccessScopeChecker := config.ProviderInstance.GetUserScopeAccessChecker()
	userScopeAccessibility := userAccessScopeChecker.Check(apiUser.GetUserScope(), c)

	if constants.ScopeAccessibilityForbidden == userScopeAccessibility {
		return contract.NewAuthError(contract.UserForbidden, nil)
	}
	return nil
}

func Authenticate(c *gin.Context) *contract.AuthError {
	if !shouldAuthenticate(c) {
		return nil
	}

	apiClient, err := authenticateApiClient(c)
	if nil != err {
		return err
	}

	c.Set(constants.ApiClient, apiClient)

	if !config.ProviderInstance.IsClientScopeAccessModelEnabled() {
		return nil
	}

	// TODO: check client FUP
	clientAccessScopeChecker := config.ProviderInstance.GetClientScopeAccessChecker()
	scopeAccessibility := clientAccessScopeChecker.Check(apiClient.GetClientScope(), c)

	if constants.ScopeAccessibilityForbidden == scopeAccessibility {
		return contract.NewAuthError(contract.ClientForbidden, nil)
	}

	if constants.ScopeAccessibilityAccessible == scopeAccessibility {
		return nil
	}

	if constants.ScopeAccessibilityOnBehalf == scopeAccessibility {
		return authenticateOnBehalf(c)
	}

	return contract.NewAuthError(contract.UnknownScopeAccessibility, map[string]constants.ScopeAccessibility{"scope": scopeAccessibility})
}

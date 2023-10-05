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
		clientId := c.Request.Header.Get(constants.ClientIdHeader)
		clientSecret := c.Request.Header.Get(constants.ClientSecretHeader)
		if config.ProviderInstance.IsCacheEnabled() {
			apiClient, err := config.ProviderInstance.GetCacheDriver().GetApiClientByIdAndSecret(clientId, clientSecret)
			if nil != apiClient {
				return apiClient, nil
			}
			if nil != err {
				log.Printf("can't get api client from cache: %v", err)
			}
		}
		apiClient, err := apiClientProvider.ProvideByIdAndSecret(clientId, clientSecret)
		if nil != err {
			return nil, err
		}
		if config.ProviderInstance.IsCacheEnabled() {
			err = config.ProviderInstance.GetCacheDriver().SetApiClientByIdAndSecret(clientId, clientSecret, apiClient)
			if nil != err {
				log.Printf("can't set api client to cache: %v", err)
			}
		}
		return apiClient, nil
	}
	if shouldAuthenticateByApiKey(c) {
		apiKey := c.Request.Header.Get(constants.ApiKeyHeader)
		if config.ProviderInstance.IsCacheEnabled() {
			apiClient, err := config.ProviderInstance.GetCacheDriver().GetApiClientByApiKey(apiKey)
			if nil != apiClient {
				return apiClient, nil
			}
			if nil != err {
				log.Printf("can't get api client from cache: %v", err)
			}
		}
		apiClient, err := apiClientProvider.ProvideByApiKey(apiKey)
		if nil != err {
			return nil, err
		}
		if config.ProviderInstance.IsCacheEnabled() {
			err = config.ProviderInstance.GetCacheDriver().SetApiClientByApiKey(apiKey, apiClient)
			if nil != err {
				log.Printf("can't set api client to cache: %v", err)
			}
		}
		return apiClient, nil
	}
	return nil, contract.NewAuthError(contract.NoCredentialsProvided, nil)
}

func authenticateApiUser(c *gin.Context) (contract.ApiUserInterface, *contract.AuthError) {
	if c.Request.Header.Get(constants.ApiUserTokenHeader) == "" {
		return nil, contract.NewAuthError(contract.UserTokenRequired, nil)
	}
	apiToken := c.Request.Header.Get(constants.ApiUserTokenHeader)
	if config.ProviderInstance.IsCacheEnabled() {
		apiUser, err := config.ProviderInstance.GetCacheDriver().GetApiUserByToken(apiToken)
		if nil != apiUser {
			return apiUser, nil
		}
		if nil != err {
			log.Printf("can't get api user from cache: %v", err)
		}
	}
	apiUserProvider := config.ProviderInstance.GetUserProvider()
	if nil == apiUserProvider {
		return nil, contract.NewInternalError(contract.UserProviderNotConfigured, nil)
	}
	apiUser, err := apiUserProvider.ProvideByToken(apiToken)
	if nil != err {
		return nil, err
	}
	if config.ProviderInstance.IsCacheEnabled() {
		err = config.ProviderInstance.GetCacheDriver().SetApiUserByToken(apiToken, apiUser)
		if nil != err {
			log.Printf("can't set api user to cache: %v", err)
		}
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

	if config.ProviderInstance.IsUserFUPEnabled() {
		userFUPChecker := config.ProviderInstance.GetUserFUPChecker()
		fupLimits := userFUPChecker.Check(apiUser.GetFUPScope(), c, apiUser.GetLogin())
		if nil != fupLimits.Error {
			return fupLimits.Error
		}
		if fupLimits.Accessible == constants.ScopeAccessibilityForbidden {
			return contract.NewFUPError(contract.RequestLimitDepleted, fupLimits.Limits)
		}
		header := fupLimits.GetLimitsHeader()
		if "" != header {
			c.Header(constants.UserFUPLimitsHeader, header)
		}
	}
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

	if config.ProviderInstance.IsClientFUPEnabled() {
		clientFUPChecker := config.ProviderInstance.GetClientFUPChecker()
		fupLimits := clientFUPChecker.Check(apiClient.GetFUPScope(), c, apiClient.GetClientId())
		if nil != fupLimits.Error {
			return fupLimits.Error
		}
		if fupLimits.Accessible == constants.ScopeAccessibilityForbidden {
			return contract.NewFUPError(contract.RequestLimitDepleted, fupLimits.Limits)
		}
		header := fupLimits.GetLimitsHeader()
		if "" != header {
			c.Header(constants.ClientFUPLimitsHeader, header)
		}
	}
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

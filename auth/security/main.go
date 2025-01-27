package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"log"
	"regexp"
)

func shouldAuthenticate(c *gin.Context) bool {
	excludeHandlers := config.ProviderInstance.GetExcludeHandlers()
	if nil != excludeHandlers && len(*excludeHandlers) > 0 {
		for _, excludeHandler := range *excludeHandlers {
			matched, err := regexp.MatchString(excludeHandler, c.Request.URL.String())
			if nil != err {
				log.Printf("can't match exclude handler pattern '%s': %v", excludeHandler, err)
			}
			if matched {
				return false
			}
		}
	}
	targetHandlers := config.ProviderInstance.GetTargetHandlers()
	if nil == targetHandlers || len(*targetHandlers) == 0 {
		return true
	}
	for _, targetHandler := range *targetHandlers {
		matched, err := regexp.MatchString(targetHandler, c.Request.URL.String())
		if nil != err {
			log.Printf("can't match target handler pattern '%s': %v", targetHandler, err)
		}
		if matched {
			return true
		}
	}
	return false
}

func shouldAuthenticateByOneOffToken(c *gin.Context) bool {
	return config.ProviderInstance.IsOneOffTokenModeEnabled() && c.Request.Header.Get(constants.OneOffTokenHeader) != ""
}

func shouldAuthenticateByApiClientAndSecret(c *gin.Context) bool {
	return c.Request.Header.Get(constants.ClientIdHeader) != "" && c.Request.Header.Get(constants.ClientSecretHeader) != ""
}

func shouldAuthenticateByApiKey(c *gin.Context) bool {
	return c.Request.Header.Get(constants.ApiKeyHeader) != ""
}

func authenticateApiClientByOneOffToken(c *gin.Context) (contract.ApiClientInterface, *contract.AuthError) {
	// check if one-off token is allowed for the current request
	targetHandlers := config.ProviderInstance.GetTargetOneOffTokenHandlers()
	if nil != targetHandlers && len(*targetHandlers) > 0 {
		inScope := false
		for _, targetHandler := range *targetHandlers {
			matched, err := regexp.MatchString(targetHandler, c.Request.URL.String())
			if nil != err {
				log.Printf("can't match one-off token target handler pattern '%s': %v", targetHandler, err)
			}
			if matched {
				inScope = true
				break
			}
		}
		if !inScope {
			return nil, contract.NewAuthError(contract.OneOffTokenNotAllowed, nil)
		}
	}

	token := c.Request.Header.Get(constants.OneOffTokenHeader)
	if !config.ProviderInstance.IsCacheEnabled() {
		return nil, contract.NewInternalError(contract.CacheDisabled, nil)
	}
	cacheDriver := config.ProviderInstance.GetCacheDriver()
	apiClient, err := cacheDriver.GetApiClientByOneOffToken(token)
	if nil != err {
		return nil, err
	}
	if nil == apiClient {
		return nil, contract.NewAuthError(contract.InvalidOneOffToken, nil)
	}
	err = cacheDriver.DeleteApiClientByOneOffToken(token)
	if nil != err {
		log.Printf("can't delete api client by one-off token: %v", err)
	}
	return apiClient, nil
}

func authenticateApiClientByApiClientAndSecret(c *gin.Context, apiClientProvider contract.ApiClientProviderInterface[contract.ApiClientInterface]) (contract.ApiClientInterface, *contract.AuthError) {
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

func authenticateApiClientByApiKey(c *gin.Context, apiClientProvider contract.ApiClientProviderInterface[contract.ApiClientInterface]) (contract.ApiClientInterface, *contract.AuthError) {
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

func authenticateApiClient(c *gin.Context) (contract.ApiClientInterface, *contract.AuthError) {
	if shouldAuthenticateByOneOffToken(c) {
		return authenticateApiClientByOneOffToken(c)
	}

	apiClientProvider := config.ProviderInstance.GetClientProvider()
	if shouldAuthenticateByApiClientAndSecret(c) {
		return authenticateApiClientByApiClientAndSecret(c, apiClientProvider)
	}

	if shouldAuthenticateByApiKey(c) {
		return authenticateApiClientByApiKey(c, apiClientProvider)
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
			c.Header(constants.RetryAfterHeader, fmt.Sprintf("%d", fupLimits.GetRetryAfter()))
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
		fupKey := apiClient.GetClientId()
		if apiClient.GetCurrentApiKey() != nil {
			// use API key as a part of the FUP key if additional API key is provided (it might have different limits)
			fupKey = fmt.Sprintf("%s:%s", fupKey, apiClient.GetCurrentApiKey().GetKey())
		}
		fupLimits := clientFUPChecker.Check(apiClient.GetFUPScope(), c, fupKey)
		if nil != fupLimits.Error {
			return fupLimits.Error
		}
		if fupLimits.Accessible == constants.ScopeAccessibilityForbidden {
			c.Header(constants.RetryAfterHeader, fmt.Sprintf("%d", fupLimits.GetRetryAfter()))
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

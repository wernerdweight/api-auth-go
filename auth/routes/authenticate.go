package routes

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"github.com/wernerdweight/api-auth-go/v2/auth/marshaller"
	"github.com/wernerdweight/events-go"
	generator "github.com/wernerdweight/token-generator-go"
	"net/http"
	"strings"
	"time"
)

func extractCredentials(header string) (string, string, *contract.AuthError) {
	encodedCredentials := header[len("Basic "):]
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if nil != err {
		return "", "", contract.NewAuthError(contract.InvalidCredentials, map[string]string{"details": err.Error()})
	}
	if !strings.Contains(string(decodedCredentials), ":") {
		return "", "", contract.NewAuthError(contract.InvalidCredentials, nil)
	}
	credentials := strings.Split(string(decodedCredentials), ":")
	return credentials[0], credentials[1], nil
}

func createToken() contract.ApiUserTokenInterface {
	tokenGenerator := generator.NewTokenGenerator("")
	token := tokenGenerator.Generate(constants.DefaultTokenLength)
	tokenClass := config.ProviderInstance.GetTokenFactory()()
	tokenClass.SetToken(token)
	tokenClass.SetExpirationDate(time.Now().Add(config.ProviderInstance.GetApiTokenExpirationInterval()))
	return tokenClass
}

func authenticateHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    contract.Unauthorized,
			"message": contract.AuthErrorCodes[contract.Unauthorized],
			"payload": nil,
		})
		return
	}

	apiClient, _ := c.Get(constants.ApiClient)
	var typedApiClient contract.ApiClientInterface
	if nil != apiClient {
		typedApiClient = apiClient.(contract.ApiClientInterface)
	}

	login, password, err := extractCredentials(authHeader)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    err.Code,
			"message": err.Err.Error(),
			"payload": err.Payload,
		})
	}

	apiUserProvider := config.ProviderInstance.GetUserProvider()
	apiUser, err := apiUserProvider.ProvideByLoginAndPassword(login, password)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    err.Code,
			"message": err.Err.Error(),
			"payload": err.Payload,
		})
		return
	}

	previousLoginAt := apiUser.GetLastLoginAt()
	token := createToken()
	now := time.Now()
	apiUser.AddApiToken(token)
	apiUser.SetLastLoginAt(&now)

	// authentication completed (issue an event for external handling)
	loginErr := events.GetEventHub().DispatchSync(&contract.AuthenticationCompletedEvent{
		ApiUser:   apiUser,
		ApiClient: typedApiClient,
	})
	if nil != loginErr {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    contract.Unauthorized,
			"message": contract.AuthErrorCodes[contract.Unauthorized],
			"payload": map[string]string{"details": loginErr.Error()},
		})
		return
	}

	err = apiUserProvider.Save(apiUser)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    err.Code,
			"message": err.Err.Error(),
			"payload": err.Payload,
		})
		return
	}

	// temporarily put previous login at back to return in the response (not to update it in the database)
	apiUser.SetLastLoginAt(previousLoginAt)
	output, err := marshaller.MarshalPublic(apiUser)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    err.Code,
			"message": err.Err.Error(),
			"payload": err.Payload,
		})
		return
	}
	c.JSON(200, output)
}

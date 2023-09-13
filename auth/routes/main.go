package routes

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"net/http"
	"strings"
	"time"
)

func Register(r *gin.Engine) {
	r.POST("/authenticate", authenticateHandler)
}

func extractCredentials(header string) (string, string, *contract.AuthError) {
	encodedCredentials := header[len("Basic "):]
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if nil != err {
		return "", "", contract.NewAuthError(contract.InvalidCredentials, map[string]error{"details": err})
	}
	if !strings.Contains(string(decodedCredentials), ":") {
		return "", "", contract.NewAuthError(contract.InvalidCredentials, nil)
	}
	credentials := strings.Split(string(decodedCredentials), ":")
	return credentials[0], credentials[1], nil
}

func generateToken(length int) string {
	// TODO: replace this with some tokenizer that supports a predefined alphabet
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	// FIXME: this is producing an ugly output
	return hex.EncodeToString(b)[0:length]
}

func createToken() contract.ApiUserTokenInterface {
	tokenClass := config.ProviderInstance.GetTokenFactory()()
	tokenClass.SetToken(generateToken(constants.DefaultTokenLength))
	tokenClass.SetExpirationDate(time.Now().Add(config.ProviderInstance.GetApiTokenExpirationInterval()))
	return tokenClass
}

func authenticateHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    contract.Unauthorized,
			"error":   contract.AuthErrorCodes[contract.Unauthorized],
			"payload": nil,
		})
		return
	}

	login, password, err := extractCredentials(authHeader)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    err.Code,
			"error":   err.Err.Error(),
			"payload": err.Payload,
		})
	}

	apiUserProvider := config.ProviderInstance.GetUserProvider()
	apiUser, err := apiUserProvider.ProvideByLoginAndPassword(login, password)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    err.Code,
			"error":   err.Err.Error(),
			"payload": err.Payload,
		})
		return
	}

	previousLoginAt := apiUser.GetLastLoginAt()
	token := createToken()
	now := time.Now()
	apiUser.AddApiToken(token)
	apiUser.SetLastLoginAt(&now)

	err = apiUserProvider.Save(apiUser)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    err.Code,
			"error":   err.Err.Error(),
			"payload": err.Payload,
		})
		return
	}

	// temporarily put previous login at back to return in the response (not to update it in the database)
	apiUser.SetLastLoginAt(previousLoginAt)
	c.JSON(200, apiUser)
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/encoder"
	"github.com/wernerdweight/events-go"
	generator "github.com/wernerdweight/token-generator-go"
	"net/http"
	"time"
)

type ResettingRequestRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResettingResetRequest struct {
	Password string `json:"password" binding:"required"`
}

func resettingRequestHandler(c *gin.Context) {
	request := ResettingRequestRequest{}
	if err := c.ShouldBindJSON(&request); nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": map[string]string{"details": err.Error()},
		})
		return
	}

	provider := config.ProviderInstance.GetUserProvider()
	apiUser, authErr := provider.ProvideByLogin(request.Email)
	if nil == apiUser {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    contract.UserNotFound,
			"error":   contract.AuthErrorCodes[contract.UserNotFound],
			"payload": nil,
		})
		return
	}
	if nil != authErr && contract.UserNotFound != authErr.Code {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    authErr.Code,
			"error":   authErr.Err.Error(),
			"payload": authErr.Payload,
		})
		return
	}
	if !apiUser.IsActive() {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    contract.UserNotActive,
			"error":   contract.AuthErrorCodes[contract.UserNotActive],
			"payload": nil,
		})
		return
	}

	// check for recent requests (prevent spam)
	if nil != apiUser.GetResetRequestedAt() && nil != apiUser.GetResetToken() {
		expirationInterval := config.ProviderInstance.GetConfirmationTokenExpirationInterval()
		expiresAt := apiUser.GetResetRequestedAt().Add(expirationInterval)
		if expiresAt.After(time.Now()) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"code":    contract.ResettingAlreadyRequested,
				"error":   contract.AuthErrorCodes[contract.ResettingAlreadyRequested],
				"payload": map[string]time.Time{"expiresAt": expiresAt},
			})
			return
		}
	}

	// generate reset token and set reset token date
	token := generator.NewTokenGenerator("").Generate(constants.DefaultTokenLength)
	now := time.Now()
	apiUser.SetResetToken(&token)
	apiUser.SetResetRequestedAt(&now)

	err := events.GetEventHub().DispatchSync(&contract.RequestResetApiUserPasswordEvent{
		ApiUser: apiUser,
	})
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": map[string]string{"details": err.Error()},
		})
		return
	}

	authErr = provider.Save(apiUser)
	if nil != authErr {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    authErr.Code,
			"error":   authErr.Err.Error(),
			"payload": authErr.Payload,
		})
		return
	}

	// call external service to send reset email (event)
	events.GetEventHub().DispatchAsync(&contract.ResettingRequestCompletedEvent{
		ApiUser: apiUser,
	})

	c.JSON(http.StatusAccepted, gin.H{
		"status": "ok",
	})
}

func resettingResetHandler(c *gin.Context) {
	token := c.Param("token")
	request := ResettingResetRequest{}
	if err := c.ShouldBindJSON(&request); nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": map[string]string{"details": err.Error()},
		})
		return
	}

	provider := config.ProviderInstance.GetUserProvider()
	apiUser, authErr := provider.ProvideByResetToken(token)
	if nil != authErr {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    authErr.Code,
			"error":   authErr.Err.Error(),
			"payload": authErr.Payload,
		})
		return
	}

	authErr = validatePassword(apiUser.GetLogin(), request.Password)
	if nil != authErr {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    authErr.Code,
			"error":   authErr.Err.Error(),
			"payload": map[string][]string{"details": authErr.Payload.([]string)},
		})
		return
	}

	encryptedPassword, err := encoder.EncryptPassword(request.Password)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    contract.EncryptionError,
			"error":   contract.AuthErrorCodes[contract.EncryptionError],
			"payload": map[string]string{"details": err.Error()},
		})
		return
	}

	apiUser.SetPassword(encryptedPassword)
	apiUser.SetResetToken(nil)
	apiUser.SetResetRequestedAt(nil)

	// call external service to set user details and other fields (event)
	err = events.GetEventHub().DispatchSync(&contract.ResetApiUserPasswordEvent{
		ApiUser: apiUser,
	})
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": map[string]string{"details": err.Error()},
		})
		return
	}

	authErr = provider.Save(apiUser)
	if nil != authErr {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    authErr.Code,
			"error":   authErr.Err.Error(),
			"payload": authErr.Payload,
		})
		return
	}

	// issue resetting done event
	events.GetEventHub().DispatchAsync(&contract.ResettingCompletedEvent{
		ApiUser: apiUser,
	})

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

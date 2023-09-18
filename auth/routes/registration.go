package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/events-go"
	"net/http"
)

type RegistrationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func registrationRequestHandler(c *gin.Context) {
	request := RegistrationRequest{}
	// TODO: get login information (email) from request body
	if err := c.ShouldBindJSON(&request); nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": err.Error(),
		})
		return
	}
	// TODO: validate login information (basics, not empty, not same as password, etc.)
	// TODO: validate login information (issue an event for external validation)
	err := events.GetEventHub().DispatchSync(&contract.ValidateLoginInformationEvent{
		Login:    request.Email,
		Password: request.Password,
	})
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": err.Error(),
		})
		return
	}
	// TODO: check for duplicates
	provider := config.ProviderInstance.GetUserProvider()
	user, err := provider.ProvideByLogin(request.Email)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    err.Code,
			"error":   err.Err,
			"payload": err.Payload,
		})
		return
	}
	if nil != user {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"code":    contract.UserAlreadyExists,
			"error":   contract.AuthErrorCodes[contract.UserAlreadyExists],
			"payload": nil,
		})
		return
	}
	// TODO: get password from request body
	// TODO: validate password
	err := validatePassword(request.Password) // TODO: move this to a separate package validator?
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": err.Error(),
		})
		return
	}
	// TODO: create new user (make sure to set active to false, set confirmation token)
	// TODO: call external service to set user details and other fields (event)
	// TODO: save user
	// TODO: call external registration service to handle confirm callback (event)
}

func registrationConfirmHandler(c *gin.Context) {
	// TODO: get token from request params
	// TODO: validate token
	// TODO: activate user
	// TODO: call external service to set user details and other fields (event)
	// TODO: save user
	// TODO: issue registration confirmed event
}

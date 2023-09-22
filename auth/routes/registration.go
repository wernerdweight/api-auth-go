package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/encoder"
	"github.com/wernerdweight/events-go"
	"net/http"
	"regexp"
)

type RegistrationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func validatePassword(email string, password string) *contract.AuthError {
	var payload []string
	if email == password {
		payload = append(payload, "password must not be the same as email")
	}
	if len(password) < 8 {
		payload = append(payload, "password must be at least 8 characters long")
	}
	if password == "password" {
		payload = append(payload, "password must not be 'password'")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		payload = append(payload, "password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		payload = append(payload, "password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`\d`).MatchString(password) {
		payload = append(payload, "password must contain at least one number")
	}
	if len(payload) > 0 {
		return contract.NewAuthError(contract.InvalidRequest, payload)
	}
	return nil
}

func registrationRequestHandler(c *gin.Context) {
	request := RegistrationRequest{}
	if err := c.ShouldBindJSON(&request); nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": map[string]string{"details": err.Error()},
		})
		return
	}

	// validate login information (issue an event for external validation)
	err := events.GetEventHub().DispatchSync(&contract.ValidateLoginInformationEvent{
		Login:    request.Email,
		Password: request.Password,
	})
	if nil != err {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    contract.InvalidRequest,
			"error":   contract.AuthErrorCodes[contract.InvalidRequest],
			"payload": map[string]string{"details": err.Error()},
		})
		return
	}

	// check for duplicates
	provider := config.ProviderInstance.GetUserProvider()
	user, authErr := provider.ProvideByLogin(request.Email)
	if nil != user {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"code":    contract.UserAlreadyExists,
			"error":   contract.AuthErrorCodes[contract.UserAlreadyExists],
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

	authErr = validatePassword(request.Email, request.Password)
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

	apiUser := provider.ProvideNew(request.Email, encryptedPassword)
	// call external service to set user details and other fields (event)
	err = events.GetEventHub().DispatchSync(&contract.CreateNewApiUserEvent{
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

	// call external registration service to handle confirm callback (event)
	events.GetEventHub().DispatchAsync(&contract.RegistrationRequestCompletedEvent{
		ApiUser: apiUser,
	})

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

func registrationConfirmHandler(c *gin.Context) {
	// TODO: get token from request params
	// TODO: validate token
	// TODO: activate user
	// TODO: call external service to set user details and other fields (event)
	// TODO: save user
	// TODO: issue registration confirmed event
}

package contract

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/events-go"
)

const (
	ValidateLoginInformationEventKey          = "api-auth-go.validate-login-information"
	CreateNewApiUserEventKey                  = "api-auth-go.create-new-api-user"
	RegistrationRequestCompletedEventKey      = "api-auth-go.registration-request-completed"
	ActivateApiUserEventKey                   = "api-auth-go.activate-api-user"
	RegistrationConfirmationCompletedEventKey = "api-auth-go.registration-confirmation-completed"
	RequestResetApiUserPasswordEventKey       = "api-auth-go.request-reset-api-user-password"
	ResetApiUserPasswordEventKey              = "api-auth-go.reset-api-user-password"
	ResettingRequestCompletedEventKey         = "api-auth-go.resetting-request-completed"
	ResettingCompletedEventKey                = "api-auth-go.resetting-completed"
	AuthenticationFailedEventKey              = "api-auth-go.authentication-failed"
)

type ValidateLoginInformationEvent struct {
	Login    string
	Password string
}

func (event *ValidateLoginInformationEvent) GetKey() events.EventKey {
	return ValidateLoginInformationEventKey
}

func (event *ValidateLoginInformationEvent) GetPayload() events.EventPayload {
	return event
}

type CreateNewApiUserEvent struct {
	ApiUser       ApiUserInterface
	ApiClient     ApiClientInterface
	PlainPassword string
}

func (event *CreateNewApiUserEvent) GetKey() events.EventKey {
	return CreateNewApiUserEventKey
}

func (event *CreateNewApiUserEvent) GetPayload() events.EventPayload {
	return event
}

type RegistrationRequestCompletedEvent struct {
	ApiUser   ApiUserInterface
	ApiClient ApiClientInterface
}

func (event *RegistrationRequestCompletedEvent) GetKey() events.EventKey {
	return RegistrationRequestCompletedEventKey
}

func (event *RegistrationRequestCompletedEvent) GetPayload() events.EventPayload {
	return event
}

type ActivateApiUserEvent struct {
	ApiUser   ApiUserInterface
	ApiClient ApiClientInterface
}

func (event *ActivateApiUserEvent) GetKey() events.EventKey {
	return ActivateApiUserEventKey
}

func (event *ActivateApiUserEvent) GetPayload() events.EventPayload {
	return event
}

type RegistrationConfirmationCompletedEvent struct {
	ApiUser   ApiUserInterface
	ApiClient ApiClientInterface
}

func (event *RegistrationConfirmationCompletedEvent) GetKey() events.EventKey {
	return RegistrationConfirmationCompletedEventKey
}

func (event *RegistrationConfirmationCompletedEvent) GetPayload() events.EventPayload {
	return event
}

type RequestResetApiUserPasswordEvent struct {
	ApiUser   ApiUserInterface
	ApiClient ApiClientInterface
}

func (event *RequestResetApiUserPasswordEvent) GetKey() events.EventKey {
	return RequestResetApiUserPasswordEventKey
}

func (event *RequestResetApiUserPasswordEvent) GetPayload() events.EventPayload {
	return event
}

type ResetApiUserPasswordEvent struct {
	ApiUser   ApiUserInterface
	ApiClient ApiClientInterface
}

func (event *ResetApiUserPasswordEvent) GetKey() events.EventKey {
	return ResetApiUserPasswordEventKey
}

func (event *ResetApiUserPasswordEvent) GetPayload() events.EventPayload {
	return event
}

type ResettingRequestCompletedEvent struct {
	ApiUser   ApiUserInterface
	ApiClient ApiClientInterface
}

func (event *ResettingRequestCompletedEvent) GetKey() events.EventKey {
	return ResettingRequestCompletedEventKey
}

func (event *ResettingRequestCompletedEvent) GetPayload() events.EventPayload {
	return event
}

type ResettingCompletedEvent struct {
	ApiUser   ApiUserInterface
	ApiClient ApiClientInterface
}

func (event *ResettingCompletedEvent) GetKey() events.EventKey {
	return ResettingCompletedEventKey
}

func (event *ResettingCompletedEvent) GetPayload() events.EventPayload {
	return event
}

type AuthenticationFailedEvent struct {
	Error    AuthError
	Context  *gin.Context
	Response gin.H
}

func (event *AuthenticationFailedEvent) GetKey() events.EventKey {
	return AuthenticationFailedEventKey
}

func (event *AuthenticationFailedEvent) GetPayload() events.EventPayload {
	return event
}

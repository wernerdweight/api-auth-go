package contract

import (
	"github.com/wernerdweight/events-go"
)

const (
	ValidateLoginInformationEventKey     = "api-auth-go.validate-login-information"
	CreateNewApiUserEventKey             = "api-auth-go.create-new-api-user"
	RegistrationRequestCompletedEventKey = "api-auth-go.registration-request-completed"
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
	ApiUser ApiUserInterface
}

func (event *CreateNewApiUserEvent) GetKey() events.EventKey {
	return CreateNewApiUserEventKey
}

func (event *CreateNewApiUserEvent) GetPayload() events.EventPayload {
	return event
}

type RegistrationRequestCompletedEvent struct {
	ApiUser ApiUserInterface
}

func (event *RegistrationRequestCompletedEvent) GetKey() events.EventKey {
	return RegistrationRequestCompletedEventKey
}

func (event *RegistrationRequestCompletedEvent) GetPayload() events.EventPayload {
	return event
}

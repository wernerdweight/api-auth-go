package contract

import (
	"github.com/wernerdweight/events-go"
)

const (
	ValidateLoginInformationEventKey = "validate_login_information"
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

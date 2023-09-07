package contract

import "errors"

type AuthErrorCode int

type AuthError struct {
	Err     error
	Code    AuthErrorCode
	Payload *interface{}
}

const (
	Unknown AuthErrorCode = iota
	ClientNotFound
	UserNotFound
	NoCredentialsProvided
)

var AuthErrorCodes = map[AuthErrorCode]string{
	Unknown:               "unknown error",
	ClientNotFound:        "client not found",
	UserNotFound:          "user not found",
	NoCredentialsProvided: "no credentials provided",
}

func NewAuthError(code AuthErrorCode, payload *interface{}) *AuthError {
	return &AuthError{
		Err:     errors.New(AuthErrorCodes[code]),
		Code:    code,
		Payload: payload,
	}
}

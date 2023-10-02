package contract

import "errors"

type AuthErrorCode int

type AuthError struct {
	Err     error
	Code    AuthErrorCode
	Payload interface{}
}

const (
	Unknown AuthErrorCode = iota
	Unauthorized
	ClientNotFound
	UserNotFound
	NoCredentialsProvided
	UserTokenRequired
	UserTokenNotFound
	UserTokenExpired
	ClientForbidden
	UserForbidden
	UnknownScopeAccessibility
	UserProviderNotConfigured
	DatabaseError
	InvalidCredentials
	InvalidRequest
	UserAlreadyExists
	EncryptionError
	UserNotActive
	ConfirmationTokenExpired
	ResettingAlreadyRequested
	ResetTokenExpired
	CacheError
	MarshallingError
)

var AuthErrorCodes = map[AuthErrorCode]string{
	Unknown:                   "unknown error",
	Unauthorized:              "unauthorized",
	ClientNotFound:            "client not found",
	UserNotFound:              "user not found",
	NoCredentialsProvided:     "no credentials provided",
	UserTokenRequired:         "user token required but not provided",
	UserTokenNotFound:         "user token not found",
	UserTokenExpired:          "user token expired",
	ClientForbidden:           "client access forbidden",
	UserForbidden:             "user access forbidden",
	UnknownScopeAccessibility: "unknown scope accessibility",
	UserProviderNotConfigured: "user provider not configured",
	DatabaseError:             "database error",
	InvalidCredentials:        "invalid credentials",
	InvalidRequest:            "invalid request",
	UserAlreadyExists:         "user already exists",
	EncryptionError:           "encryption error",
	UserNotActive:             "user not active",
	ConfirmationTokenExpired:  "confirmation token expired",
	ResettingAlreadyRequested: "resetting already requested",
	ResetTokenExpired:         "reset token expired",
	CacheError:                "cache error",
	MarshallingError:          "marshalling error",
}

func NewAuthError(code AuthErrorCode, payload interface{}) *AuthError {
	return &AuthError{
		Err:     errors.New(AuthErrorCodes[code]),
		Code:    code,
		Payload: payload,
	}
}

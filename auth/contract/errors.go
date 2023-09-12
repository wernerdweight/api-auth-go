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
	ClientForbidden
	UserForbidden
	UnknownScopeAccessibility
	UserProviderNotConfigured
	DatabaseError
	UserTokenExpired
	InvalidCredentials
)

var AuthErrorCodes = map[AuthErrorCode]string{
	Unknown:                   "unknown error",
	Unauthorized:              "unauthorized",
	ClientNotFound:            "client not found",
	UserNotFound:              "user not found",
	NoCredentialsProvided:     "no credentials provided",
	UserTokenRequired:         "user token required but not provided",
	ClientForbidden:           "client access forbidden",
	UserForbidden:             "user access forbidden",
	UnknownScopeAccessibility: "unknown scope accessibility",
	UserProviderNotConfigured: "user provider not configured",
	DatabaseError:             "database error",
	UserTokenExpired:          "user token expired",
	InvalidCredentials:        "invalid credentials",
}

func NewAuthError(code AuthErrorCode, payload interface{}) *AuthError {
	return &AuthError{
		Err:     errors.New(AuthErrorCodes[code]),
		Code:    code,
		Payload: payload,
	}
}

package contract

import "time"

type AccessScope interface{}

type ApiClientInterface interface {
	GetClientId() string
	GetClientSecret() string
	GetApiKey() string
	GetClientScope() AccessScope
}
type ApiUserInterface interface {
	AddApiToken(apiToken ApiUserTokenInterface)
	GetCurrentToken() *ApiUserTokenInterface
	GetUserScope() AccessScope
	GetLastLoginAt() time.Time
	SetLastLoginAt(lastLoginAt time.Time)
}
type ApiUserTokenInterface interface {
	SetToken(token string)
	GetToken() string
	SetExpirationDate(expirationDate time.Time)
	GetExpirationDate() time.Time
	SetApiUser(apiUser ApiUserInterface)
	GetApiUser() ApiUserInterface
}

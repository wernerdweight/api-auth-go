package entity

import (
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"time"
)

// GormApiClient is a struct that implements ApiClientInterface for GORM
type GormApiClient struct {
	ClientId     string
	ClientSecret string
	ApiKey       string
	AccessScope  *contract.AccessScope
}

func (c GormApiClient) GetClientId() string {
	return c.ClientId
}

func (c GormApiClient) GetClientSecret() string {
	return c.ClientSecret
}

func (c GormApiClient) GetApiKey() string {
	return c.ApiKey
}

func (c GormApiClient) GetClientScope() *contract.AccessScope {
	return c.AccessScope
}

// GormApiUser is a struct that implements ApiUserInterface for GORM
type GormApiUser struct {
	Login        string
	Password     string
	AccessScope  *contract.AccessScope
	LastLoginAt  time.Time
	CurrentToken *GormApiUserToken
	ApiTokens    []GormApiUserToken
}

func (u GormApiUser) AddApiToken(apiToken contract.ApiUserTokenInterface) {
	u.ApiTokens = append(u.ApiTokens, apiToken.(GormApiUserToken))
	u.CurrentToken = apiToken.(*GormApiUserToken)
}

func (u GormApiUser) GetCurrentToken() contract.ApiUserTokenInterface {
	return u.CurrentToken
}

func (u GormApiUser) GetUserScope() *contract.AccessScope {
	return u.AccessScope
}

func (u GormApiUser) GetLastLoginAt() time.Time {
	return u.LastLoginAt
}

func (u GormApiUser) SetLastLoginAt(lastLoginAt time.Time) {
	u.LastLoginAt = lastLoginAt
}

// GormApiUserToken is a struct that implements ApiUserTokenInterface for GORM
type GormApiUserToken struct {
	Token          string
	ExpirationDate time.Time
	ApiUser        GormApiUser
}

func (t GormApiUserToken) SetToken(token string) {
	t.Token = token
}

func (t GormApiUserToken) GetToken() string {
	return t.Token
}

func (t GormApiUserToken) SetExpirationDate(expirationDate time.Time) {
	t.ExpirationDate = expirationDate
}

func (t GormApiUserToken) GetExpirationDate() time.Time {
	return t.ExpirationDate
}

func (t GormApiUserToken) SetApiUser(apiUser contract.ApiUserInterface) {
	t.ApiUser = apiUser.(GormApiUser)
}

func (t GormApiUserToken) GetApiUser() contract.ApiUserInterface {
	return t.ApiUser
}

package entity

import (
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"time"
)

// MemoryApiClient is the simplest struct that implements ApiClientInterface
type MemoryApiClient struct {
	Id          string
	Secret      string
	ApiKey      string
	AccessScope *contract.AccessScope
}

func (c *MemoryApiClient) GetClientId() string {
	return c.Id
}

func (c *MemoryApiClient) GetClientSecret() string {
	return c.Secret
}

func (c *MemoryApiClient) GetApiKey() string {
	return c.ApiKey
}

func (c *MemoryApiClient) GetClientScope() *contract.AccessScope {
	return c.AccessScope
}

// MemoryApiUser is the simplest struct that implements ApiUserInterface
type MemoryApiUser struct {
	Id                string
	Login             string
	Password          string
	CurrentToken      *MemoryApiUserToken
	AccessScope       *contract.AccessScope
	ConfirmationToken string
}

func (u *MemoryApiUser) AddApiToken(apiToken contract.ApiUserTokenInterface) {
	u.CurrentToken = apiToken.(*MemoryApiUserToken)
}

func (u *MemoryApiUser) GetCurrentToken() contract.ApiUserTokenInterface {
	return u.CurrentToken
}

func (u *MemoryApiUser) GetUserScope() *contract.AccessScope {
	return u.AccessScope
}

func (u *MemoryApiUser) GetLastLoginAt() *time.Time {
	lastLoginAt := time.Now()
	return &lastLoginAt
}

func (u *MemoryApiUser) SetLastLoginAt(lastLoginAt *time.Time) {
	// no-op
}

func (u *MemoryApiUser) GetPassword() string {
	return u.Password
}

func (u *MemoryApiUser) SetPassword(password string) {
	u.Password = password
}

func (u *MemoryApiUser) SetLogin(login string) {
	u.Login = login
}

func (u *MemoryApiUser) SetConfirmationToken(confirmationToken *string) {
	// no-op
}

func (u *MemoryApiUser) GetConfirmationRequestedAt() *time.Time {
	confirmationRequestedAt := time.Now()
	return &confirmationRequestedAt
}

func (u *MemoryApiUser) SetConfirmationRequestedAt(confirmationRequestedAt *time.Time) {
	// no-op
}

func (u *MemoryApiUser) IsActive() bool {
	return true
}

func (u *MemoryApiUser) SetActive(active bool) {
	// no-op
}

// MemoryApiUserToken is the simplest struct that implements ApiUserTokenInterface
type MemoryApiUserToken struct {
	Token          string
	ExpirationDate time.Time
	ApiUser        *MemoryApiUser
}

func (t *MemoryApiUserToken) SetToken(token string) {
	t.Token = token
}

func (t *MemoryApiUserToken) GetToken() string {
	return t.Token
}

func (t *MemoryApiUserToken) SetExpirationDate(expirationDate time.Time) {
	t.ExpirationDate = expirationDate
}

func (t *MemoryApiUserToken) GetExpirationDate() time.Time {
	return t.ExpirationDate
}

func (t *MemoryApiUserToken) SetApiUser(apiUser contract.ApiUserInterface) {
	t.ApiUser = apiUser.(*MemoryApiUser)
}

func (t *MemoryApiUserToken) GetApiUser() contract.ApiUserInterface {
	return t.ApiUser
}

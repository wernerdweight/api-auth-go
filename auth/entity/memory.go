package entity

import (
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"time"
)

// MemoryApiClient is the simplest struct that implements ApiClientInterface
type MemoryApiClient struct {
	Id          string                `json:"clientId" groups:"internal"`
	Secret      string                `json:"clientSecret" groups:"internal"`
	ApiKey      string                `json:"apiKey" groups:"internal"`
	AccessScope *contract.AccessScope `json:"clientScope" groups:"internal,public"`
	FUPScope    *contract.FUPScope    `json:"fupConfig" groups:"internal"`
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

func (c *MemoryApiClient) GetFUPScope() *contract.FUPScope {
	return c.FUPScope
}

// MemoryApiUser is the simplest struct that implements ApiUserInterface
type MemoryApiUser struct {
	Id                string                `json:"id" groups:"internal,public"`
	Login             string                `json:"login" groups:"internal"`
	Password          string                `json:"password" groups:"internal"`
	CurrentToken      *MemoryApiUserToken   `json:"token" groups:"internal,public"`
	AccessScope       *contract.AccessScope `json:"userScope" groups:"internal,public"`
	ConfirmationToken string                `json:"confirmationToken" groups:"internal"`
	ResetToken        string                `json:"resetToken" groups:"internal"`
	FUPScope          *contract.FUPScope    `json:"fupConfig" groups:"internal"`
}

func (u *MemoryApiUser) AddApiToken(apiToken contract.ApiUserTokenInterface) {
	memoryApiToken := MemoryApiUserToken{
		Token:          apiToken.GetToken(),
		ExpirationDate: apiToken.GetExpirationDate(),
	}
	u.CurrentToken = &memoryApiToken
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

func (u *MemoryApiUser) GetLogin() string {
	return u.Login
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

func (u *MemoryApiUser) GetResetRequestedAt() *time.Time {
	resetRequestedAt := time.Now()
	return &resetRequestedAt
}

func (u *MemoryApiUser) SetResetRequestedAt(resetRequestedAt *time.Time) {
	// no-op
}

func (u *MemoryApiUser) GetResetToken() *string {
	return nil
}

func (u *MemoryApiUser) SetResetToken(resetToken *string) {
	// no-op
}

func (u *MemoryApiUser) GetFUPScope() *contract.FUPScope {
	return u.FUPScope
}

// MemoryApiUserToken is the simplest struct that implements ApiUserTokenInterface
type MemoryApiUserToken struct {
	Token          string         `json:"token" groups:"internal,public"`
	ExpirationDate time.Time      `json:"expirationDate" groups:"internal,public"`
	ApiUser        *MemoryApiUser `json:"-"`
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

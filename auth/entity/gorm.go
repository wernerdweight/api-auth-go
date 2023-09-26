package entity

import (
	"github.com/google/uuid"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"log"
	"time"
)

// GormApiClient is a struct that implements ApiClientInterface for GORM
type GormApiClient struct {
	ID           uuid.UUID             `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ClientId     string                `json:"-"`
	ClientSecret string                `json:"-"`
	ApiKey       string                `json:"-"`
	AccessScope  *contract.AccessScope `gorm:"type:jsonb;serializer:json" json:"clientScope"`
	CreatedAt    time.Time             `gorm:"not null;default:CURRENT_TIMESTAMP" json:"-"`
}

func (c *GormApiClient) TableName() string {
	return "api_client"
}

func (c *GormApiClient) GetClientId() string {
	return c.ClientId
}

func (c *GormApiClient) GetClientSecret() string {
	return c.ClientSecret
}

func (c *GormApiClient) GetApiKey() string {
	return c.ApiKey
}

func (c *GormApiClient) GetClientScope() *contract.AccessScope {
	return c.AccessScope
}

// GormApiUser is a struct that implements ApiUserInterface for GORM
type GormApiUser struct {
	ID                      uuid.UUID                      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Login                   string                         `gorm:"column:email" json:"-"`
	Password                string                         `json:"-"`
	AccessScope             *contract.AccessScope          `gorm:"type:jsonb;serializer:json" json:"userScope"`
	LastLoginAt             *time.Time                     `json:"lastLoginAt"`
	CurrentToken            contract.ApiUserTokenInterface `gorm:"-" json:"token"`
	ApiTokens               []GormApiUserToken             `gorm:"foreignKey:ApiUserID" json:"-"`
	CreatedAt               time.Time                      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"-"`
	Active                  bool                           `gorm:"not null;default:false" json:"-"`
	ConfirmationRequestedAt *time.Time                     `json:"-"`
	ConfirmationToken       *string                        `json:"-"`
	ResetRequestedAt        *time.Time                     `json:"-"`
	ResetToken              *string                        `json:"-"`
}

func (u *GormApiUser) TableName() string {
	return "api_user"
}

func (u *GormApiUser) AddApiToken(apiToken contract.ApiUserTokenInterface) {
	gormApiToken := GormApiUserToken{
		Token:          apiToken.GetToken(),
		ExpirationDate: apiToken.GetExpirationDate(),
	}
	u.CurrentToken = apiToken
	u.ApiTokens = append(u.ApiTokens, gormApiToken)
	log.Printf("Added token %s to user: tokens: %v", apiToken.GetToken(), u.ApiTokens)
}

func (u *GormApiUser) GetCurrentToken() contract.ApiUserTokenInterface {
	return u.CurrentToken
}

func (u *GormApiUser) GetUserScope() *contract.AccessScope {
	return u.AccessScope
}

func (u *GormApiUser) GetLastLoginAt() *time.Time {
	return u.LastLoginAt
}

func (u *GormApiUser) SetLastLoginAt(lastLoginAt *time.Time) {
	u.LastLoginAt = lastLoginAt
}

func (u *GormApiUser) GetPassword() string {
	return u.Password
}

func (u *GormApiUser) SetPassword(password string) {
	u.Password = password
}

func (u *GormApiUser) GetLogin() string {
	return u.Login
}

func (u *GormApiUser) SetLogin(login string) {
	u.Login = login
}

func (u *GormApiUser) SetConfirmationToken(confirmationToken *string) {
	u.ConfirmationToken = confirmationToken
}

func (u *GormApiUser) GetConfirmationRequestedAt() *time.Time {
	return u.ConfirmationRequestedAt
}

func (u *GormApiUser) SetConfirmationRequestedAt(confirmationRequestedAt *time.Time) {
	u.ConfirmationRequestedAt = confirmationRequestedAt
}

func (u *GormApiUser) IsActive() bool {
	return u.Active
}

func (u *GormApiUser) SetActive(active bool) {
	u.Active = active
}

func (u *GormApiUser) GetResetRequestedAt() *time.Time {
	return u.ResetRequestedAt
}

func (u *GormApiUser) SetResetRequestedAt(resetRequestedAt *time.Time) {
	u.ResetRequestedAt = resetRequestedAt
}

func (u *GormApiUser) GetResetToken() *string {
	return u.ResetToken
}

func (u *GormApiUser) SetResetToken(resetToken *string) {
	u.ResetToken = resetToken
}

// GormApiUserToken is a struct that implements ApiUserTokenInterface for GORM
type GormApiUserToken struct {
	ID             uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"-"`
	Token          string       `json:"token"`
	ExpirationDate time.Time    `json:"expirationDate"`
	ApiUser        *GormApiUser `json:"-"`
	ApiUserID      uuid.UUID    `json:"-"`
	CreatedAt      time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" json:"-"`
}

func (t *GormApiUserToken) TableName() string {
	return "api_user_token"
}

func (t *GormApiUserToken) SetToken(token string) {
	t.Token = token
}

func (t *GormApiUserToken) GetToken() string {
	return t.Token
}

func (t *GormApiUserToken) SetExpirationDate(expirationDate time.Time) {
	t.ExpirationDate = expirationDate
}

func (t *GormApiUserToken) GetExpirationDate() time.Time {
	return t.ExpirationDate
}

func (t *GormApiUserToken) SetApiUser(apiUser contract.ApiUserInterface) {
	t.ApiUser = apiUser.(*GormApiUser)
}

func (t *GormApiUserToken) GetApiUser() contract.ApiUserInterface {
	return t.ApiUser
}

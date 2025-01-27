package entity

import (
	"github.com/google/uuid"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"time"
)

// GormApiClient is a struct that implements ApiClientInterface for GORM
type GormApiClient struct {
	ID             uuid.UUID             `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id" groups:"internal,public,id"`
	ClientId       string                `json:"clientId" groups:"internal,credentials"`
	ClientSecret   string                `json:"clientSecret" groups:"internal,credentials"`
	ApiKey         string                `json:"apiKey" groups:"internal,credentials"`
	AdditionalKeys []GormApiClientKey    `gorm:"foreignKey:ApiClientID" json:"-"`
	CurrentKey     *GormApiClientKey     `gorm:"-" json:"currentKey" groups:"internal,credentials"`
	AccessScope    *contract.AccessScope `gorm:"type:jsonb;serializer:json" json:"clientScope" groups:"internal,public"`
	FUPScope       *contract.FUPScope    `gorm:"type:jsonb;serializer:json" json:"fupConfig" groups:"internal"`
	CreatedAt      time.Time             `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt" groups:"internal"`
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

func (c *GormApiClient) GetCurrentApiKey() contract.ApiClientKeyInterface {
	if c.CurrentKey != nil {
		return c.CurrentKey
	}
	return nil
}

func (c *GormApiClient) SetCurrentApiKey(key contract.ApiClientKeyInterface) {
	currentKey := &GormApiClientKey{
		Key:            key.GetKey(),
		ExpirationDate: key.GetExpirationDate(),
		AccessScope:    key.GetClientScope(),
		FUPScope:       key.GetFUPScope(),
	}
	c.CurrentKey = currentKey
	c.AccessScope = key.GetClientScope()
	c.FUPScope = key.GetFUPScope()
}

func (c *GormApiClient) GetClientScope() *contract.AccessScope {
	if c.GetCurrentApiKey() != nil {
		return c.GetCurrentApiKey().GetClientScope()
	}
	return c.AccessScope
}

func (c *GormApiClient) GetFUPScope() *contract.FUPScope {
	if c.GetCurrentApiKey() != nil {
		return c.GetCurrentApiKey().GetFUPScope()
	}
	return c.FUPScope
}

// GormApiClientKey is a struct that implements ApiClientKeyInterface for GORM
type GormApiClientKey struct {
	ID             uuid.UUID             `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id" groups:"internal,public,id"`
	Key            string                `gorm:"uniqueIndex;not null" json:"key" groups:"internal,credentials"`
	ExpirationDate *time.Time            `json:"expirationDate" groups:"internal,public"`
	ApiClient      *GormApiClient        `json:"-"`
	ApiClientID    uuid.UUID             `gorm:"not null" json:"apiClientId" groups:"internal"`
	CreatedAt      time.Time             `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt" groups:"internal"`
	AccessScope    *contract.AccessScope `gorm:"type:jsonb;serializer:json" json:"clientScope" groups:"internal,public"`
	FUPScope       *contract.FUPScope    `gorm:"type:jsonb;serializer:json" json:"fupConfig" groups:"internal"`
}

func (k *GormApiClientKey) TableName() string {
	return "api_client_key"
}

func (k *GormApiClientKey) GetKey() string {
	return k.Key
}

func (k *GormApiClientKey) GetClientScope() *contract.AccessScope {
	return k.AccessScope
}

func (k *GormApiClientKey) GetFUPScope() *contract.FUPScope {
	return k.FUPScope
}

func (k *GormApiClientKey) GetApiClient() contract.ApiClientInterface {
	return k.ApiClient
}

func (k *GormApiClientKey) GetExpirationDate() *time.Time {
	return k.ExpirationDate
}

// GormApiUser is a struct that implements ApiUserInterface for GORM
type GormApiUser struct {
	ID                      uuid.UUID                      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id" groups:"internal,public,id"`
	Login                   string                         `gorm:"column:email" json:"login" groups:"internal,credentials"`
	Password                string                         `json:"password" groups:"internal"`
	AccessScope             *contract.AccessScope          `gorm:"type:jsonb;serializer:json" json:"userScope" groups:"internal,public"`
	FUPScope                *contract.FUPScope             `gorm:"type:jsonb;serializer:json" json:"fupConfig" groups:"internal"`
	LastLoginAt             *time.Time                     `json:"lastLoginAt" groups:"internal,public"`
	CurrentToken            contract.ApiUserTokenInterface `gorm:"-" json:"token" groups:"internal,public,credentials"`
	ApiTokens               []GormApiUserToken             `gorm:"foreignKey:ApiUserID" json:"-"`
	CreatedAt               time.Time                      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt" groups:"internal"`
	Active                  bool                           `gorm:"not null;default:false" json:"active" groups:"internal"`
	ConfirmationRequestedAt *time.Time                     `json:"confirmationRequestedAt" groups:"internal"`
	ConfirmationToken       *string                        `json:"confirmationToken" groups:"internal"`
	ResetRequestedAt        *time.Time                     `json:"resetRequestedAt" groups:"internal"`
	ResetToken              *string                        `json:"resetToken" groups:"internal"`
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

func (u *GormApiUser) GetFUPScope() *contract.FUPScope {
	return u.FUPScope
}

// GormApiUserToken is a struct that implements ApiUserTokenInterface for GORM
type GormApiUserToken struct {
	ID             uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id" groups:"internal"`
	Token          string       `json:"token" groups:"internal,public"`
	ExpirationDate time.Time    `json:"expirationDate" groups:"internal,public"`
	ApiUser        *GormApiUser `json:"-"`
	ApiUserID      uuid.UUID    `json:"apiUserId" groups:"internal"`
	CreatedAt      time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt" groups:"internal"`
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

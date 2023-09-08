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

func (c MemoryApiClient) GetClientId() string {
	return c.Id
}

func (c MemoryApiClient) GetClientSecret() string {
	return c.Secret
}

func (c MemoryApiClient) GetApiKey() string {
	return c.ApiKey
}

func (c MemoryApiClient) GetClientScope() *contract.AccessScope {
	return c.AccessScope
}

// MemoryApiUser is the simplest struct that implements ApiUserInterface
type MemoryApiUser struct {
	Id          string
	Login       string
	Password    string
	AccessToken string
	AccessScope *contract.AccessScope
}

func (u MemoryApiUser) AddApiToken(apiToken contract.ApiUserTokenInterface) {
	u.AccessToken = apiToken.GetToken()
}

func (u MemoryApiUser) GetCurrentToken() *contract.ApiUserTokenInterface {
	// TODO: implement ApiUserTokenInterface and return an instance here
	return nil
}

func (u MemoryApiUser) GetUserScope() *contract.AccessScope {
	return u.AccessScope
}

func (u MemoryApiUser) GetLastLoginAt() time.Time {
	return time.Now()
}

func (u MemoryApiUser) SetLastLoginAt(lastLoginAt time.Time) {
	// no-op
}

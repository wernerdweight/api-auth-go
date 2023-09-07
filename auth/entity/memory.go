package entity

import "github.com/wernerdweight/api-auth-go/auth/contract"

// MemoryApiClient is the simplest struct that implements ApiClientInterface
type MemoryApiClient struct {
	Id          string
	Secret      string
	ApiKey      string
	AccessScope contract.AccessScope
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

func (c MemoryApiClient) GetClientScope() contract.AccessScope {
	return c.AccessScope
}

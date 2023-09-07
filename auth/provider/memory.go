package provider

import (
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/entity"
)

// MemoryApiClientProvider is the simplest implementation of the ApiClientProviderInterface
type MemoryApiClientProvider struct {
	memory []entity.MemoryApiClient
}

func (p MemoryApiClientProvider) ProvideByIdAndSecret(id string, secret string) (contract.ApiClientInterface, *contract.AuthError) {
	for _, client := range p.memory {
		if client.Id == id && client.Secret == secret {
			return &client, nil
		}
	}

	return nil, contract.NewAuthError(contract.ClientNotFound, nil)
}

func (p MemoryApiClientProvider) ProvideByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	for _, client := range p.memory {
		if client.ApiKey == apiKey {
			return &client, nil
		}
	}

	return nil, contract.NewAuthError(contract.ClientNotFound, nil)
}

func NewMemoryApiClientProvider(memory []entity.MemoryApiClient) *MemoryApiClientProvider {
	return &MemoryApiClientProvider{
		memory: memory,
	}
}

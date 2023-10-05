package contract

import "time"

type CacheDriverInterface interface {
	Init(prefix string, ttl time.Duration) *AuthError
	GetApiClientByIdAndSecret(id string, secret string) (ApiClientInterface, *AuthError)
	SetApiClientByIdAndSecret(id string, secret string, client ApiClientInterface) *AuthError
	GetApiClientByApiKey(apiKey string) (ApiClientInterface, *AuthError)
	SetApiClientByApiKey(apiKey string, client ApiClientInterface) *AuthError
	GetApiUserByToken(token string) (ApiUserInterface, *AuthError)
	SetApiUserByToken(token string, user ApiUserInterface) *AuthError
	GetFUPEntry(key string) (*FUPCacheEntry, *AuthError)
	SetFUPEntry(key string, entry *FUPCacheEntry) *AuthError
}

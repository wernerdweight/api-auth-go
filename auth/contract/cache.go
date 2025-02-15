package contract

import "time"

type CacheDriverInterface interface {
	Init(prefix string, ttl time.Duration) *AuthError
	GetApiClientByIdAndSecret(id string, secret string) (ApiClientInterface, *AuthError)
	SetApiClientByIdAndSecret(id string, secret string, client ApiClientInterface) *AuthError
	GetApiClientByApiKey(apiKey string) (ApiClientInterface, *AuthError)
	SetApiClientByApiKey(apiKey string, client ApiClientInterface) *AuthError
	GetApiClientByOneOffToken(token string) (ApiClientInterface, *AuthError)
	SetApiClientByOneOffToken(oneOffToken OneOffToken, client ApiClientInterface) *AuthError
	DeleteApiClientByOneOffToken(token string) *AuthError
	GetApiUserByToken(token string) (ApiUserInterface, *AuthError)
	SetApiUserByToken(token string, user ApiUserInterface) *AuthError
	GetFUPEntry(key string) (*FUPCacheEntry, *AuthError)
	SetFUPEntry(key string, entry *FUPCacheEntry) *AuthError
	InvalidateToken(token string) *AuthError
}

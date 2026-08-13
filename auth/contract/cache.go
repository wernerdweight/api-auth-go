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
	// IncrementFUPEntry increments the counters of the given FUP entry and returns the updated entry.
	// A counter is incremented if the stored timestamp belongs to the same period as the increment
	// (see constants.Period.GetTimestampBounds), and reset to 1 otherwise; the stored timestamp is
	// the time of the increment. An implementation has to:
	//   - do all of that atomically, since concurrent requests sharing a FUP key would otherwise
	//     overwrite each other's counters and the limit would not be enforced,
	//   - not reset a counter (and not move the stored timestamp backwards) for an increment that
	//     arrives out of order, i.e. one that decided its period before an increment that was
	//     stored first - it counts towards the newer period instead,
	//   - expire entries after constants.FUPEntryTTL of inactivity, otherwise the counters of
	//     one-off sources (per IP, per cookie) accumulate forever,
	//   - return an entry that shares no state with what is stored, since the caller reads it
	//     while other requests keep incrementing.
	IncrementFUPEntry(key string) (*FUPCacheEntry, *AuthError)
	InvalidateToken(token string) *AuthError
}

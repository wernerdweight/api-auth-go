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
	//     one-off sources (per IP, per cookie) accumulate forever - implement
	//     FUPTTLCacheDriverInterface as well to expire them as soon as the limited periods allow,
	//   - return an entry that shares no state with what is stored, since the caller reads it
	//     while other requests keep incrementing.
	IncrementFUPEntry(key string) (*FUPCacheEntry, *AuthError)
	InvalidateToken(token string) *AuthError
}

// FUPTTLCacheDriverInterface is an optional addition to CacheDriverInterface for drivers that can
// expire a FUP entry after a caller-supplied TTL. The FUP checkers derive it from the scope being
// enforced (FUPScope.GetEntryTTL), so an entry is kept only as long as the longest period the
// scope limits by needs it - a scope limiting per minute and per day keeps its entries for two
// days rather than the 35 the monthly period would need. That bounds the memory a FUP key with an
// unbounded key space (per IP, per cookie) can occupy.
//
// A driver that doesn't implement it is used through IncrementFUPEntry and keeps every entry for
// constants.FUPEntryTTL, which is correct, just less frugal. The bundled memory and Redis drivers
// implement it - a driver that wraps one of them (to add metrics, to inject failures) therefore
// has to override both increments, or the embedded implementation stays in use for this one.
type FUPTTLCacheDriverInterface interface {
	// IncrementFUPEntryWithTTL behaves exactly like CacheDriverInterface.IncrementFUPEntry, except
	// that the entry expires after the given TTL of inactivity instead of constants.FUPEntryTTL.
	IncrementFUPEntryWithTTL(key string, ttl time.Duration) (*FUPCacheEntry, *AuthError)
}

Changelog
====================================

v3.1.0
------------

### Added

- `contract.FUPTTLCacheDriverInterface`, an optional addition to `contract.CacheDriverInterface` with a single method, `IncrementFUPEntryWithTTL(key string, ttl time.Duration) (*FUPCacheEntry, *AuthError)`. The FUP checkers use it when the configured driver implements it, deriving the TTL from the scope being enforced (`contract.FUPScope.GetEntryTTL`), so an entry is kept only as long as the longest period the scope limits by needs it - two days for a scope limiting per minute and per day, instead of the 35 the monthly period needs. That bounds the memory a FUP key with an unbounded key space (per IP, per cookie) can occupy: at the 490 requests/minute a single unauthenticated source was measured at, the difference is roughly 5 GB and 300 MB of Redis.

  Nothing has to be implemented - this is not a breaking change. A driver that only implements `CacheDriverInterface` keeps being used through `IncrementFUPEntry`, with `constants.FUPEntryTTL` for every entry as before. **A custom driver that wraps a bundled one** (to add metrics, to inject failures) **has to override both increments**, though, or the embedded implementation stays in use for the TTL-aware one.

- `constants.Period.GetEntryTTL` returns how long an entry has to survive inactivity to keep counting that period correctly, and `contract.FUPScope.GetEntryTTL(path)` the longest such TTL a given scope path needs.

v3.0.0
------------

### Breaking changes

- `contract.CacheDriverInterface` has a new method, `IncrementFUPEntry(key string) (*FUPCacheEntry, *AuthError)`, which increments the counters of a FUP entry and returns the updated entry. It replaces the `GetFUPEntry` + `Increment` + `SetFUPEntry` cycle the FUP checkers used to do, which lost increments of requests that shared a FUP key and were handled concurrently.

  The bundled memory and Redis drivers implement it (the Redis one as a Lua script). **A custom cache driver has to implement it as well, and the increment has to be atomic** - otherwise concurrent requests overwrite each other's counters and more of them are let through than the limit allows. `GetFUPEntry` and `SetFUPEntry` are unchanged and remain part of the interface.

- The module path changed to `github.com/wernerdweight/api-auth-go/v3`.

The stored format of FUP entries did not change, so no data needs to be migrated and no counter is reset by the upgrade.

### Added

- Anonymous FUP limits: `Client.AnonymousFUPScope` limits requests that fail to authenticate (which resolve no client and were not limited at all before), counted under the FUP key `anonymous` and checked with `Client.AnonymousFUPChecker` - `fup.IPFUPChecker`, i.e. per IP address, unless configured otherwise. Optional; unauthenticated traffic is not limited unless the scope is set. See the README for the trade-offs (fail-open, bypassability, excluded handlers).

### Fixed

- FUP cache entries are now written with an expiration (35 days of inactivity, slightly more than the longest supported period) instead of never expiring.
- The lazy initialization of the Redis client is guarded, so concurrently handled requests can no longer each create their own client and connection pool.

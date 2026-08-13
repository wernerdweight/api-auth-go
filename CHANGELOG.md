Changelog
====================================

## v3.0.0

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

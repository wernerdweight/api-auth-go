package cache

import (
	"github.com/wernerdweight/api-auth-go/v3/auth/constants"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
	"sync"
	"time"
)

type MemoryCacheEntry[T any] struct {
	Value    T
	ExpireAt time.Time
}

// MemoryCacheDriver is the simplest implementation of the CacheDriverInterface
// Do not use this driver for multi-instance applications!
type MemoryCacheDriver struct {
	apiClientMemory map[string]MemoryCacheEntry[contract.ApiClientInterface]
	apiUserMemory   map[string]MemoryCacheEntry[contract.ApiUserInterface]
	fupMemory       map[string]MemoryCacheEntry[contract.FUPCacheEntry]
	fupLock         sync.Mutex
	prefix          string
	ttl             time.Duration
}

func (d *MemoryCacheDriver) Init(prefix string, ttl time.Duration) *contract.AuthError {
	d.prefix = prefix
	d.ttl = ttl
	return nil
}

func (d *MemoryCacheDriver) getPrefix(groupPrefix GroupType) string {
	return getPrefix(d.prefix, groupPrefix)
}

func (d *MemoryCacheDriver) GetApiClientByIdAndSecret(id string, secret string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + id + secret
	if hit, ok := d.apiClientMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByIdAndSecret(id string, secret string, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + id + secret
	d.apiClientMemory[key] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetApiClientByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + apiKey
	if hit, ok := d.apiClientMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByApiKey(apiKey string, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + apiKey
	d.apiClientMemory[key] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetApiClientByOneOffToken(token string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + "-one_off-" + token
	if hit, ok := d.apiClientMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByOneOffToken(oneOffToken contract.OneOffToken, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + "-one_off-" + oneOffToken.Value
	d.apiClientMemory[key] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: oneOffToken.Expires,
	}
	return nil
}

func (d *MemoryCacheDriver) DeleteApiClientByOneOffToken(token string) *contract.AuthError {
	delete(d.apiClientMemory, d.getPrefix(GroupTypeAuth)+"-one_off-"+token)
	return nil
}

func (d *MemoryCacheDriver) GetApiUserByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + token
	if hit, ok := d.apiUserMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiUserMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiUserByToken(token string, user contract.ApiUserInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + token
	d.apiUserMemory[key] = MemoryCacheEntry[contract.ApiUserInterface]{
		Value:    user,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

// detachFUPEntry returns a copy of the given entry that shares no state with it. Copying the entry
// itself is not enough - its Used map would still be the very same map. Entries cross the driver's
// lock in both directions (they are returned to the caller and stored on its behalf), so without
// detaching them the caller would read the map while another request increments it under the lock.
func detachFUPEntry(entry *contract.FUPCacheEntry) *contract.FUPCacheEntry {
	used := make(map[constants.Period]int, len(entry.Used))
	for period, value := range entry.Used {
		used[period] = value
	}
	return &contract.FUPCacheEntry{
		UpdatedAt: entry.UpdatedAt,
		Used:      used,
	}
}

func (d *MemoryCacheDriver) getFUPEntry(entryKey string) *contract.FUPCacheEntry {
	if hit, ok := d.fupMemory[entryKey]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return detachFUPEntry(&hit.Value)
		}
		delete(d.fupMemory, entryKey)
	}
	return &contract.FUPCacheEntry{
		UpdatedAt: time.Time{},
		Used: map[constants.Period]int{
			constants.PeriodMinutely: 0,
			constants.PeriodHourly:   0,
			constants.PeriodDaily:    0,
			constants.PeriodWeekly:   0,
			constants.PeriodMonthly:  0,
		},
	}
}

func (d *MemoryCacheDriver) setFUPEntry(entryKey string, entry *contract.FUPCacheEntry, ttl time.Duration) {
	d.fupMemory[entryKey] = MemoryCacheEntry[contract.FUPCacheEntry]{
		Value:    *detachFUPEntry(entry),
		ExpireAt: time.Now().Add(ttl),
	}
}

func (d *MemoryCacheDriver) GetFUPEntry(key string) (*contract.FUPCacheEntry, *contract.AuthError) {
	d.fupLock.Lock()
	defer d.fupLock.Unlock()
	return d.getFUPEntry(d.getPrefix(GroupTypeFUP) + key), nil
}

func (d *MemoryCacheDriver) SetFUPEntry(key string, entry *contract.FUPCacheEntry) *contract.AuthError {
	d.fupLock.Lock()
	defer d.fupLock.Unlock()
	d.setFUPEntry(d.getPrefix(GroupTypeFUP)+key, entry, constants.FUPEntryTTL)
	return nil
}

func (d *MemoryCacheDriver) IncrementFUPEntry(key string) (*contract.FUPCacheEntry, *contract.AuthError) {
	return d.IncrementFUPEntryWithTTL(key, constants.FUPEntryTTL)
}

func (d *MemoryCacheDriver) IncrementFUPEntryWithTTL(key string, ttl time.Duration) (*contract.FUPCacheEntry, *contract.AuthError) {
	d.fupLock.Lock()
	defer d.fupLock.Unlock()
	entryKey := d.getPrefix(GroupTypeFUP) + key
	entry := d.getFUPEntry(entryKey)
	entry.Increment()
	d.setFUPEntry(entryKey, entry, ttl)
	return entry, nil
}

func (d *MemoryCacheDriver) InvalidateToken(token string) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + token
	delete(d.apiUserMemory, key)
	return nil
}

func NewMemoryCacheDriver() *MemoryCacheDriver {
	return &MemoryCacheDriver{
		apiClientMemory: make(map[string]MemoryCacheEntry[contract.ApiClientInterface]),
		apiUserMemory:   make(map[string]MemoryCacheEntry[contract.ApiUserInterface]),
		fupMemory:       make(map[string]MemoryCacheEntry[contract.FUPCacheEntry]),
	}
}

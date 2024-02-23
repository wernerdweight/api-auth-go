package cache

import (
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
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
	prefix          string
	ttl             time.Duration
}

func (d *MemoryCacheDriver) Init(prefix string, ttl time.Duration) *contract.AuthError {
	d.prefix = prefix
	d.ttl = ttl
	return nil
}

func (d *MemoryCacheDriver) getPrefix(groupPrefix string) string {
	if len(groupPrefix) > 0 && groupPrefix[len(groupPrefix)-1] != '_' {
		groupPrefix += "_"
	}
	return d.prefix + groupPrefix
}

func (d *MemoryCacheDriver) GetApiClientByIdAndSecret(id string, secret string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix("auth") + id + secret
	if hit, ok := d.apiClientMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByIdAndSecret(id string, secret string, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix("auth") + id + secret
	d.apiClientMemory[key] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetApiClientByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix("auth") + apiKey
	if hit, ok := d.apiClientMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByApiKey(apiKey string, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix("auth") + apiKey
	d.apiClientMemory[key] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetApiClientByOneOffToken(token string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix("auth") + "-one_off-" + token
	if hit, ok := d.apiClientMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByOneOffToken(oneOffToken contract.OneOffToken, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix("auth") + "-one_off-" + oneOffToken.Value
	d.apiClientMemory[key] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: oneOffToken.Expires,
	}
	return nil
}

func (d *MemoryCacheDriver) DeleteApiClientByOneOffToken(token string) *contract.AuthError {
	delete(d.apiClientMemory, d.getPrefix("auth")+"-one_off-"+token)
	return nil
}

func (d *MemoryCacheDriver) GetApiUserByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	key := d.getPrefix("auth") + token
	if hit, ok := d.apiUserMemory[key]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiUserMemory, key)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiUserByToken(token string, user contract.ApiUserInterface) *contract.AuthError {
	key := d.getPrefix("auth") + token
	d.apiUserMemory[key] = MemoryCacheEntry[contract.ApiUserInterface]{
		Value:    user,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetFUPEntry(key string) (*contract.FUPCacheEntry, *contract.AuthError) {
	entryKey := d.getPrefix("fup") + key
	if hit, ok := d.fupMemory[entryKey]; ok {
		return &hit.Value, nil
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
	}, nil
}

func (d *MemoryCacheDriver) SetFUPEntry(key string, entry *contract.FUPCacheEntry) *contract.AuthError {
	d.fupMemory[d.getPrefix("fup")+key] = MemoryCacheEntry[contract.FUPCacheEntry]{
		Value: *entry,
	}
	return nil
}

func NewMemoryCacheDriver() *MemoryCacheDriver {
	return &MemoryCacheDriver{
		apiClientMemory: make(map[string]MemoryCacheEntry[contract.ApiClientInterface]),
		apiUserMemory:   make(map[string]MemoryCacheEntry[contract.ApiUserInterface]),
		fupMemory:       make(map[string]MemoryCacheEntry[contract.FUPCacheEntry]),
	}
}

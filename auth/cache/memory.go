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

func (d *MemoryCacheDriver) GetApiClientByIdAndSecret(id string, secret string) (contract.ApiClientInterface, *contract.AuthError) {
	if hit, ok := d.apiClientMemory[d.prefix+id+secret]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, d.prefix+id+secret)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByIdAndSecret(id string, secret string, client contract.ApiClientInterface) *contract.AuthError {
	d.apiClientMemory[d.prefix+id+secret] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetApiClientByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	if hit, ok := d.apiClientMemory[d.prefix+apiKey]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiClientMemory, d.prefix+apiKey)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiClientByApiKey(apiKey string, client contract.ApiClientInterface) *contract.AuthError {
	d.apiClientMemory[d.prefix+apiKey] = MemoryCacheEntry[contract.ApiClientInterface]{
		Value:    client,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetApiUserByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	if hit, ok := d.apiUserMemory[d.prefix+token]; ok {
		if hit.ExpireAt.After(time.Now()) {
			return hit.Value, nil
		}
		delete(d.apiUserMemory, d.prefix+token)
	}
	return nil, nil
}

func (d *MemoryCacheDriver) SetApiUserByToken(token string, user contract.ApiUserInterface) *contract.AuthError {
	d.apiUserMemory[d.prefix+token] = MemoryCacheEntry[contract.ApiUserInterface]{
		Value:    user,
		ExpireAt: time.Now().Add(d.ttl),
	}
	return nil
}

func (d *MemoryCacheDriver) GetFUPEntry(key string) (*contract.FUPCacheEntry, *contract.AuthError) {
	if hit, ok := d.fupMemory[d.prefix+key]; ok {
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
	d.fupMemory[d.prefix+key] = MemoryCacheEntry[contract.FUPCacheEntry]{
		Value: *entry,
	}
	return nil
}

func NewMemoryCacheDriver() *MemoryCacheDriver {
	return &MemoryCacheDriver{
		apiClientMemory: make(map[string]MemoryCacheEntry[contract.ApiClientInterface]),
		apiUserMemory:   make(map[string]MemoryCacheEntry[contract.ApiUserInterface]),
	}
}

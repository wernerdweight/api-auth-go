package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"github.com/wernerdweight/api-auth-go/v2/auth/marshaller"
	"time"
)

type RedisCacheDriver struct {
	dsn          string
	client       *redis.Client
	prefix       string
	ttl          time.Duration
	newApiClient func() contract.ApiClientInterface
	newApiUser   func() contract.ApiUserInterface
}

func (d *RedisCacheDriver) getClient() *redis.Client {
	if d.client == nil {
		opts, err := redis.ParseURL(d.dsn)
		if nil != err {
			panic(err)
		}
		d.client = redis.NewClient(opts)
	}
	return d.client
}

func (d *RedisCacheDriver) getPrefix(groupPrefix GroupType) string {
	return getPrefix(d.prefix, groupPrefix)
}

func (d *RedisCacheDriver) unmarshalClient(value string) (contract.ApiClientInterface, *contract.AuthError) {
	apiClient := d.newApiClient()
	err := json.Unmarshal([]byte(value), &apiClient)
	if nil != err {
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return apiClient, nil
}

func (d *RedisCacheDriver) unmarshalUser(value string) (contract.ApiUserInterface, *contract.AuthError) {
	apiUser := d.newApiUser()
	err := json.Unmarshal([]byte(value), &apiUser)
	if nil != err {
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return apiUser, nil
}

func (d *RedisCacheDriver) Init(prefix string, ttl time.Duration) *contract.AuthError {
	d.prefix = prefix
	d.ttl = ttl
	return nil
}

func (d *RedisCacheDriver) GetApiClientByIdAndSecret(id string, secret string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + id + secret
	value, err := d.getClient().Get(context.Background(), key).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalClient(value)
}

func (d *RedisCacheDriver) SetApiClientByIdAndSecret(id string, secret string, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + id + secret
	marshalled, authErr := marshaller.MarshalInternal(client)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), key, value, d.ttl).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetApiClientByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + apiKey
	value, err := d.getClient().Get(context.Background(), key).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalClient(value)
}

func (d *RedisCacheDriver) SetApiClientByApiKey(apiKey string, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + apiKey
	marshalled, authErr := marshaller.MarshalInternal(client)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), key, value, d.ttl).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetApiClientByOneOffToken(token string) (contract.ApiClientInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + "-one_off-" + token
	value, err := d.getClient().Get(context.Background(), key).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalClient(value)
}

func (d *RedisCacheDriver) SetApiClientByOneOffToken(oneOffToken contract.OneOffToken, client contract.ApiClientInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + "-one_off-" + oneOffToken.Value
	marshalled, authErr := marshaller.MarshalInternal(client)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), key, value, oneOffToken.Expires.Sub(time.Now())).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) DeleteApiClientByOneOffToken(token string) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + "-one_off-" + token
	err := d.getClient().Del(context.Background(), key).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetApiUserByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	key := d.getPrefix(GroupTypeAuth) + token
	value, err := d.getClient().Get(context.Background(), key).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalUser(value)
}

func (d *RedisCacheDriver) SetApiUserByToken(token string, user contract.ApiUserInterface) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + token
	marshalled, authErr := marshaller.MarshalInternal(user)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), key, value, d.ttl).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetFUPEntry(key string) (*contract.FUPCacheEntry, *contract.AuthError) {
	entryKey := d.getPrefix(GroupTypeFUP) + key
	value, err := d.getClient().Get(context.Background(), entryKey).Result()
	if nil != err {
		if redis.Nil == err {
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
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	entry := &contract.FUPCacheEntry{}
	err = json.Unmarshal([]byte(value), entry)
	if nil != err {
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return entry, nil
}

func (d *RedisCacheDriver) SetFUPEntry(key string, entry *contract.FUPCacheEntry) *contract.AuthError {
	entryKey := d.getPrefix(GroupTypeFUP) + key
	value, err := json.Marshal(entry)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), entryKey, value, 0).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) InvalidateToken(token string) *contract.AuthError {
	key := d.getPrefix(GroupTypeAuth) + token
	err := d.getClient().Del(context.Background(), key).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func NewRedisCacheDriver(dsn string, newApiClient func() contract.ApiClientInterface, newApiUser func() contract.ApiUserInterface) *RedisCacheDriver {
	return &RedisCacheDriver{
		dsn:          dsn,
		newApiClient: newApiClient,
		newApiUser:   newApiUser,
	}
}

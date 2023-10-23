package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/marshaller"
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
	value, err := d.getClient().Get(context.Background(), d.prefix+id+secret).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalClient(value)
}

func (d *RedisCacheDriver) SetApiClientByIdAndSecret(id string, secret string, client contract.ApiClientInterface) *contract.AuthError {
	marshalled, authErr := marshaller.MarshalInternal(client)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), d.prefix+id+secret, value, d.ttl).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetApiClientByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	value, err := d.getClient().Get(context.Background(), d.prefix+apiKey).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalClient(value)
}

func (d *RedisCacheDriver) SetApiClientByApiKey(apiKey string, client contract.ApiClientInterface) *contract.AuthError {
	marshalled, authErr := marshaller.MarshalInternal(client)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), d.prefix+apiKey, value, d.ttl).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetApiClientByOneOffToken(token string) (contract.ApiClientInterface, *contract.AuthError) {
	value, err := d.getClient().Get(context.Background(), d.prefix+"-one_off-"+token).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalClient(value)
}

func (d *RedisCacheDriver) SetApiClientByOneOffToken(oneOffToken contract.OneOffToken, client contract.ApiClientInterface) *contract.AuthError {
	marshalled, authErr := marshaller.MarshalInternal(client)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), d.prefix+"-one_off-"+oneOffToken.Value, value, oneOffToken.Expires.Sub(time.Now())).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) DeleteApiClientByOneOffToken(token string) *contract.AuthError {
	err := d.getClient().Del(context.Background(), d.prefix+"-one_off-"+token).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetApiUserByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	value, err := d.getClient().Get(context.Background(), d.prefix+token).Result()
	if nil != err {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return d.unmarshalUser(value)
}

func (d *RedisCacheDriver) SetApiUserByToken(token string, user contract.ApiUserInterface) *contract.AuthError {
	marshalled, authErr := marshaller.MarshalInternal(user)
	if nil != authErr {
		return authErr
	}
	value, err := json.Marshal(marshalled)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), d.prefix+token, value, d.ttl).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) GetFUPEntry(key string) (*contract.FUPCacheEntry, *contract.AuthError) {
	value, err := d.getClient().Get(context.Background(), d.prefix+key).Result()
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
	value, err := json.Marshal(entry)
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	err = d.getClient().Set(context.Background(), d.prefix+key, value, 0).Err()
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

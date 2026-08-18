package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/wernerdweight/api-auth-go/v3/auth/constants"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
	"github.com/wernerdweight/api-auth-go/v3/auth/marshaller"
	"sync"
	"time"
)

// fupIncrementScript reads, increments and writes a FUP entry in a single round trip so that
// concurrent requests sharing a FUP key can't overwrite each other's counters.
//
// KEYS[1] is the entry key, ARGV[1] the current timestamp to store, ARGV[2] the entry TTL in
// seconds, and the remaining arguments are (period, from, to) triplets - the counter of a period
// is incremented if the stored timestamp falls within its bounds and reset to 1 otherwise (see
// constants.Period.GetTimestampBounds). The stored value keeps the contract.FUPCacheEntry format.
//
// The timestamp and the bounds are sampled before the script is queued, so a request that raced
// another one across a period boundary can arrive with arguments that are already stale. Such a
// request only increments the counters (it must not reset them back to 1 and must not move the
// stored timestamp backwards), which can count it towards the newer period instead of the one it
// belongs to - it is preferable to let a limiter count a boundary request twice than to let a
// stale request wipe the counter the requests behind it are limited by.
//
// Staleness is decided by comparing the serialized timestamps, which orders them correctly except
// within a single second (time.RFC3339Nano trims trailing zeros, so a whole second sorts after the
// fractions of the same second). Two timestamps of the same second always belong to the same
// period, so the only effect there is that the stored timestamp may stay behind by less than a
// second - never that a counter is reset when it should not be.
var fupIncrementScript = redis.NewScript(`
local used = {}
local updatedAt = ''
local raw = redis.call('GET', KEYS[1])
if raw then
	local decoded = cjson.decode(raw)
	if type(decoded) == 'table' then
		if type(decoded['updatedAt']) == 'string' then
			updatedAt = decoded['updatedAt']
		end
		if type(decoded['used']) == 'table' then
			used = decoded['used']
		end
	end
end
local stale = updatedAt > ARGV[1]
for i = 3, #ARGV, 3 do
	local period = ARGV[i]
	local from = ARGV[i + 1]
	local to = ARGV[i + 2]
	local current = string.sub(updatedAt, 1, string.len(from))
	local previous = tonumber(used[period])
	if previous ~= nil and (stale or (current >= from and current <= to)) then
		used[period] = previous + 1
	else
		used[period] = 1
	end
end
local storedAt = ARGV[1]
if stale then
	storedAt = updatedAt
end
local entry = cjson.encode({ updatedAt = storedAt, used = used })
redis.call('SET', KEYS[1], entry, 'EX', ARGV[2])
return entry
`)

type RedisCacheDriver struct {
	dsn          string
	client       *redis.Client
	clientOnce   sync.Once
	prefix       string
	ttl          time.Duration
	newApiClient func() contract.ApiClientInterface
	newApiUser   func() contract.ApiUserInterface
}

// getClient initializes the client on first use. It is called from concurrently handled requests,
// so the initialization has to be guarded - otherwise each of them could create its own client
// (and its own connection pool).
func (d *RedisCacheDriver) getClient() *redis.Client {
	d.clientOnce.Do(func() {
		opts, err := redis.ParseURL(d.dsn)
		if nil != err {
			panic(err)
		}
		d.client = redis.NewClient(opts)
	})
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
	err = d.getClient().Set(context.Background(), entryKey, value, constants.FUPEntryTTL).Err()
	if nil != err {
		return contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return nil
}

func (d *RedisCacheDriver) IncrementFUPEntry(key string) (*contract.FUPCacheEntry, *contract.AuthError) {
	return d.IncrementFUPEntryWithTTL(key, constants.FUPEntryTTL)
}

func (d *RedisCacheDriver) IncrementFUPEntryWithTTL(key string, ttl time.Duration) (*contract.FUPCacheEntry, *contract.AuthError) {
	entryKey := d.getPrefix(GroupTypeFUP) + key
	updatedAt := time.Now()
	args := make([]any, 0, 2+len(constants.FUPScopePeriods)*3)
	args = append(args, updatedAt.Format(time.RFC3339Nano), int(ttl.Seconds()))
	for _, period := range constants.FUPScopePeriods {
		from, to := period.GetTimestampBounds(updatedAt)
		args = append(args, string(period), from, to)
	}
	value, err := fupIncrementScript.Run(context.Background(), d.getClient(), []string{entryKey}, args...).Text()
	if nil != err {
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	entry := &contract.FUPCacheEntry{}
	err = json.Unmarshal([]byte(value), entry)
	if nil != err {
		return nil, contract.NewInternalError(contract.CacheError, map[string]string{"details": err.Error()})
	}
	return entry, nil
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

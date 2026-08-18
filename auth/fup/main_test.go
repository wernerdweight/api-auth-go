package fup

import (
	"github.com/wernerdweight/api-auth-go/v3/auth/constants"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
	"reflect"
	"testing"
	"time"
)

func Test_mergeLimits(t *testing.T) {
	type args struct {
		limits     map[constants.Period]contract.FUPLimits
		pathLimits map[constants.Period]contract.FUPLimits
	}
	tests := []struct {
		name string
		args args
		want map[constants.Period]contract.FUPLimits
	}{
		{
			name: "Nil limits",
			args: args{
				limits:     nil,
				pathLimits: nil,
			},
			want: nil,
		},
		{
			name: "Nil path limits",
			args: args{
				limits: map[constants.Period]contract.FUPLimits{
					constants.PeriodHourly: {
						Limit:  1,
						Used:   0,
						Period: constants.PeriodHourly,
					},
				},
				pathLimits: nil,
			},
			want: map[constants.Period]contract.FUPLimits{
				constants.PeriodHourly: {
					Limit:  1,
					Used:   0,
					Period: constants.PeriodHourly,
				},
			},
		},
		{
			name: "Limits",
			args: args{
				limits: map[constants.Period]contract.FUPLimits{
					constants.PeriodHourly: {
						Limit:  2,
						Used:   0,
						Period: constants.PeriodHourly,
					},
				},
				pathLimits: map[constants.Period]contract.FUPLimits{
					constants.PeriodHourly: {
						Limit:  3,
						Used:   2,
						Period: constants.PeriodHourly,
					},
				},
			},
			want: map[constants.Period]contract.FUPLimits{
				constants.PeriodHourly: {
					Limit:  3,
					Used:   2,
					Period: constants.PeriodHourly,
				},
			},
		},
		{
			name: "Multiple Limits",
			args: args{
				limits: map[constants.Period]contract.FUPLimits{
					constants.PeriodHourly: {
						Limit:  2,
						Used:   0,
						Period: constants.PeriodHourly,
					},
					constants.PeriodDaily: {
						Limit:  20,
						Used:   18,
						Period: constants.PeriodDaily,
					},
				},
				pathLimits: map[constants.Period]contract.FUPLimits{
					constants.PeriodHourly: {
						Limit:  3,
						Used:   2,
						Period: constants.PeriodHourly,
					},
					constants.PeriodDaily: {
						Limit:  30,
						Used:   25,
						Period: constants.PeriodDaily,
					},
				},
			},
			want: map[constants.Period]contract.FUPLimits{
				constants.PeriodHourly: {
					Limit:  3,
					Used:   2,
					Period: constants.PeriodHourly,
				},
				constants.PeriodDaily: {
					Limit:  20,
					Used:   18,
					Period: constants.PeriodDaily,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeLimits(tt.args.limits, tt.args.pathLimits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeLimits() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ttlDriver implements the optional contract.FUPTTLCacheDriverInterface and records what it was
// asked for. The embedded interface is nil - anything but the increment panics, which is the point
type ttlDriver struct {
	contract.CacheDriverInterface
	ttl        time.Duration
	withoutTTL bool
}

func (d *ttlDriver) IncrementFUPEntryWithTTL(_ string, ttl time.Duration) (*contract.FUPCacheEntry, *contract.AuthError) {
	d.ttl = ttl
	return &contract.FUPCacheEntry{UpdatedAt: time.Now(), Used: map[constants.Period]int{constants.PeriodMinutely: 1}}, nil
}

func (d *ttlDriver) IncrementFUPEntry(_ string) (*contract.FUPCacheEntry, *contract.AuthError) {
	d.withoutTTL = true
	return &contract.FUPCacheEntry{UpdatedAt: time.Now(), Used: map[constants.Period]int{constants.PeriodMinutely: 1}}, nil
}

// legacyDriver is a driver that only implements contract.CacheDriverInterface, i.e. one written
// against v3.0.0 - it has to keep working, just without the shortened expiration
type legacyDriver struct {
	contract.CacheDriverInterface
	incremented bool
}

func (d *legacyDriver) IncrementFUPEntry(_ string) (*contract.FUPCacheEntry, *contract.AuthError) {
	d.incremented = true
	return &contract.FUPCacheEntry{UpdatedAt: time.Now(), Used: map[constants.Period]int{constants.PeriodMinutely: 1}}, nil
}

// Test_checkLimits_EntryTTL covers that the expiration handed to the cache comes from the scope
// being enforced, not from the longest period the package supports
func Test_checkLimits_EntryTTL(t *testing.T) {
	tests := []struct {
		name  string
		scope contract.FUPScope
		want  time.Duration
	}{
		{
			// the live anonymous scope: nothing longer than a day is limited, so a per-IP key
			// must not be kept for the 35 days the monthly period needs
			name:  "minutely and daily",
			scope: contract.FUPScope{constants.FUPIPKey: map[string]any{"minutely": 60, "daily": 5000}},
			want:  constants.PeriodDaily.GetEntryTTL(),
		},
		{
			name:  "monthly",
			scope: contract.FUPScope{constants.FUPIPKey: map[string]any{"monthly": 100000}},
			want:  constants.FUPEntryTTL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &ttlDriver{}
			scope := tt.scope
			_, scopeLimits := checkLimits(&scope, constants.AnonymousFUPKey, "192.0.2.10", constants.FUPIPKey, driver)
			if nil != scopeLimits {
				t.Fatalf("checkLimits() = %+v, want no limits to be hit", scopeLimits)
			}
			if driver.withoutTTL {
				t.Errorf("checkLimits() used IncrementFUPEntry, want the TTL-aware increment")
			}
			if tt.want != driver.ttl {
				t.Errorf("checkLimits() ttl = %v, want %v", driver.ttl, tt.want)
			}
		})
	}
}

// Test_checkLimits_LegacyDriver covers that a driver written against v3.0.0, which doesn't
// implement the optional TTL interface, is still used - it keeps constants.FUPEntryTTL of its own
func Test_checkLimits_LegacyDriver(t *testing.T) {
	driver := &legacyDriver{}
	scope := contract.FUPScope{constants.FUPIPKey: map[string]any{"minutely": 60}}
	limits, scopeLimits := checkLimits(&scope, constants.AnonymousFUPKey, "192.0.2.10", constants.FUPIPKey, driver)
	if nil != scopeLimits {
		t.Fatalf("checkLimits() = %+v, want no limits to be hit", scopeLimits)
	}
	if !driver.incremented {
		t.Errorf("checkLimits() did not increment through IncrementFUPEntry")
	}
	if got, ok := limits[constants.PeriodMinutely]; !ok || 60 != got.Limit || 1 != got.Used {
		t.Errorf("checkLimits() minutely = %+v, want limit 60 used 1", got)
	}
}

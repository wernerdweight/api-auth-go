package cache

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/wernerdweight/api-auth-go/v3/auth/constants"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
)

func newTestRedisDriver(t *testing.T) (*RedisCacheDriver, *miniredis.Miniredis) {
	t.Helper()
	server := miniredis.RunT(t)
	driver := NewRedisCacheDriver("redis://"+server.Addr(), nil, nil)
	driver.Init("test:", time.Hour)
	return driver, server
}

// storeFUPEntry writes an entry the same way the driver does, but with an arbitrary timestamp
func storeFUPEntry(t *testing.T, server *miniredis.Miniredis, key string, updatedAt time.Time, used map[constants.Period]int) {
	t.Helper()
	value, err := json.Marshal(&contract.FUPCacheEntry{UpdatedAt: updatedAt, Used: used})
	if nil != err {
		t.Fatalf("can't marshal entry: %v", err)
	}
	if err := server.Set("test:fup_"+key, string(value)); nil != err {
		t.Fatalf("can't store entry: %v", err)
	}
}

// stableNow returns the current time, making sure the assertions that follow are not evaluated
// across a minute boundary (which would reset the counters the test expects to be incremented)
func stableNow(t *testing.T) time.Time {
	t.Helper()
	const margin = 2 * time.Second
	now := time.Now()
	if now.Add(margin).Minute() != now.Minute() {
		time.Sleep(margin)
		now = time.Now()
	}
	return now
}

// earlierMinuteOfSameHour returns a timestamp of an earlier minute of the same hour as the given
// one, which is what an entry written by the previous request looks like once a minute has passed.
// There is no such minute during the first minute of an hour, hence the second return value.
func earlierMinuteOfSameHour(t *testing.T, now time.Time) (time.Time, bool) {
	t.Helper()
	if 0 == now.Minute() {
		return time.Time{}, false
	}
	return now.Add(-time.Minute), true
}

func TestRedisCacheDriver_IncrementFUPEntry(t *testing.T) {
	tests := []struct {
		name      string
		updatedAt func(t *testing.T, now time.Time) time.Time
		used      map[constants.Period]int
		want      map[constants.Period]int
	}{
		{
			name: "no entry yet",
			want: map[constants.Period]int{
				constants.PeriodMinutely: 1,
				constants.PeriodHourly:   1,
				constants.PeriodDaily:    1,
				constants.PeriodWeekly:   1,
				constants.PeriodMonthly:  1,
			},
		},
		{
			name:      "same minute",
			updatedAt: func(_ *testing.T, now time.Time) time.Time { return now },
			used: map[constants.Period]int{
				constants.PeriodMinutely: 5,
				constants.PeriodHourly:   5,
				constants.PeriodDaily:    5,
				constants.PeriodWeekly:   5,
				constants.PeriodMonthly:  5,
			},
			want: map[constants.Period]int{
				constants.PeriodMinutely: 6,
				constants.PeriodHourly:   6,
				constants.PeriodDaily:    6,
				constants.PeriodWeekly:   6,
				constants.PeriodMonthly:  6,
			},
		},
		{
			name: "another minute of the same hour",
			updatedAt: func(t *testing.T, now time.Time) time.Time {
				earlier, ok := earlierMinuteOfSameHour(t, now)
				if !ok {
					t.Skip("there is no earlier minute of the same hour to store during the first minute of an hour")
				}
				return earlier
			},
			used: map[constants.Period]int{
				constants.PeriodMinutely: 5,
				constants.PeriodHourly:   5,
				constants.PeriodDaily:    5,
				constants.PeriodWeekly:   5,
				constants.PeriodMonthly:  5,
			},
			want: map[constants.Period]int{
				// only the minutely counter is reset, the longer periods keep counting
				constants.PeriodMinutely: 1,
				constants.PeriodHourly:   6,
				constants.PeriodDaily:    6,
				constants.PeriodWeekly:   6,
				constants.PeriodMonthly:  6,
			},
		},
		{
			name:      "previous year",
			updatedAt: func(_ *testing.T, now time.Time) time.Time { return now.AddDate(-1, 0, 0) },
			used: map[constants.Period]int{
				constants.PeriodMinutely: 5,
				constants.PeriodHourly:   5,
				constants.PeriodDaily:    5,
				constants.PeriodWeekly:   5,
				constants.PeriodMonthly:  5,
			},
			want: map[constants.Period]int{
				constants.PeriodMinutely: 1,
				constants.PeriodHourly:   1,
				constants.PeriodDaily:    1,
				constants.PeriodWeekly:   1,
				constants.PeriodMonthly:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver, server := newTestRedisDriver(t)
			if nil != tt.updatedAt {
				storeFUPEntry(t, server, "key", tt.updatedAt(t, stableNow(t)), tt.used)
			}
			entry, err := driver.IncrementFUPEntry("key")
			if nil != err {
				t.Fatalf("IncrementFUPEntry() error = %v", err)
			}
			for _, period := range constants.FUPScopePeriods {
				if got := entry.GetUsed(period); got != tt.want[period] {
					t.Errorf("IncrementFUPEntry() %s = %d, want %d", period, got, tt.want[period])
				}
			}
		})
	}
}

// TestRedisCacheDriver_IncrementFUPEntry_StoredFormat makes sure the stored value stays readable by
// GetFUPEntry - the counters of already deployed keys must not be reset by upgrading the package
func TestRedisCacheDriver_IncrementFUPEntry_StoredFormat(t *testing.T) {
	driver, server := newTestRedisDriver(t)
	storeFUPEntry(t, server, "key", stableNow(t), map[constants.Period]int{
		constants.PeriodMinutely: 1,
		constants.PeriodHourly:   2,
		constants.PeriodDaily:    3,
		constants.PeriodWeekly:   4,
		constants.PeriodMonthly:  5,
	})

	if _, err := driver.IncrementFUPEntry("key"); nil != err {
		t.Fatalf("IncrementFUPEntry() error = %v", err)
	}

	entry, err := driver.GetFUPEntry("key")
	if nil != err {
		t.Fatalf("GetFUPEntry() error = %v", err)
	}
	want := map[constants.Period]int{
		constants.PeriodMinutely: 2,
		constants.PeriodHourly:   3,
		constants.PeriodDaily:    4,
		constants.PeriodWeekly:   5,
		constants.PeriodMonthly:  6,
	}
	for _, period := range constants.FUPScopePeriods {
		if got := entry.GetUsed(period); got != want[period] {
			t.Errorf("GetFUPEntry() %s = %d, want %d", period, got, want[period])
		}
	}
	if entry.UpdatedAt.IsZero() {
		t.Error("GetFUPEntry() updatedAt is zero, want the increment timestamp")
	}
}

// TestRedisCacheDriver_IncrementFUPEntry_Stale covers a request whose arguments were sampled
// before a request that reached Redis first - it must not reset the counters that request already
// advanced, nor move the stored timestamp backwards
func TestRedisCacheDriver_IncrementFUPEntry_Stale(t *testing.T) {
	driver, server := newTestRedisDriver(t)
	// the entry of a request that was sampled later and already ran, as an increment sampled
	// before the minute boundary sees it
	newer := stableNow(t).Add(time.Minute)
	storeFUPEntry(t, server, "key", newer, map[constants.Period]int{
		constants.PeriodMinutely: 1,
		constants.PeriodHourly:   5,
		constants.PeriodDaily:    5,
		constants.PeriodWeekly:   5,
		constants.PeriodMonthly:  5,
	})

	entry, err := driver.IncrementFUPEntry("key")
	if nil != err {
		t.Fatalf("IncrementFUPEntry() error = %v", err)
	}

	want := map[constants.Period]int{
		constants.PeriodMinutely: 2,
		constants.PeriodHourly:   6,
		constants.PeriodDaily:    6,
		constants.PeriodWeekly:   6,
		constants.PeriodMonthly:  6,
	}
	for _, period := range constants.FUPScopePeriods {
		if got := entry.GetUsed(period); got != want[period] {
			t.Errorf("IncrementFUPEntry() %s = %d, want %d", period, got, want[period])
		}
	}
	if entry.UpdatedAt.Before(newer) {
		t.Errorf("IncrementFUPEntry() updatedAt = %v, want no earlier than %v", entry.UpdatedAt, newer)
	}
}

func TestRedisCacheDriver_IncrementFUPEntry_TTL(t *testing.T) {
	driver, server := newTestRedisDriver(t)
	// an entry stored by an older version of the package never expires
	storeFUPEntry(t, server, "key", stableNow(t), nil)
	if ttl := server.TTL("test:fup_key"); 0 != ttl {
		t.Fatalf("TTL() = %v, want no expiration", ttl)
	}

	if _, err := driver.IncrementFUPEntry("key"); nil != err {
		t.Fatalf("IncrementFUPEntry() error = %v", err)
	}

	if ttl := server.TTL("test:fup_key"); constants.FUPEntryTTL != ttl {
		t.Errorf("TTL() = %v, want %v", ttl, constants.FUPEntryTTL)
	}
}

func TestRedisCacheDriver_SetFUPEntry_TTL(t *testing.T) {
	driver, server := newTestRedisDriver(t)
	if err := driver.SetFUPEntry("key", &contract.FUPCacheEntry{UpdatedAt: time.Now()}); nil != err {
		t.Fatalf("SetFUPEntry() error = %v", err)
	}
	if ttl := server.TTL("test:fup_key"); constants.FUPEntryTTL != ttl {
		t.Errorf("TTL() = %v, want %v", ttl, constants.FUPEntryTTL)
	}
}

// TestRedisCacheDriver_IncrementFUPEntry_Concurrent covers the reason the increment happens in a
// script - the read-modify-write cycle must not lose increments of requests sharing a FUP key
func TestRedisCacheDriver_IncrementFUPEntry_Concurrent(t *testing.T) {
	driver, _ := newTestRedisDriver(t)
	const requests = 50

	var wg sync.WaitGroup
	errs := make(chan *contract.AuthError, requests)
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := driver.IncrementFUPEntry("key"); nil != err {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("IncrementFUPEntry() error = %v", err)
	}

	entry, err := driver.GetFUPEntry("key")
	if nil != err {
		t.Fatalf("GetFUPEntry() error = %v", err)
	}
	if got := entry.GetUsed(constants.PeriodDaily); got != requests {
		t.Errorf("GetFUPEntry() daily = %d, want %d", got, requests)
	}
}

// TestRedisCacheDriver_IncrementFUPEntryWithTTL_TTL covers the expiration the FUP checkers derive
// from the scope - a scope limiting nothing longer than a day must not keep per-IP keys for 35 days
func TestRedisCacheDriver_IncrementFUPEntryWithTTL_TTL(t *testing.T) {
	driver, server := newTestRedisDriver(t)
	ttl := constants.PeriodDaily.GetEntryTTL()

	entry, err := driver.IncrementFUPEntryWithTTL("key", ttl)
	if nil != err {
		t.Fatalf("IncrementFUPEntryWithTTL() error = %v", err)
	}
	if got := entry.GetUsed(constants.PeriodDaily); 1 != got {
		t.Errorf("IncrementFUPEntryWithTTL() daily = %d, want 1", got)
	}
	if got := server.TTL("test:fup_key"); ttl != got {
		t.Errorf("TTL() = %v, want %v", got, ttl)
	}

	// the expiration is refreshed on every increment, so the entry only has to outlive the gap
	// between two requests
	server.FastForward(ttl / 2)
	if _, err := driver.IncrementFUPEntryWithTTL("key", ttl); nil != err {
		t.Fatalf("IncrementFUPEntryWithTTL() error = %v", err)
	}
	if got := server.TTL("test:fup_key"); ttl != got {
		t.Errorf("TTL() = %v, want %v after a second increment", got, ttl)
	}

	// while a gap longer than the TTL releases the key instead of keeping it around
	server.FastForward(ttl + time.Minute)
	if server.Exists("test:fup_key") {
		t.Errorf("Exists() = true, want the entry to be released after %v of inactivity", ttl)
	}
}

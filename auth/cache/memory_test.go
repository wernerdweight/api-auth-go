package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/wernerdweight/api-auth-go/v3/auth/constants"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
)

func newTestMemoryDriver(t *testing.T) *MemoryCacheDriver {
	t.Helper()
	driver := NewMemoryCacheDriver()
	driver.Init("test:", time.Hour)
	return driver
}

// TestMemoryCacheDriver_IncrementFUPEntry_Concurrent mirrors the Redis driver's concurrency test -
// requests sharing a FUP key must not lose increments (and, under -race, must not share state)
func TestMemoryCacheDriver_IncrementFUPEntry_Concurrent(t *testing.T) {
	driver := newTestMemoryDriver(t)
	const requests = 50

	var wg sync.WaitGroup
	errs := make(chan *contract.AuthError, requests)
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			entry, err := driver.IncrementFUPEntry("key")
			if nil != err {
				errs <- err
				return
			}
			// the returned entry is read outside of the driver (checkLimits does the same), so it
			// must not be the state a concurrent increment writes to
			entry.GetUsed(constants.PeriodDaily)
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

// TestMemoryCacheDriver_FUPEntriesAreDetached covers the same thing deterministically - entries
// that crossed the driver's lock share no state with what is stored
func TestMemoryCacheDriver_FUPEntriesAreDetached(t *testing.T) {
	driver := newTestMemoryDriver(t)

	incremented, err := driver.IncrementFUPEntry("key")
	if nil != err {
		t.Fatalf("IncrementFUPEntry() error = %v", err)
	}
	incremented.Used[constants.PeriodDaily] = 42

	stored, err := driver.GetFUPEntry("key")
	if nil != err {
		t.Fatalf("GetFUPEntry() error = %v", err)
	}
	if got := stored.GetUsed(constants.PeriodDaily); 1 != got {
		t.Errorf("GetFUPEntry() daily = %d, want %d (the returned entry is still the stored one)", got, 1)
	}

	stored.Used[constants.PeriodDaily] = 43
	if _, err := driver.IncrementFUPEntry("key"); nil != err {
		t.Fatalf("IncrementFUPEntry() error = %v", err)
	}
	again, err := driver.GetFUPEntry("key")
	if nil != err {
		t.Fatalf("GetFUPEntry() error = %v", err)
	}
	if got := again.GetUsed(constants.PeriodDaily); 2 != got {
		t.Errorf("GetFUPEntry() daily = %d, want %d (the entry GetFUPEntry returned is still the stored one)", got, 2)
	}
}

func TestMemoryCacheDriver_SetFUPEntryIsDetached(t *testing.T) {
	driver := newTestMemoryDriver(t)
	entry := &contract.FUPCacheEntry{
		UpdatedAt: time.Now(),
		Used:      map[constants.Period]int{constants.PeriodDaily: 7},
	}
	if err := driver.SetFUPEntry("key", entry); nil != err {
		t.Fatalf("SetFUPEntry() error = %v", err)
	}

	entry.Used[constants.PeriodDaily] = 99

	stored, err := driver.GetFUPEntry("key")
	if nil != err {
		t.Fatalf("GetFUPEntry() error = %v", err)
	}
	if got := stored.GetUsed(constants.PeriodDaily); 7 != got {
		t.Errorf("GetFUPEntry() daily = %d, want %d (the caller's entry is still the stored one)", got, 7)
	}
}

// TestMemoryCacheDriver_FUPEntryExpiration covers that FUP counters are released instead of
// growing unbounded (per-IP counters of one-off visitors)
func TestMemoryCacheDriver_FUPEntryExpiration(t *testing.T) {
	driver := newTestMemoryDriver(t)
	if _, err := driver.IncrementFUPEntry("key"); nil != err {
		t.Fatalf("IncrementFUPEntry() error = %v", err)
	}

	entryKey := driver.getPrefix(GroupTypeFUP) + "key"
	stored := driver.fupMemory[entryKey]
	stored.ExpireAt = time.Now().Add(-time.Second)
	driver.fupMemory[entryKey] = stored

	entry, err := driver.GetFUPEntry("key")
	if nil != err {
		t.Fatalf("GetFUPEntry() error = %v", err)
	}
	if got := entry.GetUsed(constants.PeriodDaily); 0 != got {
		t.Errorf("GetFUPEntry() daily = %d, want %d", got, 0)
	}
	if _, ok := driver.fupMemory[entryKey]; ok {
		t.Error("the expired entry is still in memory")
	}
}

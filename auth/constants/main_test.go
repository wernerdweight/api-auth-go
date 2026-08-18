package constants

import (
	"testing"
	"time"
)

// TestPeriod_GetTimestampBounds asserts that the bounds accept exactly the timestamps that
// GetFormatToCompare considers to be within the same period - the two are used interchangeably
// (the bounds where the comparison happens in the cache, GetFormatToCompare where it happens in Go)
func TestPeriod_GetTimestampBounds(t *testing.T) {
	// a Monday that belongs to the first ISO week of the following year
	reference := time.Date(2025, 12, 29, 23, 58, 30, 0, time.UTC)
	for _, period := range FUPScopePeriods {
		t.Run(string(period), func(t *testing.T) {
			from, to := period.GetTimestampBounds(reference)
			if len(from) != len(to) || from > to {
				t.Fatalf("GetTimestampBounds() = (%q, %q), want bounds of the same length in ascending order", from, to)
			}
			for offset := -400 * 24 * time.Hour; offset <= 400*24*time.Hour; offset += 37 * time.Minute {
				other := reference.Add(offset)
				value := other.Format(time.RFC3339Nano)[:len(from)]
				got := value >= from && value <= to
				want := period.GetFormatToCompare(reference) == period.GetFormatToCompare(other)
				if got != want {
					t.Errorf("%q within (%q, %q) = %v, want %v (%s)", value, from, to, got, want, other)
				}
			}
		})
	}
}

func TestPeriod_GetTimestampBounds_Values(t *testing.T) {
	// a Wednesday
	reference := time.Date(2026, 8, 12, 14, 3, 5, 0, time.UTC)
	tests := []struct {
		period Period
		from   string
		to     string
	}{
		{PeriodMinutely, "2026-08-12T14:03", "2026-08-12T14:03"},
		{PeriodHourly, "2026-08-12T14", "2026-08-12T14"},
		{PeriodDaily, "2026-08-12", "2026-08-12"},
		{PeriodWeekly, "2026-08-10", "2026-08-16"},
		{PeriodMonthly, "2026-08", "2026-08"},
		{Period("unknown"), "", ""},
	}
	for _, tt := range tests {
		t.Run(string(tt.period), func(t *testing.T) {
			from, to := tt.period.GetTimestampBounds(reference)
			if from != tt.from || to != tt.to {
				t.Errorf("GetTimestampBounds() = (%q, %q), want (%q, %q)", from, to, tt.from, tt.to)
			}
		})
	}
}

func TestPeriod_GetTimestampBounds_WeekBoundaries(t *testing.T) {
	tests := []struct {
		name string
		day  time.Time
		from string
		to   string
	}{
		{"monday", time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC), "2026-08-10", "2026-08-16"},
		{"sunday", time.Date(2026, 8, 16, 23, 59, 59, 0, time.UTC), "2026-08-10", "2026-08-16"},
		{"next monday", time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC), "2026-08-17", "2026-08-23"},
		// the ISO week spanning the turn of the year belongs to both calendar years
		{"turn of the year", time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC), "2025-12-29", "2026-01-04"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to := PeriodWeekly.GetTimestampBounds(tt.day)
			if from != tt.from || to != tt.to {
				t.Errorf("GetTimestampBounds() = (%q, %q), want (%q, %q)", from, to, tt.from, tt.to)
			}
		})
	}
}

// TestPeriod_GetEntryTTL covers that every period keeps a counter alive for longer than the period
// it counts - an entry expiring within its own period would reset the counter early and let more
// requests through than the limit allows
func TestPeriod_GetEntryTTL(t *testing.T) {
	tests := []struct {
		period  Period
		want    time.Duration
		atLeast time.Duration
	}{
		{period: PeriodMinutely, want: time.Hour, atLeast: time.Minute},
		{period: PeriodHourly, want: time.Hour * 25, atLeast: time.Hour},
		{period: PeriodDaily, want: time.Hour * 24 * 2, atLeast: time.Hour * 24},
		{period: PeriodWeekly, want: time.Hour * 24 * 8, atLeast: time.Hour * 24 * 7},
		{period: PeriodMonthly, want: FUPEntryTTL, atLeast: time.Hour * 24 * 31},
		// an unknown period must not shorten the expiration
		{period: Period("yearly"), want: FUPEntryTTL, atLeast: time.Hour * 24 * 31},
	}
	for _, tt := range tests {
		t.Run(string(tt.period), func(t *testing.T) {
			got := tt.period.GetEntryTTL()
			if tt.want != got {
				t.Errorf("GetEntryTTL() = %v, want %v", got, tt.want)
			}
			if got <= tt.atLeast {
				t.Errorf("GetEntryTTL() = %v, want longer than the period itself (%v)", got, tt.atLeast)
			}
		})
	}
}

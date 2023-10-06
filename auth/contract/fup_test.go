package contract

import (
	"github.com/stretchr/testify/assert"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"net/http"
	"testing"
	"time"
)

func TestFUPCacheEntry_GetUsed(t *testing.T) {
	type fields struct {
		UpdatedAt time.Time
		Used      map[constants.Period]int
	}
	type args struct {
		period constants.Period
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Empty entry",
			fields: fields{
				UpdatedAt: time.Time{},
				Used:      nil,
			},
			args: args{
				period: constants.PeriodHourly,
			},
			want: 0,
		},
		{
			name: "Entry with used",
			fields: fields{
				UpdatedAt: time.Time{},
				Used: map[constants.Period]int{
					constants.PeriodHourly: 1,
				},
			},
			args: args{
				period: constants.PeriodHourly,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &FUPCacheEntry{
				UpdatedAt: tt.fields.UpdatedAt,
				Used:      tt.fields.Used,
			}
			assert.Equalf(t, tt.want, e.GetUsed(tt.args.period), "GetUsed(%v)", tt.args.period)
		})
	}
}

func TestFUPCacheEntry_Increment(t *testing.T) {
	type fields struct {
		UpdatedAt time.Time
		Used      map[constants.Period]int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Empty entry",
			fields: fields{
				UpdatedAt: time.Time{},
				Used:      nil,
			},
			want: 1,
		},
		{
			name: "Entry with used",
			fields: fields{
				UpdatedAt: time.Now(),
				Used: map[constants.Period]int{
					constants.PeriodHourly: 1,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &FUPCacheEntry{
				UpdatedAt: tt.fields.UpdatedAt,
				Used:      tt.fields.Used,
			}
			e.Increment()
			assert.Equalf(t, tt.want, e.Used[constants.PeriodHourly], "e.Used[constants.PeriodHourly]")
		})
	}
}

func TestFUPScopeLimits_GetLimitsHeader(t *testing.T) {
	type fields struct {
		Accessible constants.ScopeAccessibility
		Limits     map[constants.Period]FUPLimits
		Error      *AuthError
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Empty limits",
			fields: fields{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits:     nil,
				Error:      nil,
			},
			want: "",
		},
		{
			name: "Limits",
			fields: fields{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits: map[constants.Period]FUPLimits{
					constants.PeriodHourly: {
						Limit: 1,
						Used:  0,
					},
				},
				Error: nil,
			},
			want: "{\"hourly\":{\"limit\":1,\"used\":0}}",
		},
		{
			name: "Error",
			fields: fields{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits:     nil,
				Error: &AuthError{
					Err:     nil,
					Code:    MarshallingError,
					Payload: nil,
					Status:  http.StatusInternalServerError,
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &FUPScopeLimits{
				Accessible: tt.fields.Accessible,
				Limits:     tt.fields.Limits,
				Error:      tt.fields.Error,
			}
			assert.Equalf(t, tt.want, l.GetLimitsHeader(), "GetLimitsHeader()")
		})
	}
}

func TestFUPScopeLimits_GetRetryAfter(t *testing.T) {
	type fields struct {
		Accessible constants.ScopeAccessibility
		Limits     map[constants.Period]FUPLimits
		Error      *AuthError
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Forbidden",
			fields: fields{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits:     nil,
				Error:      nil,
			},
			want: -1,
		},
		{
			name: "Accessible",
			fields: fields{
				Accessible: constants.ScopeAccessibilityAccessible,
				Limits:     nil,
				Error:      nil,
			},
			want: -1,
		},
		{
			name: "Error",
			fields: fields{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits:     nil,
				Error: &AuthError{
					Err:     nil,
					Code:    MarshallingError,
					Payload: nil,
					Status:  http.StatusInternalServerError,
				},
			},
			want: -1,
		},
		{
			name: "Limits",
			fields: fields{
				Accessible: constants.ScopeAccessibilityForbidden,
				Limits: map[constants.Period]FUPLimits{
					constants.PeriodHourly: {
						Limit: 1,
						Used:  2,
					},
				},
				Error: nil,
			},
			want: int(time.Until(constants.PeriodHourly.GetResetTime()).Seconds()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &FUPScopeLimits{
				Accessible: tt.fields.Accessible,
				Limits:     tt.fields.Limits,
				Error:      tt.fields.Error,
			}
			assert.Equalf(t, tt.want, l.GetRetryAfter(), "GetRetryAfter()")
		})
	}
}

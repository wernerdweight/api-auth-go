package fup

import (
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"reflect"
	"testing"
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

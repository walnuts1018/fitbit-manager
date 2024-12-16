package domain

import (
	"reflect"
	"testing"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func TestNewFitbitTimeRange(t *testing.T) {
	type args struct {
		start synchro.Time[tz.AsiaTokyo]
		end   synchro.Time[tz.AsiaTokyo]
	}
	tests := []struct {
		name string
		args args
		want []FitbitTimeRange[tz.AsiaTokyo]
	}{
		{
			name: "same day",
			args: args{
				start: synchro.New[tz.AsiaTokyo](2024, 1, 1, 0, 0, 0, 0),
				end:   synchro.New[tz.AsiaTokyo](2024, 1, 1, 22, 0, 0, 0),
			},
			want: []FitbitTimeRange[tz.AsiaTokyo]{
				{
					date: fitbitDate{
						year:  2024,
						month: 1,
						day:   1,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   22,
						minute: 0,
					},
				},
			},
		},
		{
			name: "different day",
			args: args{
				start: synchro.New[tz.AsiaTokyo](2024, 1, 1, 0, 0, 0, 0),
				end:   synchro.New[tz.AsiaTokyo](2024, 1, 2, 15, 0, 0, 0),
			},
			want: []FitbitTimeRange[tz.AsiaTokyo]{
				{
					date: fitbitDate{
						year:  2024,
						month: 1,
						day:   1,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   23,
						minute: 59,
					},
				},
				{
					date: fitbitDate{
						year:  2024,
						month: 1,
						day:   2,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   15,
						minute: 0,
					},
				},
			},
		},
		{
			name: "different month",
			args: args{
				start: synchro.New[tz.AsiaTokyo](2024, 1, 31, 0, 0, 0, 0),
				end:   synchro.New[tz.AsiaTokyo](2024, 2, 2, 15, 0, 0, 0),
			},
			want: []FitbitTimeRange[tz.AsiaTokyo]{
				{
					date: fitbitDate{
						year:  2024,
						month: 1,
						day:   31,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   23,
						minute: 59,
					},
				},
				{
					date: fitbitDate{
						year:  2024,
						month: 2,
						day:   1,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   23,
						minute: 59,
					},
				},
				{
					date: fitbitDate{
						year:  2024,
						month: 2,
						day:   2,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   15,
						minute: 0,
					},
				},
			},
		},
		{
			name: "startTimeとendTimeが違う日で、endTimeが00:00のとき",
			args: args{
				start: synchro.New[tz.AsiaTokyo](2024, 1, 1, 0, 0, 0, 0), // 2024-01-01 00:00:00
				end:   synchro.New[tz.AsiaTokyo](2024, 1, 2, 0, 0, 1, 0), // 2024-01-02 00:00:01
			},
			want: []FitbitTimeRange[tz.AsiaTokyo]{
				{
					date: fitbitDate{
						year:  2024,
						month: 1,
						day:   1,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   23,
						minute: 59,
					},
				},
			},
		},
		{
			name: "startTimeとendTimeが違う日で、startTimeが23:59のとき",
			args: args{
				start: synchro.New[tz.AsiaTokyo](2024, 1, 1, 23, 59, 0, 0), // 2024-01-01 23:59:00
				end:   synchro.New[tz.AsiaTokyo](2024, 1, 2, 1, 0, 0, 0),   // 2024-01-02 01:00:00
			},
			want: []FitbitTimeRange[tz.AsiaTokyo]{
				{
					date: fitbitDate{
						year:  2024,
						month: 1,
						day:   2,
					},
					startTime: fitbitTime{
						hour:   0,
						minute: 0,
					},
					endTime: fitbitTime{
						hour:   1,
						minute: 0,
					},
				},
			},
		},
		{
			name: "startTimeとendTimeが違う日で、startTimeが23:59でendTimeが00:00のとき",
			args: args{
				start: synchro.New[tz.AsiaTokyo](2024, 1, 1, 23, 59, 0, 0), // 2024-01-01 23:59:00
				end:   synchro.New[tz.AsiaTokyo](2024, 1, 2, 0, 0, 0, 0),   // 2024-01-02 00:00:00
			},
			want: []FitbitTimeRange[tz.AsiaTokyo]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFitbitTimeRange(tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFitbitTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

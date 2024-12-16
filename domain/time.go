package domain

import (
	"fmt"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type fitbitDate struct {
	year  int
	month time.Month
	day   int
}

type fitbitTime struct {
	hour   int
	minute int
}

type FitbitTimeRange struct {
	date      fitbitDate
	startTime fitbitTime
	endTime   fitbitTime
}

func NewFitbitTimeRange[T synchro.TimeZone](start, end synchro.Time[T]) []FitbitTimeRange {
	startUTC := synchro.ConvertTz[tz.UTC](start)
	endUTC := synchro.ConvertTz[tz.UTC](end)

	now := startUTC
	timeRanges := make([]FitbitTimeRange, 0)

	for datetimeBefore(now, endUTC) {
		if dateEqual(now, endUTC) {
			timeRanges = append(timeRanges, FitbitTimeRange{
				date: fitbitDate{
					year:  now.Year(),
					month: now.Month(),
					day:   now.Day(),
				},
				startTime: fitbitTime{
					hour:   now.Hour(),
					minute: now.Minute(),
				},
				endTime: fitbitTime{
					hour:   endUTC.Hour(),
					minute: endUTC.Minute(),
				},
			})
			break
		} else {
			if !(now.Hour() == 23 && now.Minute() == 59) {
				timeRanges = append(timeRanges, FitbitTimeRange{
					date: fitbitDate{
						year:  now.Year(),
						month: now.Month(),
						day:   now.Day(),
					},
					startTime: fitbitTime{
						hour:   now.Hour(),
						minute: now.Minute(),
					},
					endTime: fitbitTime{
						hour:   23,
						minute: 59,
					},
				})
			}
		}
		now = synchro.New[tz.UTC](now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0)
	}

	return timeRanges
}

func (r FitbitTimeRange) Date() string {
	return fmt.Sprintf("%04d-%02d-%02d", r.date.year, r.date.month, r.date.day)
}

func (r FitbitTimeRange) StartTime() string {
	return fmt.Sprintf("%02d:%02d", r.startTime.hour, r.startTime.minute)
}

func (r FitbitTimeRange) EndTime() string {
	return fmt.Sprintf("%02d:%02d", r.endTime.hour, r.endTime.minute)
}

func dateEqual[T synchro.TimeZone](a, b synchro.Time[T]) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

// 秒以下を無視
func datetimeBefore[T synchro.TimeZone](a, b synchro.Time[T]) bool {
	newA := synchro.New[tz.UTC](a.Year(), a.Month(), a.Day(), a.Hour(), a.Minute(), 0, 0)
	newB := synchro.New[tz.UTC](b.Year(), b.Month(), b.Day(), b.Hour(), b.Minute(), 0, 0)
	return newA.Before(newB)
}

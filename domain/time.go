package domain

import (
	"fmt"
	"time"

	"github.com/Code-Hex/synchro"
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

type FitbitTimeRange[T synchro.TimeZone] struct {
	date      fitbitDate
	startTime fitbitTime
	endTime   fitbitTime
}

func NewFitbitTimeRange[T synchro.TimeZone](start, end synchro.Time[T]) []FitbitTimeRange[T] {
	now := start
	timeRanges := make([]FitbitTimeRange[T], 0)

	for datetimeBefore(now, end) {
		if dateEqual(now, end) {
			timeRanges = append(timeRanges, FitbitTimeRange[T]{
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
					hour:   end.Hour(),
					minute: end.Minute(),
				},
			})
			break
		} else {
			if !(now.Hour() == 23 && now.Minute() == 59) {
				timeRanges = append(timeRanges, FitbitTimeRange[T]{
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
		now = synchro.New[T](now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0)
	}

	return timeRanges
}

func (r FitbitTimeRange[T]) Date() string {
	return fmt.Sprintf("%04d-%02d-%02d", r.date.year, r.date.month, r.date.day)
}

func (r FitbitTimeRange[T]) StartTime() string {
	return fmt.Sprintf("%02d:%02d", r.startTime.hour, r.startTime.minute)
}

func (r FitbitTimeRange[T]) EndTime() string {
	return fmt.Sprintf("%02d:%02d", r.endTime.hour, r.endTime.minute)
}

func dateEqual[T synchro.TimeZone](a, b synchro.Time[T]) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

// 秒以下を無視
func datetimeBefore[T synchro.TimeZone](a, b synchro.Time[T]) bool {
	newA := synchro.New[T](a.Year(), a.Month(), a.Day(), a.Hour(), a.Minute(), 0, 0)
	newB := synchro.New[T](b.Year(), b.Month(), b.Day(), b.Hour(), b.Minute(), 0, 0)
	return newA.Before(newB)
}

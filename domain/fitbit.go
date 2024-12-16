package domain

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type HeartData struct {
	Time  synchro.Time[tz.AsiaTokyo]
	Value int `json:"value"`
}

type HeartDetail string

const (
	HeartDetailOneSecond     HeartDetail = "1sec"
	HeartDetailOneMinute     HeartDetail = "1min"
	HeartDetailFiveMinute    HeartDetail = "5min"
	HeartDetailFifteenMinute HeartDetail = "15min"
)

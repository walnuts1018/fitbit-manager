package domain

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type HeartData struct {
	Datatime *synchro.Time[tz.AsiaTokyo]
	Time     string `json:"time"`
	Value    int    `json:"value"`
}

type HeartDetail string

const (
	HeartDetailOneSecond     HeartDetail = "1sec"
	HeartDetailOneMinute     HeartDetail = "1min"
	HeartDetailFiveMinute    HeartDetail = "5min"
	HeartDetailFifteenMinute HeartDetail = "15min"
)

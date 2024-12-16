package config

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type RecordStartDatetime synchro.Time[tz.AsiaTokyo]

func ParseRecordStartDatetime(v string) (RecordStartDatetime, error) {
	t, err := synchro.ParseISO[tz.AsiaTokyo](v)
	return RecordStartDatetime(t), err
}

package timeJST

import "time"

var JST *time.Location

func init() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	JST = jst
}

func Now() time.Time {
	return time.Now().In(JST)
}

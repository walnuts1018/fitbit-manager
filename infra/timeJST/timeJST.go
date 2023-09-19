package timeJST

import "time"

func Now() time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	return time.Now().In(jst)
}

var JTC = time.FixedZone("JST", 9*60*60)

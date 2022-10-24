package clock

import "time"

var TimeNowFunc = time.Now

type Time = time.Time

func Now() Time {
	return TimeNowFunc()
}

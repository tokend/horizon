package changes

import "time"

func unixToTime(t int64) time.Time {
	return time.Unix(t, 0).UTC()
}

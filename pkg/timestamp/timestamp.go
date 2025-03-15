package timestamp

import "time"

func Time() int64 {
	return time.Now().UnixMicro()
}

package common

import "time"

type HeartBeatData struct {
	URL           URL
	HeartbeatTime time.Time
}

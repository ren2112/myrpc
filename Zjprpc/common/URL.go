package common

import "time"

type URL struct {
	InterfaceName string
	HostName      string
	Port          int
	LastHeartbeat time.Time
}

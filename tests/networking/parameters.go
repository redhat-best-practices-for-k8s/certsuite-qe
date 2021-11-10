package networking

import "time"

const (
	WaitingTime   = 5 * time.Minute
	RetryInterval = 5
)

var (
	TestNetworkingNameSpace = "networking-tests"
)

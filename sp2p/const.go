package sp2p

import "time"

// 网络传播类型
const (
	UNICAST   byte = iota
	BROADCAST
	MULTICAST
)

const (
	nodesBackupKey = "nbk:"
	DELIMITER      = "\n\r"
	MAX_BUF_LEN    = 1024 * 16
)

// Timeouts
const (
	respTimeout  = 500 * time.Millisecond
	sendTimeout  = 500 * time.Millisecond
	deadlineTime = 5 * time.Second
	expiration   = 20 * time.Second

	ntpFailureThreshold = 32               // Continuous timeouts after which to check NTP
	ntpWarningCooldown  = 10 * time.Minute // Minimum amount of time to pass before repeating NTP warning
	driftThreshold      = 10 * time.Second // Allowed clock drift before warning user
)

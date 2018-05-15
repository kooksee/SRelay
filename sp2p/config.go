package sp2p

import (
	"time"

	"github.com/kooksee/srelay/config"
)

type KConfig struct {
	*config.Config

	// 接收数据的最大缓存区
	MaxBufLen int

	// ntp服务器检测超时次数
	NtpFailureThreshold int
	// Minimum amount of time to pass before repeating NTP warning
	//在重复NTP警告之前需要经过的最短时间
	NtpWarningCooldown time.Duration
	// ntpPool is the NTP server to query for the current time
	NtpPool string
	// Number of measurements to do against the NTP server
	NtpChecks int

	// Kademlia bucket size
	BucketSize int

	// 网络响应超时时间
	RespTimeout time.Duration
	// 网络发送超时时间
	SendTimeout time.Duration
	// 过期时间
	Expiration time.Duration
	// Allowed clock drift before warning user
	DriftThreshold time.Duration

	// 节点列表备份时间
	NodesBackupTime time.Duration

	// 节点长度
	NodeIDBits int

	// Kademlia concurrency factor
	Alpha int
	// 节点响应的数量
	NodeResponseNumber int
	// 节点广播的数量
	NodeBroadcastNumber int
	// 节点分区的数量
	NodePartitionNumber int

	// 节点ID长度
	HashBits int
	// K桶的数量
	NBuckets int

	// version
	Version string
}

func DefaultKConfig() *KConfig {
	return &KConfig{
		MaxBufLen:           MAX_BUF_LEN,
		NtpFailureThreshold: ntpFailureThreshold,
		NtpWarningCooldown:  ntpWarningCooldown,
		NtpPool:             ntpPool,
		NtpChecks:           ntpChecks,
		BucketSize:          bucketSize,
		RespTimeout:         respTimeout,
		SendTimeout:         sendTimeout,
		Expiration:          expiration,
		DriftThreshold:      driftThreshold,
		NodesBackupTime:     nodesBackupTime,
		NodeIDBits:          NodeIDBits,
		Alpha:               alpha,
		NodeResponseNumber:  responseNodeNumber,
		NodeBroadcastNumber: broadcastNodeNumber,
		NodePartitionNumber: 8,
		HashBits:            hashBits,
		NBuckets:            nBuckets,
	}
}

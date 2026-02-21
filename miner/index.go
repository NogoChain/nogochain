package miner

import (
	"nogochain/miner/stratum"
)

// 导出Stratum相关
var (
	NewStratumServer = stratum.NewServer
)

// 导出类型

type (
	Job           = stratum.Job
	Share         = stratum.Share
	StratumServer = stratum.Server
)

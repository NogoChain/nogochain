package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// 定义监控指标
var (
	// 区块相关指标
	BlockHeight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_block_height",
			Help: "Current block height",
		},
	)

	BlockProcessingTime = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "nogochain_block_processing_time_seconds",
			Help:    "Time to process a block",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
	)

	BlockSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_block_size_bytes",
			Help: "Current block size in bytes",
		},
	)

	// 交易相关指标
	TransactionCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nogochain_transaction_count_total",
			Help: "Total number of transactions",
		},
	)

	TransactionProcessingTime = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "nogochain_transaction_processing_time_seconds",
			Help:    "Time to process a transaction",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
	)

	// 网络相关指标
	PeerCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_peer_count",
			Help: "Current number of peers",
		},
	)

	NetworkLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_network_latency_seconds",
			Help: "Network latency",
		},
	)

	// 节点性能指标
	CPUUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_cpu_usage_percent",
			Help: "CPU usage percentage",
		},
	)

	MemoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	DiskUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_disk_usage_bytes",
			Help: "Disk usage in bytes",
		},
	)

	// 矿池相关指标
	PoolHashRate = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_pool_hash_rate_ghs",
			Help: "Pool hash rate in GH/s",
		},
	)

	PoolMiners = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nogochain_pool_miners_count",
			Help: "Number of miners in the pool",
		},
	)

	PoolShares = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nogochain_pool_shares_total",
			Help: "Total number of shares submitted",
		},
	)

	// 错误指标
	ErrorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nogochain_error_count_total",
			Help: "Total number of errors",
		},
	)
)

// InitMetrics 初始化所有监控指标
func InitMetrics() {
	prometheus.MustRegister(
		BlockHeight,
		BlockProcessingTime,
		BlockSize,
		TransactionCount,
		TransactionProcessingTime,
		PeerCount,
		NetworkLatency,
		CPUUsage,
		MemoryUsage,
		DiskUsage,
		PoolHashRate,
		PoolMiners,
		PoolShares,
		ErrorCount,
	)
}

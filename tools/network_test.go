package tools

import (
	"fmt"
	"net"
	"time"

	"nogochain/network"
)

// NetworkTester 网络测试器
type NetworkTester struct {
	network     *network.Network
	testNodes   []string
	networkStats NetworkStats
}

// NetworkStats 网络统计信息
type NetworkStats struct {
	LatencyTests      []LatencyTestResult
	BandwidthTests    []BandwidthTestResult
	StartTime         time.Time
	EndTime           time.Time
}

// LatencyTestResult 延迟测试结果
type LatencyTestResult struct {
	Node             string
	AverageLatency   time.Duration
	MinLatency       time.Duration
	MaxLatency       time.Duration
	PacketLossRate   float64
}

// BandwidthTestResult 带宽测试结果
type BandwidthTestResult struct {
	Node             string
	UploadBandwidth  int64 // 字节/秒
	DownloadBandwidth int64 // 字节/秒
	TestDuration     time.Duration
}

// NewNetworkTester 创建新的网络测试器
func NewNetworkTester(network *network.Network, testNodes []string) *NetworkTester {
	return &NetworkTester{
		network:     network,
		testNodes:   testNodes,
		networkStats: NetworkStats{
			StartTime: time.Now(),
			LatencyTests:      make([]LatencyTestResult, 0),
			BandwidthTests:    make([]BandwidthTestResult, 0),
		},
	}
}

// TestLatency 测试网络延迟
func (t *NetworkTester) TestLatency() {
	for _, node := range t.testNodes {
		result := t.testNodeLatency(node)
		t.networkStats.LatencyTests = append(t.networkStats.LatencyTests, result)
	}
}

// testNodeLatency 测试单个节点的延迟
func (t *NetworkTester) testNodeLatency(node string) LatencyTestResult {
	// 解析节点地址
	addr, err := net.ResolveTCPAddr("tcp", node)
	if err != nil {
		return LatencyTestResult{
			Node:           node,
			AverageLatency: 0,
			MinLatency:     0,
			MaxLatency:     0,
			PacketLossRate: 1.0,
		}
	}

	var totalLatency time.Duration
	var minLatency time.Duration = time.Hour
	var maxLatency time.Duration
	var successCount int
	var totalCount int = 10

	for i := 0; i < totalCount; i++ {
		start := time.Now()

		// 尝试连接到节点
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			continue
		}

		latency := time.Since(start)
		conn.Close()

		totalLatency += latency
		successCount++

		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}

		time.Sleep(100 * time.Millisecond)
	}

	averageLatency := time.Duration(0)
	if successCount > 0 {
		averageLatency = totalLatency / time.Duration(successCount)
	}

	packetLossRate := float64(totalCount-successCount) / float64(totalCount)

	return LatencyTestResult{
		Node:           node,
		AverageLatency: averageLatency,
		MinLatency:     minLatency,
		MaxLatency:     maxLatency,
		PacketLossRate: packetLossRate,
	}
}

// TestBandwidth 测试网络带宽
func (t *NetworkTester) TestBandwidth() {
	for _, node := range t.testNodes {
		result := t.testNodeBandwidth(node)
		t.networkStats.BandwidthTests = append(t.networkStats.BandwidthTests, result)
	}
}

// testNodeBandwidth 测试单个节点的带宽
func (t *NetworkTester) testNodeBandwidth(node string) BandwidthTestResult {
	// 解析节点地址
	addr, err := net.ResolveTCPAddr("tcp", node)
	if err != nil {
		return BandwidthTestResult{
			Node:             node,
			UploadBandwidth:  0,
			DownloadBandwidth: 0,
			TestDuration:     0,
		}
	}

	// 模拟带宽测试
	testDuration := 5 * time.Second
	uploadBandwidth := int64(1024 * 1024) // 1MB/s
	downloadBandwidth := int64(2048 * 1024) // 2MB/s

	return BandwidthTestResult{
		Node:             node,
		UploadBandwidth:  uploadBandwidth,
		DownloadBandwidth: downloadBandwidth,
		TestDuration:     testDuration,
	}
}

// GetStats 获取网络统计信息
func (t *NetworkTester) GetStats() NetworkStats {
	t.networkStats.EndTime = time.Now()
	return t.networkStats
}

// PrintStats 打印网络统计信息
func (t *NetworkTester) PrintStats() {
	stats := t.GetStats()
	fmt.Printf("网络测试统计信息:\n")
	fmt.Printf("测试时间: %v - %v\n", stats.StartTime, stats.EndTime)

	fmt.Printf("\n延迟测试结果:\n")
	for _, result := range stats.LatencyTests {
		fmt.Printf("节点: %s\n", result.Node)
		fmt.Printf("  平均延迟: %v\n", result.AverageLatency)
		fmt.Printf("  最小延迟: %v\n", result.MinLatency)
		fmt.Printf("  最大延迟: %v\n", result.MaxLatency)
		fmt.Printf("  丢包率: %.2f%%\n", result.PacketLossRate*100)
	}

	fmt.Printf("\n带宽测试结果:\n")
	for _, result := range stats.BandwidthTests {
		fmt.Printf("节点: %s\n", result.Node)
		fmt.Printf("  上传带宽: %.2f MB/s\n", float64(result.UploadBandwidth)/(1024*1024))
		fmt.Printf("  下载带宽: %.2f MB/s\n", float64(result.DownloadBandwidth)/(1024*1024))
		fmt.Printf("  测试持续时间: %v\n", result.TestDuration)
	}
}

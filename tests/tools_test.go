package tests

import (
	"testing"
	"time"

	"nogochain/tools"
)

// TestTools 测试tools包
func TestTools(t *testing.T) {
	// 测试NewAutoTestRunner
	runner := tools.NewAutoTestRunner(1 * time.Hour)
	if runner == nil {
		t.Errorf("NewAutoTestRunner 返回 nil")
	}

	// 测试NewTxGenerator
	generator := tools.NewTxGenerator(runner.Network, 100, 30*time.Second)
	if generator == nil {
		t.Errorf("NewTxGenerator 返回 nil")
	}

	// 测试NewBlockSyncTester
	syncTester := tools.NewBlockSyncTester(runner.Blockchain, runner.Network)
	if syncTester == nil {
		t.Errorf("NewBlockSyncTester 返回 nil")
	}

	// 测试NewNetworkTester
	testNodes := []string{"127.0.0.1:8545"}
	networkTester := tools.NewNetworkTester(runner.Network, testNodes)
	if networkTester == nil {
		t.Errorf("NewNetworkTester 返回 nil")
	}
}

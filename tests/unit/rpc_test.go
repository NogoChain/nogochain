package unit

import (
	"testing"

	"nogochain/network/config"
	"nogochain/rpc"
)

// TestRPC_NewServer 测试创建新RPC服务器
func TestRPC_NewServer(t *testing.T) {
	// 创建RPC配置
	rpcConfig := &config.RPCConfig{
		Enabled: false,
		Port:    8545,
		JWT: &config.JWTConfig{
			Enabled:   false,
			Secret:    "",
			TokenFile: "",
		},
	}

	// 创建新RPC服务器
	server := rpc.NewServer(rpcConfig)

	// 验证服务器不为空
	if server == nil {
		t.Fatal("RPC服务器为nil")
	}
}

// TestRPC_Start 测试启动RPC服务器
func TestRPC_Start(t *testing.T) {
	// 创建RPC配置
	rpcConfig := &config.RPCConfig{
		Enabled: false,
		Port:    8545,
		JWT: &config.JWTConfig{
			Enabled:   false,
			Secret:    "",
			TokenFile: "",
		},
	}

	// 创建新RPC服务器
	server := rpc.NewServer(rpcConfig)

	// 启动服务器
	err := server.Start()
	if err != nil {
		t.Fatalf("启动RPC服务器失败: %v", err)
	}
}

// TestRPC_Stop 测试停止RPC服务器
func TestRPC_Stop(t *testing.T) {
	// 创建RPC配置
	rpcConfig := &config.RPCConfig{
		Enabled: false,
		Port:    8545,
		JWT: &config.JWTConfig{
			Enabled:   false,
			Secret:    "",
			TokenFile: "",
		},
	}

	// 创建新RPC服务器
	server := rpc.NewServer(rpcConfig)

	// 启动服务器
	err := server.Start()
	if err != nil {
		t.Fatalf("启动RPC服务器失败: %v", err)
	}

	// 停止服务器
	err = server.Stop()
	if err != nil {
		t.Fatalf("停止RPC服务器失败: %v", err)
	}
}

package config

// Config 网络配置
type Config struct {
	// P2P配置
	Port      int      `json:"port"`
	MaxPeers  int      `json:"maxPeers"`
	Bootnodes []string `json:"bootnodes"`

	// 发现配置
	Discovery *DiscoveryConfig `json:"discovery"`

	// 同步配置
	Sync *SyncConfig `json:"sync"`

	// RPC配置
	RPC *RPCConfig `json:"rpc"`

	// 日志配置
	Log *LogConfig `json:"log"`

	// 监控配置
	Metrics *MetricsConfig `json:"metrics"`
}

// DiscoveryConfig 节点发现配置
type DiscoveryConfig struct {
	Enabled   bool     `json:"enabled"`
	TTL       int      `json:"ttl"`
	Interval  int      `json:"interval"`
	Bootnodes []string `json:"bootnodes"`
}

// SyncConfig 同步配置
type SyncConfig struct {
	Enabled      bool `json:"enabled"`
	FastSync     bool `json:"fastSync"`
	BlockBatch   int  `json:"blockBatch"`
	MaxForkDepth int  `json:"maxForkDepth"`
}

// RPCConfig RPC配置
type RPCConfig struct {
	Enabled bool       `json:"enabled"`
	Port    int        `json:"port"`
	Host    string     `json:"host"`
	JWT     *JWTConfig `json:"jwt"`
}

// JWTConfig JWT认证配置
type JWTConfig struct {
	Enabled   bool   `json:"enabled"`
	Secret    string `json:"secret"`
	TokenFile string `json:"tokenFile"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `json:"level"`
	File     string `json:"file"`
	MaxSize  int    `json:"maxSize"` // MB
	MaxAge   int    `json:"maxAge"`  // 天
	Compress bool   `json:"compress"`
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Host    string `json:"host"`
}

// DefaultConfig 默认网络配置
func DefaultConfig() *Config {
	return &Config{
		Port:      30303,
		MaxPeers:  50,
		Bootnodes: []string{},
		Discovery: &DiscoveryConfig{
			Enabled:   true,
			TTL:       30,
			Interval:  30,
			Bootnodes: []string{},
		},
		Sync: &SyncConfig{
			Enabled:      true,
			FastSync:     true,
			BlockBatch:   128,
			MaxForkDepth: 100,
		},
		RPC: &RPCConfig{
			Enabled: true,
			Port:    8545,
			Host:    "127.0.0.1",
			JWT: &JWTConfig{
				Enabled:   true,
				Secret:    "test-secret-key-for-jwt-authentication", // 测试用密钥
				TokenFile: "jwt-token.txt",
			},
		},
		Log: &LogConfig{
			Level:    "info",
			File:     "logs/nogochain.log",
			MaxSize:  100, // 100MB
			MaxAge:   7,   // 7天
			Compress: true,
		},
		Metrics: &MetricsConfig{
			Enabled: true,
			Port:    9090,
			Host:    "127.0.0.1",
		},
	}
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 这里可以从文件加载配置
	// 目前返回默认配置
	return DefaultConfig()
}

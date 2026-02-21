package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"nogochain/core/blockchain"
	"nogochain/core/synchronizer"
	"nogochain/metrics"
	"nogochain/network"
	"nogochain/network/config"
)

// initLogger 初始化日志系统
func initLogger(cfg *config.LogConfig) {
	// 创建日志目录
	logDir := filepath.Dir(cfg.File)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	// 配置日志轮转
	logWriter := &lumberjack.Logger{
		Filename: cfg.File,
		MaxSize:  cfg.MaxSize, // MB
		MaxAge:   cfg.MaxAge,  // 天
		Compress: cfg.Compress,
	}

	// 设置zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(logWriter).With().Timestamp().Caller().Logger()

	// 设置日志级别
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}

// initMetrics 初始化Prometheus监控
func initMetrics(cfg *config.MetricsConfig) {
	if !cfg.Enabled {
		return
	}

	// 初始化监控指标
	metrics.InitMetrics()

	// 启动系统资源监控
	metrics.StartSystemMetrics()

	// 启动监控服务器
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Info().Str("addr", addr).Msg("Starting Prometheus metrics server")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Error().Err(err).Msg("Failed to start metrics server")
		}
	}()
}

func main() {
	fmt.Println("NogoChain (EVM+NogoPow)")
	fmt.Println("ChainID: 318, Symbol: NOGO, Decimals: 18")
	fmt.Println("Starting NogoChain node...")

	// 解析命令行参数
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	// 初始化网络配置
	var netConfig *config.Config
	if *configFile != "" {
		// 从文件加载配置
		data, err := ioutil.ReadFile(*configFile)
		if err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}

		// 解析JSON配置
		netConfig = &config.Config{}
		if err := json.Unmarshal(data, netConfig); err != nil {
			fmt.Printf("Failed to parse config file: %v\n", err)
			os.Exit(1)
		}

		// 确保所有配置项都有默认值
		if netConfig.Discovery == nil {
			netConfig.Discovery = &config.DiscoveryConfig{
				Enabled:   true,
				TTL:       30,
				Interval:  30,
				Bootnodes: []string{},
			}
		}
		if netConfig.Sync == nil {
			netConfig.Sync = &config.SyncConfig{
				Enabled:      true,
				FastSync:     true,
				BlockBatch:   128,
				MaxForkDepth: 100,
			}
		}
		if netConfig.RPC == nil {
			netConfig.RPC = &config.RPCConfig{
				Enabled: true,
				Port:    8545,
				Host:    "127.0.0.1",
				JWT: &config.JWTConfig{
					Enabled:   true,
					Secret:    "test-secret-key-for-jwt-authentication",
					TokenFile: "jwt-token.txt",
				},
			}
		}
		if netConfig.Log == nil {
			netConfig.Log = &config.LogConfig{
				Level:    "info",
				File:     "logs/nogochain.log",
				MaxSize:  100,
				MaxAge:   7,
				Compress: true,
			}
		}
		if netConfig.Metrics == nil {
			netConfig.Metrics = &config.MetricsConfig{
				Enabled: true,
				Port:    9090,
				Host:    "127.0.0.1",
			}
		}
	} else {
		// 使用默认配置
		netConfig = config.DefaultConfig()
	}

	// 初始化日志系统
	initLogger(netConfig.Log)

	log.Info().Msg("NogoChain node starting...")

	// 初始化区块链
	bc := blockchain.NewBlockchain(nil)
	log.Info().Str("genesisBlock", bc.Genesis().Hash().String()).Msg("Blockchain initialized")
	log.Info().Str("currentHead", bc.CurrentHead().Hash().String()).Uint64("height", bc.CurrentHead().NumberU64()).Msg("Current blockchain status")

	// 初始化交易池
	_ = blockchain.NewTransactionPool()
	log.Info().Msg("Transaction pool initialized")

	log.Info().Msg("Network config initialized")

	// 初始化网络管理器
	net := network.NewNetwork(netConfig, bc)
	log.Info().Msg("Network manager initialized")

	// 启动网络
	if err := net.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start network")
	}
	log.Info().Msg("Network started successfully")

	// 初始化同步器
	sync := synchronizer.NewSynchronizer(bc, synchronizer.FastSync)
	log.Info().Msg("Synchronizer initialized")

	// 启动同步器
	sync.Start()
	log.Info().Msg("Synchronizer started")

	// 模拟添加一个对等节点
	peer := &synchronizer.Peer{
		ID:          "peer1",
		Head:        bc.CurrentHead().Hash(),
		Td:          1000000,
		BlockNumber: bc.CurrentHead().NumberU64(),
	}
	sync.AddPeer(peer)
	log.Info().Str("peerID", peer.ID).Msg("Added peer")

	// 开始区块同步
	if err := net.SyncBlocks(); err != nil {
		log.Error().Err(err).Msg("Failed to start sync")
	} else {
		log.Info().Msg("Block sync started")
	}

	// 初始化监控
	initMetrics(netConfig.Metrics)

	log.Info().Msg("Node started successfully!")
	log.Info().Msg("NogoChain is ready for transactions and block processing")

	// 保持程序运行
	select {}
}

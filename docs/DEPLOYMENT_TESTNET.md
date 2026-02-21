# NogoChain Test Network Deployment Tutorial

## 1. Test Network Overview

### 1.1 Test Network Purpose
- Verify NogoChain core functionality
- Test network performance and stability
- Develop and test smart contracts
- Simulate production environment scenarios
- Collect user feedback and issues

### 1.2 Test Network Features
- Lower hardware requirements
- Fast synchronization and block generation
- Test-specific Chain ID
- Simulated mining difficulty
- Developer-friendly configuration

### 1.3 Test Network Parameters
| Parameter | Value | Description |
|-----------|-------|-------------|
| Chain ID | 319 | Test network chain ID |
| Block Time | 2 seconds | For fast testing |
| Mining Reward | 100 NOGO | Test tokens |
| Difficulty Adjustment | Every 5 blocks | Fast adaptation to hashrate changes |
| Maximum Block Size | 8MB | Support more transactions |
| Gas Price | 1 Gwei | Low-cost testing |

## 2. Hardware Requirements

### 2.1 Minimum Configuration
- **CPU**: 2-core processor
- **Memory**: 8GB RAM
- **Storage**: 100GB SSD
- **Network**: 50Mbps bandwidth
- **Operating System**: Windows 10+ or Linux (Ubuntu 18.04+)

### 2.2 Recommended Configuration
- **CPU**: 4-core processor
- **Memory**: 16GB RAM
- **Storage**: 250GB SSD
- **Network**: 100Mbps bandwidth
- **Operating System**: Windows 10+ or Linux (Ubuntu 20.04+)

## 3. Network Configuration

### 3.1 Port Settings
| Service | Default Port | Purpose | Need Public Access |
|---------|-------------|---------|-------------------|
| Node P2P | 30305 | Inter-node communication | Yes |
| Node RPC | 8548 | JSON-RPC interface | No (local access only) |
| Node WebSocket | 8549 | WebSocket interface | No (local access only) |
| Pool Stratum | 3336 | Miner connection port | Yes |

### 3.2 Firewall Configuration
- Allow inbound connections for P2P ports and Stratum port
- Restrict RPC ports to local access only
- Temporarily disable firewall for quick testing (test environment)

## 4. Quick Start Guide

### 4.1 One-Click Start
1. Download NogoChain test network version
2. Extract to the specified directory
3. Run `testnet/start_node.bat` script
4. The script will automatically complete all configuration and startup steps

### 4.2 Manual Startup Steps

#### 4.2.1 Environment Preparation
1. Install Go 1.22+ (if not installed)
2. Clone repository: `git clone https://github.com/nogochain/nogochain.git`
3. Enter directory: `cd nogochain`
4. Install dependencies: `go mod tidy`

#### 4.2.2 Compile Components
1. Compile node: `go build -o build/nogochain cmd/nogochain/main.go`
2. Compile mining pool: `go build -o build/nogopool cmd/nogopool/main.go`
3. Compile CLI tool: `go build -o build/nogocli cmd/nogocli/main.go`
4. Compile mining tool: `go build -o build/nogominer cmd/nogominer/main.go`

#### 4.2.3 Configuration File Setup
1. Copy test network configuration: `cp -r testnet/config/ .`
2. Modify configuration parameters in `node_config.json`
3. Modify configuration parameters in `mining_pool_config.json`

#### 4.2.4 Start Node
1. Initialize genesis block: `./nogochain.exe --datadir testnet/data init testnet/config/genesis.json`
2. Start node: `./nogochain.exe --datadir testnet/data --config testnet/config/node_config.json`

#### 4.2.5 Start Mining Pool
1. Start mining pool: `./nogopool.exe --datadir testnet/data --config testnet/config/mining_pool_config.json`

#### 4.2.6 Start Miner
1. Start miner: `./nogominer.exe --stratum 127.0.0.1:3336 --wallet 0xYourTestWalletAddress`

## 5. Testing Methods and Tools

### 5.1 Basic Functionality Testing

#### 5.1.1 Node Synchronization Test
1. Start multiple nodes
2. Monitor node synchronization status
3. Verify block height consistency
4. Test network partition and recovery

#### 5.1.2 Transaction Test
1. Generate test accounts
2. Send transactions
3. Verify transaction confirmation
4. Test high-concurrency transactions

#### 5.1.3 Smart Contract Test
1. Deploy test contracts
2. Call contract methods
3. Test contract events
4. Verify contract execution results

### 5.2 Performance Testing

#### 5.2.1 Block Processing Performance
1. Measure block generation time
2. Test block verification speed
3. Evaluate transaction processing capability
4. Analyze system resource usage

#### 5.2.2 Network Performance
1. Test inter-node communication delay
2. Evaluate network bandwidth usage
3. Measure synchronization speed
4. Test network topology changes

#### 5.2.3 Storage Performance
1. Measure database read/write speed
2. Evaluate storage growth trends
3. Test data compression effects
4. Analyze storage optimization space

### 5.3 Security Testing

#### 5.3.1 Vulnerability Scanning
1. Run security scanning tools
2. Check for known vulnerabilities
3. Test RPC interface security
4. Evaluate network protocol security

#### 5.3.2 Attack Simulation
1. Test denial of service attacks
2. Simulate double-spend attacks
3. Test 51% attack defense
4. Evaluate consensus algorithm security

### 5.4 Testing Tools

#### 5.4.1 Built-in Tools
- **nogocli**: Command-line tool for node management and query
- **nogotest**: Automated testing tool
- **tx_generator**: Transaction generator for performance testing

#### 5.4.2 Third-party Tools
- **Hardhat**: Smart contract development and testing framework
- **Foundry**: Ethereum testing framework
- **Remix**: Online smart contract editor and testing tool
- **Postman**: API testing tool
- **Grafana**: Monitoring and visualization tool

## 6. Test Network Configuration Guide

### 6.1 Configuration File Description

#### 6.1.1 Node Configuration File (node_config.json)
```json
{
  "chain_id": 319,
  "data_dir": "testnet/data",
  "log_dir": "testnet/logs",
  "p2p_port": 30305,
  "rpc_port": 8548,
  "ws_port": 8549,
  "max_peers": 20,
  "min_gas_price": "1000000000",
  "gas_limit": "16000000",
  "metrics": true,
  "metrics_port": 6061,
  "dev_mode": true,
  "dev_block_time": 2000
}
```

#### 6.1.2 Mining Pool Configuration File (mining_pool_config.json)
```json
{
  "pool_name": "NogoChain Test Pool",
  "node_url": "http://localhost:8548",
  "stratum_port": 3336,
  "payment_interval": 300,
  "min_payout": "10000000000000000000",
  "fee": 0.005,
  "reward_address": "0xYourTestRewardAddress",
  "test_mode": true,
  "simulate_hashrate": true
}
```

### 6.2 Environment Variable Configuration

#### 6.2.1 Windows Environment
```batch
set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=off
set NOGOCHAIN_DATA_DIR=testnet/data
set NOGOCHAIN_LOG_LEVEL=debug
```

#### 6.2.2 Linux Environment
```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off
export NOGOCHAIN_DATA_DIR=testnet/data
export NOGOCHAIN_LOG_LEVEL=debug
```

### 6.3 Development Mode Configuration

#### 6.3.1 Quick Development Setup
1. Enable `dev_mode`: Accelerate testing
2. Set `dev_block_time`: Adjust block time
3. Configure `debug_api`: Enable debugging API
4. Disable `metrics`: Reduce resource usage

#### 6.3.2 Test Network Reset
1. Stop all services
2. Delete data directory: `rm -rf testnet/data`
3. Reinitialize genesis block
4. Restart services

## 7. Test Result Analysis Guide

### 7.1 Log Analysis

#### 7.1.1 Node Logs
- **INFO**: Normal operation information
- **DEBUG**: Detailed debugging information
- **WARN**: Warning information
- **ERROR**: Error information
- **FATAL**: Fatal error information

#### 7.1.2 Log Analysis Tools
- **grep**: Search for specific keywords
- **tail**: View latest logs
- **logrotate**: Manage log file size
- **ELK Stack**: Centralized log analysis

### 7.2 Performance Metrics Analysis

#### 7.2.1 Key Metrics
- **TPS (Transactions Per Second)**: Transaction processing capability per second
- **Block Time**: Block generation time
- **Confirmation Time**: Transaction confirmation time
- **Resource Usage**: CPU, memory, disk, network usage
- **Sync Speed**: Block synchronization speed

#### 7.2.2 Metrics Collection Tools
- **Prometheus**: Metrics collection and storage
- **Grafana**: Metrics visualization and alerting
- **top/htop**: System resource monitoring
- **iostat**: Disk I/O monitoring

### 7.3 Test Report Template

#### 7.3.1 Basic Information
- Test network version
- Test time and duration
- Test environment configuration
- Test participants

#### 7.3.2 Test Results
- Functionality test results
- Performance test results
- Security test results
- Issues and defects

#### 7.3.3 Analysis and Recommendations
- Performance bottleneck analysis
- Security risk assessment
- Optimization suggestions
- Next test plan

### 7.4 Common Issue Analysis

#### 7.4.1 Synchronization Issues
- **Symptom**: Node cannot sync to latest block
- **Possible Causes**: Network connection issues, node configuration errors, corrupted block data
- **Solutions**: Check network connection, verify configuration files, reset data directory

#### 7.4.2 Performance Degradation
- **Symptom**: Decreased TPS, increased block time
- **Possible Causes**: Insufficient system resources, network congestion, database performance issues
- **Solutions**: Increase system resources, optimize network configuration, adjust database parameters

#### 7.4.3 Transaction Failure
- **Symptom**: Transactions are rejected or unconfirmed
- **Possible Causes**: Insufficient balance, low gas price, contract execution failure
- **Solutions**: Check account balance, adjust gas price, verify contract code

#### 7.4.4 Mining Pool Connection Issues
- **Symptom**: Miners cannot connect to mining pool
- **Possible Causes**: Network connection issues, port configuration errors, mining pool service not running
- **Solutions**: Check network connection, verify port configuration, restart mining pool service

## 8. Test Network Maintenance

### 8.1 Daily Maintenance
- Monitor node status
- Check log files
- Back up test data
- Update test network version

### 8.2 Network Reset
1. Notify test network users in advance
2. Stop all services
3. Clear test data
4. Deploy new version
5. Reinitialize network
6. Notify users that network has been reset

### 8.3 Version Upgrade
1. Download latest test version
2. Stop test services
3. Replace binary files
4. Update configuration files
5. Start services and verify

## 9. Developer Guide

### 9.1 Smart Contract Development
1. Install development tools: Hardhat, Foundry, or Remix
2. Configure test network connection
3. Write smart contracts
4. Deploy to test network
5. Test contract functionality

### 9.2 Application Development
1. Install NogoChain SDK
2. Configure API connection
3. Implement basic functionality
4. Test application performance
5. Optimize user experience

### 9.3 Contributing Code
1. Fork the repository
2. Create a branch
3. Implement features or fix issues
4. Run tests
5. Submit Pull Request

## 10. Test Network Incentives

### 10.1 Test Rewards
- Users participating in the test network can receive test token rewards
- Users who discover security vulnerabilities can receive additional rewards
- Users who provide valuable feedback can receive contribution rewards

### 10.2 Reward Application
1. Register test network account
2. Participate in test activities
3. Submit test results and feedback
4. Verify contributions
5. Claim rewards

## 11. Frequently Asked Questions

### 11.1 Technical Questions

#### 11.1.1 Node Cannot Start
**Q**: Node startup fails with "Error: Failed to load genesis file"
**A**: Check if the genesis block file path is correct and ensure the file format is correct

#### 11.1.2 Mining Pool Connection Failed
**Q**: Mining pool cannot connect to node with "Error: RPC connection failed"
**A**: Check if the node RPC port is configured correctly and ensure the node RPC service is running

#### 11.1.3 Transaction Unconfirmed
**Q**: Transaction remains unconfirmed for a long time after submission
**A**: Check if the gas price is sufficient, view node synchronization status, and confirm the network is normal

#### 11.1.4 Performance Test Results Abnormal
**Q**: Performance test results do not match expectations
**A**: Check test environment configuration, ensure sufficient system resources, and verify test methods are correct

### 11.2 Configuration Questions

#### 11.2.1 Port Conflict
**Q**: Service startup shows "Error: Address already in use"
**A**: Check if the port is occupied and modify port settings in the configuration file

#### 11.2.2 Insufficient Memory
**Q**: Runtime shows "Error: Out of memory"
**A**: Increase system memory, adjust service memory limits, and optimize configuration parameters

#### 11.2.3 Insufficient Disk Space
**Q**: Runtime shows "Error: No space left on device"
**A**: Increase disk space, enable data compression, and clean up old data

### 11.3 Network Questions

#### 11.3.1 Connection Limit
**Q**: Node connection count reaches limit
**A**: Adjust `max_peers` parameter, optimize network configuration, and ensure sufficient network bandwidth

#### 11.3.2 Network Delay
**Q**: High inter-node communication delay
**A**: Check network connection quality, optimize network topology, and use more stable network connections

## 12. Contact Information

### 12.1 Developer Community
- GitHub: https://github.com/nogochain
- Discord: https://discord.gg/nogochain
- Telegram: https://t.me/nogochain
- Twitter: https://twitter.com/nogochain

### 12.2 Technical Support
- Development Forum: https://dev.nogochain.org
- Technical Documentation: https://docs.nogochain.org
- Support Email: dev@nogochain.org

### 12.3 Test Network Monitoring
- Network Status: https://monitor.nogochain.org
- Block Explorer: https://testnet-explorer.nogochain.org
- Test Network Statistics: https://stats.nogochain.org

## 13. Appendix

### 13.1 Test Network Startup Scripts

#### 13.1.1 Windows Startup Script (start_testnet.bat)
```batch
@echo off

:: Set environment variables
set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=off
set NOGOCHAIN_DATA_DIR=testnet/data
set NOGOCHAIN_LOG_LEVEL=debug

:: Create directories
if not exist "testnet" mkdir "testnet"
if not exist "testnet/data" mkdir "testnet/data"
if not exist "testnet/logs" mkdir "testnet/logs"
if not exist "testnet/config" mkdir "testnet/config"

:: Copy configuration files
if not exist "testnet/config/genesis.json" copy "testnet/config/genesis.json" "testnet/config/genesis.json"
if not exist "testnet/config/node_config.json" copy "testnet/config/node_config.json" "testnet/config/node_config.json"
if not exist "testnet/config/mining_pool_config.json" copy "testnet/config/mining_pool_config.json" "testnet/config/mining_pool_config.json"

:: Start node
echo 启动测试网络节点...
start "NogoChain Testnet Node" /MIN nogochain.exe --datadir testnet/data --config testnet/config/node_config.json

:: Start mining pool
echo 启动测试网络矿池...
start "NogoChain Testnet Pool" /MIN nogopool.exe --datadir testnet/data --config testnet/config/mining_pool_config.json

echo 测试网络启动完成！
echo 查看日志文件了解运行状态：testnet/logs/
pause
```

#### 13.1.2 Linux Startup Script (start_testnet.sh)
```bash
#!/bin/bash

# Set environment variables
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off
export NOGOCHAIN_DATA_DIR=testnet/data
export NOGOCHAIN_LOG_LEVEL=debug

# Create directories
mkdir -p testnet/data testnet/logs testnet/config

# Copy configuration files
cp -n testnet/config/genesis.json testnet/config/genesis.json
cp -n testnet/config/node_config.json testnet/config/node_config.json
cp -n testnet/config/mining_pool_config.json testnet/config/mining_pool_config.json

# Start node
echo "启动测试网络节点..."
./nogochain --datadir testnet/data --config testnet/config/node_config.json > testnet/logs/node.log 2>&1 &

# Start mining pool
echo "启动测试网络矿池..."
./nogopool --datadir testnet/data --config testnet/config/mining_pool_config.json > testnet/logs/pool.log 2>&1 &

echo "测试网络启动完成！"
echo "查看日志文件了解运行状态：testnet/logs/"
```

### 13.2 Test Tool Commands

#### 13.2.1 Transaction Generator
```bash
# Generate 1000 transactions
./tx_generator --count 1000 --from 0xSenderAddress --to 0xReceiverAddress --value 1000000000000000000 --gas 21000 --gasprice 1000000000
```

#### 13.2.2 Performance Test Tool
```bash
# Run performance test
./performance_test --duration 3600 --tps 1000 --concurrency 100
```

#### 13.2.3 Network Test Tool
```bash
# Test network delay
./network_test --nodes 10 --messages 1000 --size 1024
```

### 13.3 Test Network Resources

#### 13.3.1 Official Resources
- Test Network Documentation: https://docs.nogochain.org/testnet
- Test Token Faucet: https://faucet.nogochain.org
- Test Network Status: https://status.nogochain.org

#### 13.3.2 Third-party Resources
- Remix IDE: https://remix.ethereum.org
- Hardhat: https://hardhat.org
- Foundry: https://getfoundry.sh
- MetaMask: https://metamask.io

## 14. Version History

| Version | Date | Changes |
|---------|------|----------|
| v0.1.0 | 2024-01-01 | Test network initial version |
| v0.1.1 | 2024-01-15 | Performance optimization |
| v0.1.2 | 2024-02-01 | Security update |
| v0.2.0 | 2024-03-01 | Feature enhancement |
| v0.2.1 | 2024-03-15 | Bug fixes |

## 15. Disclaimer

This tutorial is for testing and development purposes only. Tokens in the test network have no actual value and are only used for functionality testing. During deployment, please ensure compliance with relevant laws and regulations, and pay attention to security protection measures. The NogoChain team is not responsible for any losses caused by using this tutorial.

---

**Document Version**: v0.2.1
**Last Updated**: 2024-03-15
**Author**: NogoChain Team
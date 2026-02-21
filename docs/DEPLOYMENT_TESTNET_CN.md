# NogoChain 测试网络部署教程

## 1. 测试网络概述

### 1.1 测试网络目的
- 验证 NogoChain 核心功能
- 测试网络性能和稳定性
- 开发和测试智能合约
- 模拟生产环境场景
- 收集用户反馈和问题

### 1.2 测试网络特点
- 较低的硬件要求
- 快速同步和区块生成
- 测试专用的 Chain ID
- 模拟挖矿难度
- 开发者友好的配置

### 1.3 测试网络参数
| 参数 | 值 | 说明 |
|------|-----|------|
| Chain ID | 319 | 测试网络链 ID |
| 区块时间 | 2 秒 | 快速测试 |
| 挖矿奖励 | 100 NOGO | 测试用代币 |
| 难度调整 | 每 5 个区块 | 快速适应哈希率变化 |
| 最大区块大小 | 8MB | 支持更多交易 |
| Gas 价格 | 1 Gwei | 低成本测试 |

## 2. 硬件要求

### 2.1 最低配置
- **CPU**: 2 核处理器
- **内存**: 8GB RAM
- **存储**: 100GB SSD
- **网络**: 50Mbps 带宽
- **操作系统**: Windows 10+ 或 Linux (Ubuntu 18.04+)

### 2.2 推荐配置
- **CPU**: 4 核处理器
- **内存**: 16GB RAM
- **存储**: 250GB SSD
- **网络**: 100Mbps 带宽
- **操作系统**: Windows 10+ 或 Linux (Ubuntu 20.04+)

## 3. 网络配置

### 3.1 端口设置
| 服务 | 默认端口 | 用途 | 是否需要公网开放 |
|------|---------|------|----------------|
| 节点 P2P | 30305 | 节点间通信 | 是 |
| 节点 RPC | 8548 | JSON-RPC 接口 | 否（仅本地访问） |
| 节点 WebSocket | 8549 | WebSocket 接口 | 否（仅本地访问） |
| 矿池 Stratum | 3336 | 矿工连接端口 | 是 |

### 3.2 防火墙配置
- 允许 P2P 端口和 Stratum 端口的入站连接
- 限制 RPC 端口仅允许本地访问
- 临时禁用防火墙以便快速测试（测试环境）

## 4. 快速启动指南

### 4.1 一键启动
1. 下载 NogoChain 测试网络版本
2. 解压到指定目录
3. 运行 `testnet/start_node.bat` 脚本
4. 脚本会自动完成所有配置和启动步骤

### 4.2 手动启动步骤

#### 4.2.1 环境准备
1. 安装 Go 1.22+（如果未安装）
2. 克隆代码库：`git clone https://github.com/nogochain/nogochain.git`
3. 进入目录：`cd nogochain`
4. 安装依赖：`go mod tidy`

#### 4.2.2 编译组件
1. 编译节点：`go build -o build/nogochain cmd/nogochain/main.go`
2. 编译矿池：`go build -o build/nogopool cmd/nogopool/main.go`
3. 编译 CLI 工具：`go build -o build/nogocli cmd/nogocli/main.go`
4. 编译挖矿工具：`go build -o build/nogominer cmd/nogominer/main.go`

#### 4.2.3 配置文件设置
1. 复制测试网络配置：`cp -r testnet/config/ .`
2. 修改 `node_config.json` 中的配置参数
3. 修改 `mining_pool_config.json` 中的配置参数

#### 4.2.4 启动节点
1. 初始化创世区块：`./nogochain.exe --datadir testnet/data init testnet/config/genesis.json`
2. 启动节点：`./nogochain.exe --datadir testnet/data --config testnet/config/node_config.json`

#### 4.2.5 启动矿池
1. 启动矿池：`./nogopool.exe --datadir testnet/data --config testnet/config/mining_pool_config.json`

#### 4.2.6 启动矿工
1. 启动矿工：`./nogominer.exe --stratum 127.0.0.1:3336 --wallet 0xYourTestWalletAddress`

## 5. 测试方法和工具

### 5.1 基本功能测试

#### 5.1.1 节点同步测试
1. 启动多个节点
2. 监控节点同步状态
3. 验证区块高度一致
4. 测试网络分区和恢复

#### 5.1.2 交易测试
1. 生成测试账户
2. 发送交易
3. 验证交易确认
4. 测试高并发交易

#### 5.1.3 智能合约测试
1. 部署测试合约
2. 调用合约方法
3. 测试合约事件
4. 验证合约执行结果

### 5.2 性能测试

#### 5.2.1 区块处理性能
1. 测量区块生成时间
2. 测试区块验证速度
3. 评估交易处理能力
4. 分析系统资源使用

#### 5.2.2 网络性能
1. 测试节点间通信延迟
2. 评估网络带宽使用
3. 测量同步速度
4. 测试网络拓扑变化

#### 5.2.3 存储性能
1. 测量数据库读写速度
2. 评估存储增长趋势
3. 测试数据压缩效果
4. 分析存储优化空间

### 5.3 安全测试

#### 5.3.1 漏洞扫描
1. 运行安全扫描工具
2. 检查已知漏洞
3. 测试 RPC 接口安全
4. 评估网络协议安全性

#### 5.3.2 攻击模拟
1. 测试拒绝服务攻击
2. 模拟双花攻击
3. 测试 51% 攻击防御
4. 评估共识算法安全性

### 5.4 测试工具

#### 5.4.1 内置工具
- **nogocli**: 命令行工具，用于节点管理和查询
- **nogotest**: 自动化测试工具
- **tx_generator**: 交易生成器，用于性能测试

#### 5.4.2 第三方工具
- **Hardhat**: 智能合约开发和测试框架
- **Foundry**: 以太坊测试框架
- **Remix**: 在线智能合约编辑器和测试工具
- **Postman**: API 测试工具
- **Grafana**: 监控和可视化工具

## 6. 测试网络配置指南

### 6.1 配置文件说明

#### 6.1.1 节点配置文件 (node_config.json)
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

#### 6.1.2 矿池配置文件 (mining_pool_config.json)
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

### 6.2 环境变量配置

#### 6.2.1 Windows 环境
```batch
set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=off
set NOGOCHAIN_DATA_DIR=testnet/data
set NOGOCHAIN_LOG_LEVEL=debug
```

#### 6.2.2 Linux 环境
```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off
export NOGOCHAIN_DATA_DIR=testnet/data
export NOGOCHAIN_LOG_LEVEL=debug
```

### 6.3 开发模式配置

#### 6.3.1 快速开发设置
1. 启用 `dev_mode`：加速测试
2. 设置 `dev_block_time`：调整区块时间
3. 配置 `debug_api`：启用调试 API
4. 关闭 `metrics`：减少资源使用

#### 6.3.2 测试网络重置
1. 停止所有服务
2. 删除数据目录：`rm -rf testnet/data`
3. 重新初始化创世区块
4. 重启服务

## 7. 测试结果分析指南

### 7.1 日志分析

#### 7.1.1 节点日志
- **INFO**: 正常操作信息
- **DEBUG**: 详细调试信息
- **WARN**: 警告信息
- **ERROR**: 错误信息
- **FATAL**: 致命错误信息

#### 7.1.2 日志分析工具
- **grep**: 搜索特定关键词
- **tail**: 查看最新日志
- **logrotate**: 管理日志文件大小
- **ELK Stack**: 集中式日志分析

### 7.2 性能指标分析

#### 7.2.1 关键指标
- **TPS (Transactions Per Second)**: 每秒交易处理能力
- **Block Time**: 区块生成时间
- **Confirmation Time**: 交易确认时间
- **Resource Usage**: CPU、内存、磁盘、网络使用
- **Sync Speed**: 区块同步速度

#### 7.2.2 指标收集工具
- **Prometheus**: 指标收集和存储
- **Grafana**: 指标可视化和告警
- **top/htop**: 系统资源监控
- **iostat**: 磁盘 I/O 监控

### 7.3 测试报告模板

#### 7.3.1 基本信息
- 测试网络版本
- 测试时间和持续时间
- 测试环境配置
- 测试参与者

#### 7.3.2 测试结果
- 功能测试结果
- 性能测试结果
- 安全测试结果
- 问题和缺陷

#### 7.3.3 分析和建议
- 性能瓶颈分析
- 安全风险评估
- 优化建议
- 下一步测试计划

### 7.4 常见问题分析

#### 7.4.1 同步问题
- **症状**: 节点无法同步到最新区块
- **可能原因**: 网络连接问题、节点配置错误、区块数据损坏
- **解决方案**: 检查网络连接、验证配置文件、重置数据目录

#### 7.4.2 性能下降
- **症状**: TPS 降低、区块时间增加
- **可能原因**: 系统资源不足、网络拥堵、数据库性能问题
- **解决方案**: 增加系统资源、优化网络配置、调整数据库参数

#### 7.4.3 交易失败
- **症状**: 交易被拒绝或未确认
- **可能原因**: 余额不足、Gas 价格过低、合约执行失败
- **解决方案**: 检查账户余额、调整 Gas 价格、验证合约代码

#### 7.4.4 矿池连接问题
- **症状**: 矿工无法连接到矿池
- **可能原因**: 网络连接问题、端口配置错误、矿池服务未运行
- **解决方案**: 检查网络连接、验证端口配置、重启矿池服务

## 8. 测试网络维护

### 8.1 日常维护
- 监控节点状态
- 检查日志文件
- 备份测试数据
- 更新测试网络版本

### 8.2 网络重置
1. 提前通知测试网络用户
2. 停止所有服务
3. 清除测试数据
4. 部署新版本
5. 重新初始化网络
6. 通知用户网络已重置

### 8.3 版本升级
1. 下载最新测试版本
2. 停止测试服务
3. 替换二进制文件
4. 更新配置文件
5. 启动服务并验证

## 9. 开发者指南

### 9.1 智能合约开发
1. 安装开发工具：Hardhat、Foundry 或 Remix
2. 配置测试网络连接
3. 编写智能合约
4. 部署到测试网络
5. 测试合约功能

### 9.2 应用开发
1. 安装 NogoChain SDK
2. 配置 API 连接
3. 实现基本功能
4. 测试应用性能
5. 优化用户体验

### 9.3 贡献代码
1. Fork 代码库
2. 创建分支
3. 实现功能或修复问题
4. 运行测试
5. 提交 Pull Request

## 10. 测试网络激励

### 10.1 测试奖励
- 参与测试网络的用户可以获得测试代币奖励
- 发现安全漏洞的用户可以获得额外奖励
- 提供有价值反馈的用户可以获得贡献奖励

### 10.2 奖励申请
1. 注册测试网络账户
2. 参与测试活动
3. 提交测试结果和反馈
4. 验证贡献
5. 领取奖励

## 11. 常见问题解答

### 11.1 技术问题

#### 11.1.1 节点无法启动
**Q**: 节点启动失败，显示 "Error: Failed to load genesis file"
**A**: 检查创世区块文件路径是否正确，确保文件格式无误

#### 11.1.2 矿池连接失败
**Q**: 矿池无法连接到节点，显示 "Error: RPC connection failed"
**A**: 检查节点 RPC 端口是否正确配置，确保节点 RPC 服务正在运行

#### 11.1.3 交易未确认
**Q**: 交易提交后长时间未确认
**A**: 检查 Gas 价格是否足够，查看节点同步状态，确认网络是否正常

#### 11.1.4 性能测试结果异常
**Q**: 性能测试结果与预期不符
**A**: 检查测试环境配置，确保系统资源充足，验证测试方法是否正确

### 11.2 配置问题

#### 11.2.1 端口冲突
**Q**: 启动服务时显示 "Error: Address already in use"
**A**: 检查端口是否被占用，修改配置文件中的端口设置

#### 11.2.2 内存不足
**Q**: 运行时显示 "Error: Out of memory"
**A**: 增加系统内存，调整服务内存限制，优化配置参数

#### 11.2.3 磁盘空间不足
**Q**: 运行时显示 "Error: No space left on device"
**A**: 增加磁盘空间，启用数据压缩，清理旧数据

### 11.3 网络问题

#### 11.3.1 连接数限制
**Q**: 节点连接数达到上限
**A**: 调整 `max_peers` 参数，优化网络配置，确保网络带宽充足

#### 11.3.2 网络延迟
**Q**: 节点间通信延迟高
**A**: 检查网络连接质量，优化网络拓扑，使用更稳定的网络连接

## 12. 联系方式

### 12.1 开发者社区
- GitHub: https://github.com/nogochain
- Discord: https://discord.gg/nogochain
- Telegram: https://t.me/nogochain
- Twitter: https://twitter.com/nogochain

### 12.2 技术支持
- 开发论坛: https://dev.nogochain.org
- 技术文档: https://docs.nogochain.org
- 支持邮箱: dev@nogochain.org

### 12.3 测试网络监控
- 网络状态: https://monitor.nogochain.org
- 区块浏览器: https://testnet-explorer.nogochain.org
- 测试网络统计: https://stats.nogochain.org

## 13. 附录

### 13.1 测试网络启动脚本

#### 13.1.1 Windows 启动脚本 (start_testnet.bat)
```batch
@echo off

:: 设置环境变量
set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=off
set NOGOCHAIN_DATA_DIR=testnet/data
set NOGOCHAIN_LOG_LEVEL=debug

:: 创建目录
if not exist "testnet" mkdir "testnet"
if not exist "testnet/data" mkdir "testnet/data"
if not exist "testnet/logs" mkdir "testnet/logs"
if not exist "testnet/config" mkdir "testnet/config"

:: 复制配置文件
if not exist "testnet/config/genesis.json" copy "testnet/config/genesis.json" "testnet/config/genesis.json"
if not exist "testnet/config/node_config.json" copy "testnet/config/node_config.json" "testnet/config/node_config.json"
if not exist "testnet/config/mining_pool_config.json" copy "testnet/config/mining_pool_config.json" "testnet/config/mining_pool_config.json"

:: 启动节点
echo 启动测试网络节点...
start "NogoChain Testnet Node" /MIN nogochain.exe --datadir testnet/data --config testnet/config/node_config.json

:: 启动矿池
echo 启动测试网络矿池...
start "NogoChain Testnet Pool" /MIN nogopool.exe --datadir testnet/data --config testnet/config/mining_pool_config.json

echo 测试网络启动完成！
echo 查看日志文件了解运行状态：testnet/logs/
pause
```

#### 13.1.2 Linux 启动脚本 (start_testnet.sh)
```bash
#!/bin/bash

# 设置环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off
export NOGOCHAIN_DATA_DIR=testnet/data
export NOGOCHAIN_LOG_LEVEL=debug

# 创建目录
mkdir -p testnet/data testnet/logs testnet/config

# 复制配置文件
cp -n testnet/config/genesis.json testnet/config/genesis.json
cp -n testnet/config/node_config.json testnet/config/node_config.json
cp -n testnet/config/mining_pool_config.json testnet/config/mining_pool_config.json

# 启动节点
echo "启动测试网络节点..."
./nogochain --datadir testnet/data --config testnet/config/node_config.json > testnet/logs/node.log 2>&1 &

# 启动矿池
echo "启动测试网络矿池..."
./nogopool --datadir testnet/data --config testnet/config/mining_pool_config.json > testnet/logs/pool.log 2>&1 &

echo "测试网络启动完成！"
echo "查看日志文件了解运行状态：testnet/logs/"
```

### 13.2 测试工具命令

#### 13.2.1 交易生成器
```bash
# 生成 1000 个交易
./tx_generator --count 1000 --from 0xSenderAddress --to 0xReceiverAddress --value 1000000000000000000 --gas 21000 --gasprice 1000000000
```

#### 13.2.2 性能测试工具
```bash
# 运行性能测试
./performance_test --duration 3600 --tps 1000 --concurrency 100
```

#### 13.2.3 网络测试工具
```bash
# 测试网络延迟
./network_test --nodes 10 --messages 1000 --size 1024
```

### 13.3 测试网络资源

#### 13.3.1 官方资源
- 测试网络文档: https://docs.nogochain.org/testnet
- 测试代币 faucet: https://faucet.nogochain.org
- 测试网络状态: https://status.nogochain.org

#### 13.3.2 第三方资源
- Remix IDE: https://remix.ethereum.org
- Hardhat: https://hardhat.org
- Foundry: https://getfoundry.sh
- MetaMask: https://metamask.io

## 14. 版本历史

| 版本 | 日期 | 变更内容 |
|------|------|----------|
| v0.1.0 | 2024-01-01 | 测试网络初始版本 |
| v0.1.1 | 2024-01-15 | 性能优化 |
| v0.1.2 | 2024-02-01 | 安全更新 |
| v0.2.0 | 2024-03-01 | 功能增强 |
| v0.2.1 | 2024-03-15 | Bug 修复 |

## 15. 免责声明

本教程仅供测试和开发目的使用。测试网络中的代币没有实际价值，仅用于功能测试。部署过程中请确保遵循相关法律法规，注意安全防护措施。NogoChain 团队不对因使用本教程而导致的任何损失负责。

---

**文档版本**: v0.2.1
**最后更新**: 2024-03-15
**作者**: NogoChain 团队
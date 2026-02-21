# NogoChain 测试网络配置

本文档介绍了NogoChain测试网络的配置和使用方法。

## 测试网络特性

- **链ID**: 31888
- **代币符号**: TESTNOGO
- **目标区块时间**: 10秒
- **初始难度**: 1000
- **基础区块奖励**: 20 TESTNOGO
- **减半间隔**: 10,000个区块
- **最低区块奖励**: 0.5 TESTNOGO

## 目录结构

```
testnet/
├── config/                  # 配置文件目录
│   ├── testnet_params.go    # 测试网络参数
│   ├── genesis.json         # 创世区块配置
│   ├── node_config.json     # 节点配置
│   └── mining_pool_config.json # 矿池配置
├── logs/                    # 日志目录
├── data/                    # 区块链数据目录
├── start_node.bat           # 节点启动脚本
├── start_miner.bat          # 挖矿启动脚本
└── README.md                # 本文档
```

## 配置说明

### 1. 测试网络参数 (testnet_params.go)

包含测试网络的基础链参数、区块奖励参数、Gas参数和共识参数。

### 2. 创世区块配置 (genesis.json)

定义了测试网络的初始状态，包括：
- 链配置参数
- 初始账户余额
- 创世区块属性

### 3. 节点配置 (node_config.json)

包含节点运行所需的所有配置：
- P2P网络配置
- RPC服务配置
- 同步配置
- 挖矿配置
- 日志和监控配置

### 4. 矿池配置 (mining_pool_config.json)

配置测试网络矿池的参数：
- 矿池基本设置
- 难度调整参数
- 节点连接配置
- 支付设置
- 统计信息配置

## 使用方法

### 1. 启动节点

```bash
cd testnet
./start_node.bat
```

该脚本会：
- 创建必要的目录结构
- 初始化区块链数据（首次运行）
- 启动测试网络节点

### 2. 启动挖矿

在另一个终端中运行：

```bash
cd testnet
./start_miner.bat
```

### 3. 连接到测试网络

可以使用以下RPC URL连接到测试网络：
- `http://127.0.0.1:8546`

### 4. 测试账户

创世区块中预配置了以下测试账户（每个账户有大量TESTNOGO）：
- `0x71c7656ec7ab88b098defb751b7401b5f6d8976f`
- `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
- `0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc`

## 调整配置

### 修改区块奖励

编辑 `config/testnet_params.go` 文件中的 `TestBlockReward` 常量。

### 修改难度参数

编辑 `config/testnet_params.go` 文件中的难度相关参数：
- `TestInitialDifficulty` - 初始难度
- `TestTargetBlockTime` - 目标区块时间
- `TestDifficultyAdjustmentInterval` - 难度调整间隔

### 修改矿池设置

编辑 `config/mining_pool_config.json` 文件中的矿池参数。

## 监控

- **节点日志**: `testnet/logs/nogochain-testnet.log`
- **挖矿日志**: `testnet/logs/miner.log`
- **矿池日志**: `testnet/logs/pool.log`
- **监控指标**: `http://127.0.0.1:9091/metrics`

## 常见问题

### 1. 节点无法启动

- 检查端口是否被占用
- 检查配置文件格式是否正确
- 检查日志文件中的错误信息

### 2. 挖矿没有收益

- 确保节点已完全同步
- 检查矿工地址是否正确
- 检查网络连接是否正常

### 3. 矿池连接失败

- 确保节点RPC服务已启用
- 检查矿池配置中的RPC URL是否正确
- 检查防火墙设置

## 技术支持

如需技术支持，请参考项目根目录下的 `docs/` 目录中的文档，或联系开发团队。
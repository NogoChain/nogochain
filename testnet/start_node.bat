@echo off

rem 设置测试网络环境变量
set CHAIN_ID=31888
set NETWORK_ID=31888
set DATA_DIR=testnet/data
set LOG_DIR=testnet/logs
set GENESIS_FILE=testnet/config/genesis.json
set NODE_CONFIG=testnet/config/node_config.json

rem 创建必要的目录
if not exist "%DATA_DIR%" mkdir "%DATA_DIR%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

rem 初始化区块链数据
if not exist "%DATA_DIR%\chaindata" (
    echo 初始化测试网络区块链数据...
    ..\bin\nogod --datadir "%DATA_DIR%" init "%GENESIS_FILE%"
    if %errorlevel% neq 0 (
        echo 初始化失败，请检查创世区块配置文件
        pause
        exit /b 1
    )
    echo 初始化完成
)

echo 启动测试网络节点...
echo 节点配置: %NODE_CONFIG%
echo 数据目录: %DATA_DIR%
echo 日志目录: %LOG_DIR%
echo 按 Ctrl+C 停止节点

rem 启动节点
..\bin\nogod ^
    --datadir "%DATA_DIR%" ^
    --config "%NODE_CONFIG%" ^
    --networkid %NETWORK_ID% ^
    --port 30304 ^
    --rpc ^
    --rpcaddr 127.0.0.1 ^
    --rpcport 8546 ^
    --rpcapi eth,net,web3,nogo ^
    --nodiscover ^
    --verbosity 3 ^
    --metrics ^
    --metrics.addr 127.0.0.1 ^
    --metrics.port 9091

pause
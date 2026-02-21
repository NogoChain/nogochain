@echo off

rem 设置测试网络挖矿环境变量
set MINER_ADDRESS=0x71c7656ec7ab88b098defb751b7401b5f6d8976f
set RPC_URL=http://127.0.0.1:8546
set THREADS=2
set LOG_FILE=testnet/logs/miner.log

rem 创建日志目录
if not exist "testnet/logs" mkdir "testnet/logs"

echo 启动测试网络挖矿...
echo 矿工地址: %MINER_ADDRESS%
echo RPC URL: %RPC_URL%
echo 线程数: %THREADS%
echo 日志文件: %LOG_FILE%
echo 按 Ctrl+C 停止挖矿

rem 启动挖矿
..\bin\nogominer ^
    --address %MINER_ADDRESS% ^
    --rpc %RPC_URL% ^
    --threads %THREADS% ^
    --log %LOG_FILE% ^
    --verbosity 3

pause
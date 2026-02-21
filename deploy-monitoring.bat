@echo off

REM 部署Prometheus和Grafana监控系统

echo 部署NogoChain监控系统
echo ======================
echo.

REM 检查是否存在Prometheus和Grafana
if not exist "prometheus.exe" (
    echo 错误: Prometheus未找到
    echo 请先下载Prometheus: https://prometheus.io/download/
    echo 并将prometheus.exe放在当前目录
    pause
    exit /b 1
)

if not exist "grafana-server.exe" (
    echo 错误: Grafana未找到
    echo 请先下载Grafana: https://grafana.com/grafana/download
    echo 并将grafana-server.exe放在当前目录
    pause
    exit /b 1
)

REM 创建监控配置目录
mkdir "monitoring\prometheus" 2>nul
mkdir "monitoring\grafana\data" 2>nul
mkdir "monitoring\grafana\logs" 2>nul
mkdir "monitoring\grafana\plugins" 2>nul

REM 复制Prometheus配置文件
copy "prometheus.yml" "monitoring\prometheus\" /y

REM 启动Prometheus
echo 启动Prometheus...
start "Prometheus" "prometheus.exe" --config.file="monitoring\prometheus\prometheus.yml" --storage.tsdb.path="monitoring\prometheus\data"
sleep 5

echo Prometheus启动成功!

REM 启动Grafana
echo 启动Grafana...
start "Grafana" "grafana-server.exe" --homepath="." --config=".\conf\defaults.ini" --packaging=zip cfg:default.paths.provisioning=".\conf\provisioning" cfg:default.paths.data="monitoring\grafana\data" cfg:default.paths.logs="monitoring\grafana\logs" cfg:default.paths.plugins="monitoring\grafana\plugins"
sleep 5

echo Grafana启动成功!

echo 监控系统部署完成!
echo ======================
echo Prometheus地址: http://localhost:9090
echo Grafana地址: http://localhost:3000
echo 默认用户名/密码: admin/admin
echo.
echo 下一步操作:
echo 1. 登录Grafana (http://localhost:3000)
echo 2. 添加Prometheus数据源 (URL: http://localhost:9090)
echo 3. 导入仪表板: 上传 grafana-dashboard.json 文件
echo 4. 启动NogoChain节点和矿池
echo.
echo 监控指标将在节点和矿池启动后开始采集
echo.
pause

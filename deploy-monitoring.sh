#!/bin/bash

# 部署Prometheus和Grafana监控系统

echo "部署NogoChain监控系统"
echo "======================"

# 检查是否安装了Prometheus和Grafana
if ! command -v prometheus &> /dev/null; then
    echo "错误: Prometheus未安装"
    echo "请先安装Prometheus: https://prometheus.io/download/"
    exit 1
fi

if ! command -v grafana-server &> /dev/null; then
    echo "错误: Grafana未安装"
    echo "请先安装Grafana: https://grafana.com/grafana/download"
    exit 1
fi

# 创建监控配置目录
mkdir -p monitoring/{prometheus,grafana}

# 复制Prometheus配置文件
cp prometheus.yml monitoring/prometheus/

# 启动Prometheus
echo "启动Prometheus..."
prometheus --config.file=monitoring/prometheus/prometheus.yml --storage.tsdb.path=monitoring/prometheus/data &
prometheus_pid=$!
sleep 5

# 检查Prometheus是否启动成功
if ! kill -0 $prometheus_pid 2> /dev/null; then
    echo "错误: Prometheus启动失败"
    exit 1
fi
echo "Prometheus启动成功!"

# 启动Grafana
echo "启动Grafana..."
grafana-server --homepath=/usr/share/grafana --config=/etc/grafana/grafana.ini --packaging=deb cfg:default.paths.provisioning=/etc/grafana/provisioning cfg:default.paths.data=monitoring/grafana/data cfg:default.paths.logs=monitoring/grafana/logs cfg:default.paths.plugins=monitoring/grafana/plugins &
grafana_pid=$!
sleep 5

# 检查Grafana是否启动成功
if ! kill -0 $grafana_pid 2> /dev/null; then
    echo "错误: Grafana启动失败"
    exit 1
fi
echo "Grafana启动成功!"

echo "监控系统部署完成!"
echo "======================"
echo "Prometheus地址: http://localhost:9090"
echo "Grafana地址: http://localhost:3000"
echo "默认用户名/密码: admin/admin"
echo ""
echo "下一步操作:"
echo "1. 登录Grafana (http://localhost:3000)"
echo "2. 添加Prometheus数据源 (URL: http://localhost:9090)"
echo "3. 导入仪表板: 上传 grafana-dashboard.json 文件"
echo "4. 启动NogoChain节点和矿池"
echo ""
echo "监控指标将在节点和矿池启动后开始采集"

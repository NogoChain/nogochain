#!/bin/bash

set -e

# 脚本版本
SCRIPT_VERSION="1.0.0"

# 输出目录
OUTPUT_DIR="bin"

# 编译组件列表
COMPONENTS=(
    "nogochain"
    "nogocli"
    "nogod"
    "nogominer"
    "nogopool"
    "nogostratumminer"
    "nogotest"
)

# 显示帮助信息
show_help() {
    echo "NogoChain 编译脚本 v$SCRIPT_VERSION"
    echo ""
    echo "用法: ./build.sh [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help           显示帮助信息"
    echo "  -v, --version        显示脚本版本"
    echo "  -o, --output <dir>   设置输出目录 (默认: $OUTPUT_DIR)"
    echo "  -c, --component <name>  仅编译指定组件"
    echo "  -a, --all            编译所有组件 (默认)"
    echo "  -d, --debug          启用调试模式"
    echo ""
    echo "示例:"
    echo "  ./build.sh              # 编译所有组件到默认目录"
    echo "  ./build.sh -o build     # 编译所有组件到 build 目录"
    echo "  ./build.sh -c nogochain # 仅编译 nogochain 组件"
}

# 显示版本信息
show_version() {
    echo "NogoChain 编译脚本 v$SCRIPT_VERSION"
}

# 验证依赖
verify_dependencies() {
    echo "验证依赖..."
    
    # 检查 Go 是否已安装
    if ! command -v go &> /dev/null; then
        echo "错误: Go 未安装，请安装 Go 1.22+"
        exit 1
    fi
    
    # 检查 Go 版本
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo "Go 版本: $GO_VERSION"
    
    # 检查 CGO_ENABLED
    CGO_STATUS=$(go env CGO_ENABLED)
    echo "CGO_ENABLED: $CGO_STATUS"
    if [ "$CGO_STATUS" != "0" ]; then
        echo "警告: CGO 已启用，请设置 CGO_ENABLED=0"
    fi
    
    # 验证 go.mod
    echo "验证 go.mod..."
    go mod verify
    
    echo "依赖验证成功!"
}

# 创建输出目录
create_output_dir() {
    if [ ! -d "$OUTPUT_DIR" ]; then
        echo "创建输出目录: $OUTPUT_DIR"
        mkdir -p "$OUTPUT_DIR"
    fi
}

# 编译单个组件
build_component() {
    local component=$1
    echo "编译组件: $component"
    
    local cmd_path="cmd/$component"
    local output_path="$OUTPUT_DIR/$component"
    
    if [ ! -d "$cmd_path" ]; then
        echo "错误: 组件目录不存在: $cmd_path"
        return 1
    fi
    
    # 编译组件
    if go build -o "$output_path" "$cmd_path"; then
        echo "✓ 组件编译成功: $component"
        return 0
    else
        echo "✗ 组件编译失败: $component"
        return 1
    fi
}

# 编译所有组件
build_all_components() {
    local success=0
    
    for component in "${COMPONENTS[@]}"; do
        if ! build_component "$component"; then
            success=1
        fi
    done
    
    return $success
}

# 主函数
main() {
    local target_component=""
    local debug_mode=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                show_version
                exit 0
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -c|--component)
                target_component="$2"
                shift 2
                ;;
            -a|--all)
                target_component=""
                shift
                ;;
            -d|--debug)
                debug_mode=true
                shift
                ;;
            *)
                echo "错误: 未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 启用调试模式
    if [ "$debug_mode" = true ]; then
        set -x
    fi
    
    # 验证依赖
    verify_dependencies
    
    # 创建输出目录
    create_output_dir
    
    # 编译组件
    if [ -n "$target_component" ]; then
        # 仅编译指定组件
        if build_component "$target_component"; then
            echo ""
            echo "编译完成! 输出目录: $OUTPUT_DIR"
            exit 0
        else
            echo ""
            echo "编译失败!"
            exit 1
        fi
    else
        # 编译所有组件
        if build_all_components; then
            echo ""
            echo "所有组件编译完成! 输出目录: $OUTPUT_DIR"
            exit 0
        else
            echo ""
            echo "部分组件编译失败!"
            exit 1
        fi
    fi
}

# 执行主函数
main "$@"
# NogoChain 编译脚本
# Version: 1.0.0
# Author: NogoChain Team

# 默认参数
$OutputDir = "bin"
$Component = $null
$Debug = $false

# 解析命令行参数
for ($i=0; $i -lt $args.Length; $i++) {
    switch ($args[$i]) {
        "-h" {
            Write-Host "NogoChain 编译脚本 v1.0.0"
            Write-Host ""
            Write-Host "用法: .\build.ps1 [选项]"
            Write-Host ""
            Write-Host "选项:"
            Write-Host "  -h, --help           显示帮助信息"
            Write-Host "  -v, --version        显示脚本版本"
            Write-Host "  -o, --output <dir>   设置输出目录 (默认: bin)"
            Write-Host "  -c, --component <name>  仅编译指定组件"
            Write-Host "  -d, --debug          启用调试模式"
            Write-Host ""
            Write-Host "示例:"
            Write-Host "  .\build.ps1              # 编译所有组件到默认目录"
            Write-Host "  .\build.ps1 -o build     # 编译所有组件到 build 目录"
            Write-Host "  .\build.ps1 -c nogochain # 仅编译 nogochain 组件"
            exit 0
        }
        "--help" {
            Write-Host "NogoChain 编译脚本 v1.0.0"
            Write-Host ""
            Write-Host "用法: .\build.ps1 [选项]"
            Write-Host ""
            Write-Host "选项:"
            Write-Host "  -h, --help           显示帮助信息"
            Write-Host "  -v, --version        显示脚本版本"
            Write-Host "  -o, --output <dir>   设置输出目录 (默认: bin)"
            Write-Host "  -c, --component <name>  仅编译指定组件"
            Write-Host "  -d, --debug          启用调试模式"
            Write-Host ""
            Write-Host "示例:"
            Write-Host "  .\build.ps1              # 编译所有组件到默认目录"
            Write-Host "  .\build.ps1 -o build     # 编译所有组件到 build 目录"
            Write-Host "  .\build.ps1 -c nogochain # 仅编译 nogochain 组件"
            exit 0
        }
        "-v" {
            Write-Host "NogoChain 编译脚本 v1.0.0"
            exit 0
        }
        "--version" {
            Write-Host "NogoChain 编译脚本 v1.0.0"
            exit 0
        }
        "-o" {
            $i++
            if ($i -lt $args.Length) {
                $OutputDir = $args[$i]
            }
        }
        "--output" {
            $i++
            if ($i -lt $args.Length) {
                $OutputDir = $args[$i]
            }
        }
        "-c" {
            $i++
            if ($i -lt $args.Length) {
                $Component = $args[$i]
            }
        }
        "--component" {
            $i++
            if ($i -lt $args.Length) {
                $Component = $args[$i]
            }
        }
        "-d" {
            $Debug = $true
        }
        "--debug" {
            $Debug = $true
        }
    }
}

# 启用调试模式
if ($Debug) {
    $DebugPreference = "Continue"
}

# 验证依赖
Write-Host "验证依赖..."

# 检查 Go 是否已安装
if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Host "错误: Go 未安装，请安装 Go 1.22+"
    exit 1
}

# 检查 Go 版本
$goVersionOutput = go version
Write-Host "Go 版本: $goVersionOutput"

# 检查 CGO_ENABLED
$cgoStatus = go env CGO_ENABLED
Write-Host "CGO_ENABLED: $cgoStatus"
if ($cgoStatus -ne "0") {
    Write-Host "警告: CGO 已启用，请设置 CGO_ENABLED=0"
}

# 验证 go.mod
Write-Host "验证 go.mod..."
go mod verify
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: go.mod 验证失败"
    exit 1
}

Write-Host "依赖验证成功!"

# 创建输出目录
if (-not (Test-Path -Path $OutputDir -PathType Container)) {
    Write-Host "创建输出目录: $OutputDir"
    New-Item -Path $OutputDir -ItemType Directory -Force | Out-Null
}

# 定义组件列表
$Components = @(
    "nogochain",
    "nogocli",
    "nogod",
    "nogominer",
    "nogopool",
    "nogostratumminer",
    "nogotest"
)

# 编译指定组件
if ($Component) {
    Write-Host "编译组件: $Component"
    $CmdPath = "cmd\$Component"
    $OutputPath = "$OutputDir\$Component.exe"
    
    if (-not (Test-Path -Path $CmdPath -PathType Container)) {
        Write-Host "错误: 组件目录不存在: $CmdPath"
        exit 1
    }
    
    go build -o "$OutputPath" "$CmdPath"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "组件编译成功: $Component"
        Write-Host "编译完成! 输出目录: $OutputDir"
        exit 0
    } else {
        Write-Host "组件编译失败: $Component"
        exit 1
    }
}

# 编译所有组件
Write-Host "编译所有组件..."
$Success = $true

foreach ($c in $Components) {
    Write-Host "编译组件: $c"
    $CmdPath = "cmd\$c"
    $OutputPath = "$OutputDir\$c.exe"
    
    if (-not (Test-Path -Path $CmdPath -PathType Container)) {
        Write-Host "错误: 组件目录不存在: $CmdPath"
        $Success = $false
    } else {
        go build -o "$OutputPath" "$CmdPath"
        if ($LASTEXITCODE -eq 0) {
            Write-Host "组件编译成功: $c"
        } else {
            Write-Host "组件编译失败: $c"
            $Success = $false
        }
    }
}

if ($Success) {
    Write-Host "所有组件编译完成! 输出目录: $OutputDir"
    exit 0
} else {
    Write-Host "部分组件编译失败!"
    exit 1
}
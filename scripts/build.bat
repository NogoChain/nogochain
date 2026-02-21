@echo off
REM NogoChain 编译脚本
REM Version: 1.0.0
REM Author: NogoChain Team

setlocal

REM 默认输出目录
set OUTPUT_DIR=bin

REM 显示帮助信息
if "%1"=="-h" goto help
if "%1"=="--help" goto help
if "%1"=="/h" goto help

REM 显示版本信息
if "%1"=="-v" goto version
if "%1"=="--version" goto version
if "%1"=="/v" goto version

REM 设置输出目录
if "%1"=="-o" (set OUTPUT_DIR=%2 && shift && shift) 
if "%1"=="--output" (set OUTPUT_DIR=%2 && shift && shift)
if "%1"=="/o" (set OUTPUT_DIR=%2 && shift && shift)

REM 仅编译指定组件
set TARGET_COMPONENT=
if "%1"=="-c" (set TARGET_COMPONENT=%2 && shift && shift)
if "%1"=="--component" (set TARGET_COMPONENT=%2 && shift && shift)
if "%1"=="/c" (set TARGET_COMPONENT=%2 && shift && shift)

REM 验证依赖
echo 验证依赖...

REM 检查 Go 是否已安装
go version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo 错误: Go 未安装，请安装 Go 1.22+
    exit /b 1
)

REM 检查 Go 版本
echo Go 版本: 
for /f "tokens=3" %%i in ('go version') do echo %%i

REM 检查 CGO_ENABLED
echo CGO_ENABLED: 
for /f %%i in ('go env CGO_ENABLED') do echo %%i

REM 验证 go.mod
echo 验证 go.mod...
go mod verify
if %ERRORLEVEL% neq 0 (
    echo 错误: go.mod 验证失败
    exit /b 1
)

echo 依赖验证成功!

REM 创建输出目录
if not exist "%OUTPUT_DIR%" (
    echo 创建输出目录: %OUTPUT_DIR%
    mkdir "%OUTPUT_DIR%"
)

REM 定义组件列表
set COMPONENTS=nogochain nogocli nogod nogominer nogopool nogostratumminer nogotest

REM 编译指定组件
if not "%TARGET_COMPONENT%"=="" (
    echo 编译组件: %TARGET_COMPONENT%
    set CMD_PATH=cmd\%TARGET_COMPONENT%
    set OUTPUT_PATH=%OUTPUT_DIR%\%TARGET_COMPONENT%.exe
    
    if not exist "%CMD_PATH%" (
        echo 错误: 组件目录不存在: %CMD_PATH%
        exit /b 1
    )
    
    go build -o "%OUTPUT_PATH%" "%CMD_PATH%"
    if %ERRORLEVEL% equ 0 (
        echo 组件编译成功: %TARGET_COMPONENT%
        echo 编译完成! 输出目录: %OUTPUT_DIR%
        exit /b 0
    ) else (
        echo 组件编译失败: %TARGET_COMPONENT%
        exit /b 1
    )
)

REM 编译所有组件
echo 编译所有组件...
set SUCCESS=1

for %%c in (%COMPONENTS%) do (
    echo 编译组件: %%c
    set CMD_PATH=cmd\%%c
    set OUTPUT_PATH=%OUTPUT_DIR%\%%c.exe
    
    if not exist "%CMD_PATH%" (
        echo 错误: 组件目录不存在: %CMD_PATH%
        set SUCCESS=0
    ) else (
        go build -o "%OUTPUT_PATH%" "%CMD_PATH%"
        if %ERRORLEVEL% equ 0 (
            echo 组件编译成功: %%c
        ) else (
            echo 组件编译失败: %%c
            set SUCCESS=0
        )
    )
)

if %SUCCESS% equ 1 (
    echo 所有组件编译完成! 输出目录: %OUTPUT_DIR%
    exit /b 0
) else (
    echo 部分组件编译失败!
    exit /b 1
)

:help
echo NogoChain 编译脚本 v1.0.0
echo
echo 用法: build.bat [选项]
echo
echo 选项:
echo   -h, --help           显示帮助信息
echo   -v, --version        显示脚本版本
echo   -o, --output <dir>   设置输出目录 (默认: bin)
echo   -c, --component <name>  仅编译指定组件
echo
echo 示例:
echo   build.bat              # 编译所有组件到默认目录
echo   build.bat -o build     # 编译所有组件到 build 目录
echo   build.bat -c nogochain # 仅编译 nogochain 组件
exit /b 0

:version
echo NogoChain 编译脚本 v1.0.0
exit /b 0

endlocal
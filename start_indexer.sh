#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Indexer 启动脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 1. 编译 indexer
echo -e "${YELLOW}[1/3] 正在编译 indexer...${NC}"
go build -o bin/indexer ./indexer/main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}编译失败！${NC}"
    exit 1
fi
echo -e "${GREEN}编译成功！${NC}"
echo ""

# 2. 选择环境
echo -e "${YELLOW}[2/3] 请选择环境:${NC}"
echo "  1) testnet"
echo "  2) mainnet"
read -p "请输入选项 (1 或 2): " env_choice

case $env_choice in
    1)
        ENV_DIR="config/testnet"
        ENV_NAME="testnet"
        ;;
    2)
        ENV_DIR="config/mainnet"
        ENV_NAME="mainnet"
        ;;
    *)
        echo -e "${RED}无效选项，退出${NC}"
        exit 1
        ;;
esac

# 检查目录是否存在
if [ ! -d "$ENV_DIR" ]; then
    echo -e "${RED}目录 $ENV_DIR 不存在！${NC}"
    exit 1
fi

# 3. 查找所有 yaml 配置文件（排除 config.yaml）
echo -e "${YELLOW}[3/3] 扫描配置文件...${NC}"
CONFIG_FILES=($(find "$ENV_DIR" -name "*.yaml" -type f ! -name "config.yaml" | sort))

if [ ${#CONFIG_FILES[@]} -eq 0 ]; then
    echo -e "${RED}在 $ENV_DIR 目录下未找到配置文件！${NC}"
    exit 1
fi

echo -e "${GREEN}找到 ${#CONFIG_FILES[@]} 个配置文件:${NC}"
for i in "${!CONFIG_FILES[@]}"; do
    filename=$(basename "${CONFIG_FILES[$i]}")
    echo "  $((i+1)). $filename"
done
echo ""

# 4. 启动 indexer 进程
echo -e "${BLUE}开始启动 indexer 进程...${NC}"
echo ""

# 创建日志目录
mkdir -p log/indexer

# 存储进程 PID
PIDS=()

for config_file in "${CONFIG_FILES[@]}"; do
    filename=$(basename "$config_file" .yaml)
    echo -e "${YELLOW}启动: $filename${NC}"
    
    # 在后台启动 indexer
    ./bin/indexer -f "$config_file" > "log/indexer/${filename}.log" 2>&1 &
    PID=$!
    PIDS+=($PID)
    
    echo -e "${GREEN}  ✓ 已启动 (PID: $PID)${NC}"
    echo "  配置文件: $config_file"
    echo "  日志文件: log/indexer/${filename}.log"
    echo ""
    
    # 稍微延迟，避免同时启动太多进程
    sleep 1
done

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}所有 indexer 进程已启动！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "运行中的进程:"
for i in "${!PIDS[@]}"; do
    filename=$(basename "${CONFIG_FILES[$i]}" .yaml)
    echo "  PID ${PIDS[$i]}: $filename"
done
echo ""
echo "查看日志:"
for config_file in "${CONFIG_FILES[@]}"; do
    filename=$(basename "$config_file" .yaml)
    echo "  tail -f log/indexer/${filename}.log"
done
echo ""
echo "停止所有进程:"
echo "  pkill -f './bin/indexer'"
echo "  或: kill ${PIDS[*]}"
echo ""


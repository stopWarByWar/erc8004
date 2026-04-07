#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 安全选项
set -euo pipefail

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

get_yaml_top_value() {
  # Get top-level "key: value" from yaml (first match)
  # Usage: get_yaml_top_value <file> <key>
  local file="$1"
  local key="$2"
  awk -v k="$key" '
    $0 ~ "^[[:space:]]*"k":[[:space:]]*" {
      sub("^[[:space:]]*"k":[[:space:]]*", "", $0)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", $0)
      print $0
      exit
    }
  ' "$file"
}

get_yaml_section_value() {
  # Get 2-space-indented value within a top-level section.
  # Usage: get_yaml_section_value <file> <section> <key>
  local file="$1"
  local section="$2"
  local key="$3"
  awk -v s="$section" -v k="$key" '
    BEGIN { in_section=0 }
    $0 ~ "^[[:space:]]*"s":[[:space:]]*$" { in_section=1; next }
    in_section==1 && $0 ~ "^[^[:space:]]" { in_section=0 }
    in_section==1 && $0 ~ "^[[:space:]]{2}"k":[[:space:]]*" {
      sub("^[[:space:]]{2}"k":[[:space:]]*", "", $0)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", $0)
      print $0
      exit
    }
  ' "$file"
}

has_any_processor_enabled() {
  local file="$1"
  local r i v c m
  r="$(get_yaml_section_value "$file" "reputation" "run" || true)"
  i="$(get_yaml_section_value "$file" "identity" "run" || true)"
  v="$(get_yaml_section_value "$file" "validation" "run" || true)"
  c="$(get_yaml_section_value "$file" "commerce" "run" || true)"
  m="$(get_yaml_section_value "$file" "comment" "run" || true)"
  [[ "$r" == "true" || "$i" == "true" || "$v" == "true" || "$c" == "true" || "$m" == "true" ]]
}

describe_enabled_processors() {
  local file="$1"
  local out=()
  local r i v c m
  r="$(get_yaml_section_value "$file" "reputation" "run" || true)"
  i="$(get_yaml_section_value "$file" "identity" "run" || true)"
  v="$(get_yaml_section_value "$file" "validation" "run" || true)"
  c="$(get_yaml_section_value "$file" "commerce" "run" || true)"
  m="$(get_yaml_section_value "$file" "comment" "run" || true)"

  [[ "$i" == "true" ]] && out+=("identity")
  [[ "$r" == "true" ]] && out+=("reputation")
  [[ "$v" == "true" ]] && out+=("validation")
  [[ "$c" == "true" ]] && out+=("commerce")
  [[ "$m" == "true" ]] && out+=("comment")

  if [[ ${#out[@]} -eq 0 ]]; then
    echo "none"
  else
    local IFS=","
    echo "${out[*]}"
  fi
}

check_commerce_db_schema_if_needed() {
  local file="$1"
  local commerce_run
  commerce_run="$(get_yaml_section_value "$file" "commerce" "run" || true)"
  if [[ "$commerce_run" != "true" ]]; then
    return 0
  fi

  if ! command -v psql >/dev/null 2>&1; then
    echo -e "${YELLOW}  ! psql 未安装，跳过 DB 缺表检查（commerce.run=true）。${NC}"
    return 0
  fi

  local dns
  dns="$(get_yaml_top_value "$file" "dns" || true)"
  if [[ -z "${dns:-}" ]]; then
    echo -e "${YELLOW}  ! 未配置 dns，跳过 DB 缺表检查（commerce.run=true）。${NC}"
    return 0
  fi

  # Connectivity check
  if ! psql "$dns" -v ON_ERROR_STOP=1 -qtAc "select 1" >/dev/null 2>&1; then
    echo -e "${RED}  ✗ DB 连接失败（commerce.run=true）。请检查 dns：${dns}${NC}"
    return 1
  fi

  # Schema check (commerce tables)
  local reg
  reg="$(psql "$dns" -v ON_ERROR_STOP=1 -qtAc "select to_regclass('public.commerce_actions')" 2>/dev/null | tr -d '[:space:]' || true)"
  if [[ "$reg" != "commerce_actions" ]]; then
    echo -e "${RED}  ✗ 缺少表 public.commerce_actions（commerce.run=true）。${NC}"
    echo -e "${YELLOW}    请先执行安全建表：${NC}"
    echo -e "${YELLOW}    psql \"${dns}\" -f migrations/202604071600_commerce_init_safe.psql${NC}"
    return 1
  fi
  return 0
}

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
    cfg_name="$(get_yaml_top_value "${CONFIG_FILES[$i]}" "name" || true)"
    cfg_name="${cfg_name:-$filename}"
    enabled="$(describe_enabled_processors "${CONFIG_FILES[$i]}")"
    echo "  $((i+1)). $filename  (${cfg_name})  [run: ${enabled}]"
done
echo ""

# 3.1 选择性启动（默认全起）
echo -e "${YELLOW}启动模式：默认全起；可选输入要启动的编号（逗号分隔），或直接回车。${NC}"
read -p "请输入要启动的配置编号（例如 1,3,5）或 all（默认 all）: " selection
selection="${selection:-}"

SELECTED_FILES=()
if [[ -z "$selection" || "$selection" == "all" || "$selection" == "ALL" ]]; then
  SELECTED_FILES=("${CONFIG_FILES[@]}")
else
  # Normalize separators to space
  selection="${selection//,/ }"
  for idx in $selection; do
    if ! [[ "$idx" =~ ^[0-9]+$ ]]; then
      echo -e "${RED}无效编号：$idx${NC}"
      exit 1
    fi
    if (( idx < 1 || idx > ${#CONFIG_FILES[@]} )); then
      echo -e "${RED}编号越界：$idx${NC}"
      exit 1
    fi
    SELECTED_FILES+=("${CONFIG_FILES[$((idx-1))]}")
  done
fi

# 4. 启动 indexer 进程
echo -e "${BLUE}开始启动 indexer 进程...${NC}"
echo ""

# 创建日志目录
mkdir -p log/indexer

# 存储进程 PID
PIDS=()

for config_file in "${SELECTED_FILES[@]}"; do
    filename=$(basename "$config_file" .yaml)
    echo -e "${YELLOW}启动: $filename${NC}"

    # 如果配置里所有 processor 都是 run=false，跳过以免启动空转进程
    if ! has_any_processor_enabled "$config_file"; then
        echo -e "${YELLOW}  ! 跳过：该配置中所有 processor 均为 run=false${NC}"
        echo ""
        continue
    fi

    # 若启用 commerce，启动前检查 DB 是否已建表
    if ! check_commerce_db_schema_if_needed "$config_file"; then
        echo -e "${RED}  ✗ 跳过：DB 缺表或不可用（commerce.run=true）${NC}"
        echo ""
        continue
    fi
    
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
    filename=$(basename "${SELECTED_FILES[$i]:-unknown}" .yaml)
    echo "  PID ${PIDS[$i]}: $filename"
done
echo ""
echo "查看日志:"
for config_file in "${SELECTED_FILES[@]}"; do
    filename=$(basename "$config_file" .yaml)
    echo "  tail -f log/indexer/${filename}.log"
done
echo ""
echo "停止所有进程:"
echo "  pkill -f './bin/indexer'"
echo "  或: kill ${PIDS[*]}"
echo ""


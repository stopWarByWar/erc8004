# ERC-8004 BAS (Blockchain Agent Service)

ERC-8004 标准的区块链 AI Agent 身份注册与验证系统。索引区块链智能合约事件，存储 Agent 身份/信誉数据于 PostgreSQL，并提供 REST API。

## 项目结构

```
├── abi/                      # Solidity 智能合约 ABI 绑定
├── agentCard/                # Agent Profile Card 数据结构
├── cmd/                      # CLI 工具
│   └── commerce_seed/        # 商业行为测试种子工具
├── config/                   # 配置文件
├── docs/
│   ├── designs/              # 设计文档
│   ├── specs/                # 规格说明
│   ├── exec-plans/           # 执行计划跟踪
│   └── erc-8183/             # ERC-8183 规范
├── helper/                   # AWS S3 上传辅助
├── indexer/                  # 区块链事件索引器
│   └── processor/            # 事件处理器（Commerce/Identity/Reputation/Validation）
├── logger/                   # 结构化日志
├── migrations/               # 数据库迁移
├── model/                    # GORM 数据模型
└── server/                   # Gin REST API
    ├── api/
    │   ├── handle/            # HTTP 处理器
    │   ├── logic/             # 业务逻辑
    │   └── types/             # API 类型定义
    └── main.go               # API 服务入口
```

## 快速开始

### 编译

```bash
# API 服务
go build ./server/main.go

# Indexer 服务
go build ./indexer/main.go

# CLI 工具
go build ./cmd/commerce_seed/main.go
go build ./cmd/update_desc_vector.go
```

### 运行

```bash
# API 服务
go run ./server/main.go -f ./server/api/logic/config.yaml

# Indexer 服务
go run ./indexer/main.go -f ./config/conf.yaml
./start_indexer.sh   # 交互式多链索引器启动器

# 批量更新向量嵌入
go run ./cmd/update_desc_vector.go -f ./config/conf.yaml
```

### 测试

```bash
go test ./...
go test -v -run TestName ./path/to/package/
go test -cover ./...
go test -race ./...
```

## ERC-8183 Commerce Token-Aware Budget

商业行为记录现在支持 ERC-20 token 感知：

### 新增字段

**`commerce_actions` 表：**
- `payment_token` — ERC-20 代币地址（从 `BudgetSet` 事件解析）
- `token_symbol` — 代币符号（从链上 `symbol()` 读取并缓存）
- `budget_usd` — BudgetSet 时刻的 USD 快照价值

**`commerce_scores` / `commerce_scores_global` 表：**
- `total_volume_usd` — USD 计价总交易量
- `weighted_volume_usd_sum` — Σ(budget_usd × outcome_score)

### Indexer 功能

- **Token 缓存**：相同 token 只调用一次 `symbol()` + `decimals()`
- **Job Budget 缓存**：BudgetSet 写入缓存，后续 funded/submitted/终态 事件复用
- **Price 快照**：BudgetSet 时刻计算 budget_usd，后续不变
- **CoinGecko API**：3 级降级（Pro → Free → 默认 1）

### CLI 参数

```bash
# commerce_seed 新增 flags
-payment-token string      # ERC-20 token 地址（默认 ETH）
-provider-agent-id string  # createJob 时传入 provider agent ID（默认 0）
```

### API 扩展

**GET `/agent/commerce/actions`** 新增过滤参数：
- `payment_token` — 按 ERC-20 地址精确过滤
- `token_symbol` — 按代币符号过滤（USDC/ETH/WETH）
- `min_budget_usd` — 最小 USD 价值
- `max_budget_usd` — 最大 USD 价值

**GET `/agent/commerce/scores`** 新增返回字段：
- `total_volume_usd` — USD 计价总交易量
- `weighted_score_usd` — Σ(budget_usd × outcome_score) / Σ(budget_usd)

## 架构

### 双服务模式

- **Server** (`server/main.go`)：Gin REST API
- **Indexer** (`indexer/main.go`)：轮询区块链事件

### 数据流

```
HTTP Request → handle/ → logic/ → model/ → PostgreSQL
Chain Event  → indexer/processor/ → model/ → PostgreSQL
```

### 多链支持

所有 Agent 以 `(chain_id, identity_registry)` 为作用域。链配置在 `config/testnet/` 和 `config/mainnet/`。

## 相关文档

- [CLAUDE.md](./CLAUDE.md) — 完整项目指南
- [ERC-8183 Commerce Reputation 设计](../docs/designs/erc8183-commerce-reputation.md)
- [Token-Aware Budget 规格](../docs/specs/erc8183-commerce-token-aware.md)
- [执行计划](../docs/exec-plans/erc8183-commerce-token-aware.md)

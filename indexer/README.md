# Indexer

区块链事件索引器，监听智能合约事件并将数据写入 PostgreSQL。

## 架构

```
┌─────────────────────────────────────────────────────────────┐
│                      indexer/main.go                        │
│                  多链配置加载 + Fan-out 模式                 │
└────────────────┬──────────────────────────────────────────┘
                 │
        ┌────────┴────────┬──────────────┐
        ▼                ▼              ▼
┌───────────────┐ ┌──────────────┐ ┌───────────────┐
│ Identity      │ │ Reputation   │ │ Commerce      │
│ Processor     │ │ Processor    │ │ Processor     │
│               │ │              │ │               │
│ • Registered  │ │ • NewFeedback│ │ • JobCreated  │
│ • UriUpdated   │ │ • Response   │ │ • JobFunded   │
│ • TransferOwn  │ │ • Feedback   │ │ • JobSubmitted│
│ • SetMetaData   │ │   Revoked    │ │ • JobCompleted│
│               │ │              │ │ • JobRejected │
│               │ │              │ │ • JobExpired  │
│               │ │              │ │ • BudgetSet   │
└───────────────┘ └──────────────┘ └───────────────┘
        │                │              │
        └────────────────┴──────────────┘
                           │
                    model/ → PostgreSQL
```

## Processor 类型

| Processor | 合约 | 事件 |
|-----------|------|------|
| `IdentityProcessor` | IdentityRegistry | Registered, UriUpdated, TransferOwnership, SetMetaData |
| `ReputationProcessor` | ReputationRegistry | NewFeedback, ResponseAppended, FeedbackRevoked |
| `ValidationRegistryProcessor` | ValidationRegistry | ValidationCreated, ValidationResponsed |
| `CommerceProcessor` | AgenticCommerce | JobCreated, JobFunded, JobSubmitted, JobCompleted, JobRejected, JobExpired, ProviderSet, BudgetSet |

## 事件处理流程

1. **轮询**：每个 processor 定时轮询区块，范围 `[execBlock+1, execBlock+fetchBlockInterval]`
2. **过滤**：调用 `ethclient.FilterLogs` 按合约地址和事件 Topic 过滤
3. **去重**：按 `blockNumber + txIndex` 跳过已处理事件
4. **写入**：解析事件数据，写入 PostgreSQL
5. **断点**：更新 `execBlock/execIndex` 到 DB，重启后从这里恢复

## Fan-out 模式

IdentityProcessor 处理完事件后，通过 channel 通知 ReputationProcessor 和 ValidationRegistryProcessor，用于处理依赖关系（如 Reputation 依赖 Identity 的 AgentUID）。

## Commerce Token-Aware Budget

`processor/token.go` 提供：
- **Token 缓存**：`globalTokenInfoCache` 缓存 ERC-20 symbol/decimals
- **价格缓存**：`globalPriceCache` 按 token+时间戳缓存 CoinGecko 价格
- **Job Budget 缓存**：`jobBudgetCache` 缓存 BudgetSet 快照

## 配置

配置文件在 `config/` 目录，按链环境分为 `testnet/` 和 `mainnet/`。

启动时通过 `-f` 指定配置文件：
```bash
go run ./indexer/main.go -f ./config/testnet/base_sepolia.yaml
```

## 编译

```bash
go build ./indexer/main.go
```

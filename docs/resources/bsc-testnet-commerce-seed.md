# BSC Testnet：自造 AgenticCommerce 链上数据

用于在 **chain_id = 97**（BSC Testnet）上对已部署的 **AgenticCommerce** 发交易，产生 `JobCreated` / `JobFunded` / `JobSubmitted` / `JobCompleted` 等事件，供 Indexer（`indexer/processor/commerce.go`）写入 `commerce_actions`。

## 前置条件（一次性核对）

| 项 | 说明 |
|----|------|
| RPC | 可写 HTTPS RPC；勿把私钥写入仓库 |
| Gas | Client / Provider / Evaluator 钱包有足够 **tBNB** |
| 合约 | `AgenticCommerce` 地址与 Indexer 配置里 `commerce.addr` **完全一致**（须为 `0x` 开头的 20 字节地址） |
| Payment token | 与链上 `paymentToken()` 一致；**Client** 持有足够测试代币 |
| Approve | `fund` 前需对 commerce 合约 `approve` 至少 **`-budget`** 数量（脚本在 `happy-path` / `approve-fund` 中处理） |
| 角色 | **Client**：`createJob`、`setBudget`、`fund`；**Provider**：`submit`；**Evaluator**：`complete`（可用 `-evaluator-key`，默认同 `provider-key`） |
| 数据库 | 已执行 `migrations/202604071600_commerce_init_safe.psql`（或团队约定的 commerce 表迁移） |
| Indexer | 对应链配置里 `commerce.run: true`，`start_block` 略小于发交易区块（脚本结束会打印建议值） |

## 配置文件（必需）

脚本**只从配置文件读取**（不再从环境变量读取）。推荐使用本地文件 `config/commerce_seed.local.yaml`（已加入 `.gitignore`，不会被提交）。

先生成 3 个钱包并写入本地配置：

```bash
go run ./cmd/commerce_seed/ -init-config -out ./config/commerce_seed.local.yaml
```

然后编辑该文件里的 `rpc_url` 与 `commerce`（合约地址），并确保 client 钱包有足够 tBNB 与 paymentToken 余额。

## 命令

### 运行前检查余额/授权（推荐先跑）

```bash
go run ./cmd/commerce_seed/ -f ./config/commerce_seed.local.yaml -command check
```

会检查：
- 三个账户的 **tBNB**（native）余额
- `paymentToken()` 对应的 **client token 余额**
- client → commerce 的 **allowance** 是否 ≥ `budget`

如果你在分步调试时带了 `-job-id`，还会额外读取 `getJob(jobId).budget`，提前发现 `BudgetMismatch` 风险：

```bash
go run ./cmd/commerce_seed/ -f ./config/commerce_seed.local.yaml -command check -job-id 123
```

### happy-path（推荐）

```bash
go run ./cmd/commerce_seed/ \
  -f ./config/commerce_seed.local.yaml \
  -command happy-path
```

说明：

- `happy-path` 顺序为：`createJob` → `setBudget`（与 `-budget` 一致）→ `ERC20 approve` → `fund` → `submit` → `complete`。
- `-expired-at` 默认 0 表示「当前时间 + 7 天」。
- 分步调试：`-command create-job` → `set-budget`（需 `-job-id`）→ `approve-fund` → `submit-only` → `complete-only`。

## Indexer 配置对齐

在Indexer 使用的 YAML（结构见 `config/conf.go` 的 `IndexerConfig.Commerce`）中设置：

- `commerce.addr`：与 `-commerce` 相同。
- `commerce.start_block`：使用脚本输出的 **suggested commerce.start_block**，或略小若干块以覆盖全部相关日志。
- `commerce.run: true`。

`rpc_url` 需与发交易时使用的网络一致（BSC Testnet）。

## 验证

**链上**：在 BscScan Testnet 查看交易收据中的事件 topic。

**数据库**：

```sql
SELECT count(*) FROM commerce_actions
WHERE chain_id = '97' AND commerce_contract = '0xYourAgenticCommerce';
```

**API**（示例，按实际路由与参数调整）：`GET agent/identity/commerce/actions?...` — `uid` 可由 Indexer 为各钱包创建 stub agent 后从 `agents` 表反查。

## 常见失败

- `BudgetMismatch`：`fund` 的 `expectedBudget` 必须与链上 `getJob(jobId).budget` 一致；先用 `setBudget` 或 `happy-path`。
- `approve` 不足 / token 错误：确认 `paymentToken()` 与钱包余额、spender 为 commerce 地址。
- 合约 `pause` 或 hook 白名单：非零 `-hook` 时需满足合约规则。

# ERC-8183 Commerce（展示细化 + Indexer 协同）执行清单

> 设计基线：[/Users/pan/work/rido/bas/erc8004/docs/designs/erc8183-commerce-reputation.md](/Users/pan/work/rido/bas/erc8004/docs/designs/erc8183-commerce-reputation.md)
>
> 本文是“可勾选”的落地清单：确保 **Indexer 写入事实表**、**Model 查询能力**、**API 参数/返回** 与 UI 展示第 8 章要求一致。

---

## 0. 范围与验收口径

- [ ] **UI 需要的每一列**（时间/Job/行为/信号/金额/对手方/前序状态/Hook/Tx）在 `commerce_actions` 都有对应字段，且由 Indexer 在合适事件写入。
- [ ] `GET agent/commerce/actions` 支持第 8.4~8.7 的过滤与 8.7 的排序，分页稳定。
- [ ] `GET agent/commerce/scores` 返回 0~3 个 role 的 score（role 维度），无数据时返回空数组或空对象但结构稳定。
- [ ] Indexer 出现临时 RPC 错误（`getJob` 失败、节点抖动）不会永久卡死在同一条 log 上。
- [ ] 写入幂等（重复扫链/重启/回放不会产生重复行），并有唯一约束兜底。

---

## 1. 展示列 ↔ DB 字段 ↔ 链上事件映射（alignment matrix）

> 说明：这是“展示字段是否能被索引出来”的最小真相表。若本表中某行“链上事件/handler”为空，意味着 UI 要展示但 indexer 并未生产该字段，需要补齐。

### 1.1 通用列

| UI 列名 | `commerce_actions` 字段 | 生产来源（事件 / handler） | 备注 |
|---|---|---|---|
| 时间 | `block_timestamp` | 所有事件写入时从 log 取 | 统一为链上 block time |
| Job | `job_id` + (`chain_id`,`commerce_contract`) | 所有事件 | 用于拼 job 详情链接 |
| 行为 | `action` | 事件类型映射 | 标签色由前端决定 |
| 信号 | `signal_polarity`,`signal_weight` | `DetermineSignal(action, previous_status, role)` | 与触发器聚合一致 |
| 金额 | `job_budget` | `JobFunded`/终态(`getJob`) | 是否对中间事件补值见 1.2 |
| 对手方 | `counterparty` | 终态 fan-out | 非终态是否需要由 UI 决策 |
| 前序状态 | `previous_status` | `JobCompleted/Rejected/Expired`（推断） | 仅终态有值 |
| Hook | `hook_address` | `JobCreated` 已有；其它需补齐 | UI 过滤依赖 |
| Tx | `tx_hash` | 所有事件 | 区块浏览器链接 |

### 1.2 Role 维度：事件覆盖（按设计第 8 章）

| Role | UI 期望行为（最小集合） | 链上事件 | 当前 indexer 是否已处理 |
|---|---|---|---|
| provider | `job_submitted`, `job_completed`, `job_rejected`, `job_expired` | `JobSubmitted`, `JobCompleted`, `JobRejected`, `JobExpired` | 已有（终态通过 fan-out） |
| client | `job_created`, `job_funded`, `job_completed`, `job_rejected`, `job_expired` | `JobCreated`, `JobFunded`, `JobCompleted`, `JobRejected`, `JobExpired` | 已有（created/funded 独立写入） |
| evaluator | `job_completed`, `job_rejected`, `job_expired` | `JobCompleted`, `JobRejected`, `JobExpired` | 已有（终态通过 fan-out） |

### 1.3 合约 ABI 额外事件（候选：用于“时间线审计”）

来自 [/Users/pan/work/rido/bas/erc8004/abi/AgenticCommerce.json](/Users/pan/work/rido/bas/erc8004/abi/AgenticCommerce.json)：

| 事件 | UI 是否需要展示到行为列表 | 需要写入哪些列 | 决策/备注 |
|---|---|---|---|
| `ProviderSet(jobId, provider)` | [ ] 是 / [ ] 否 | `action`,`agent_uid`,`agent_address`,`job_id`,`counterparty?`,`hook_address?` | 若展示“谁改了 provider” |
| `BudgetSet(jobId, amount)` | [ ] 是 / [ ] 否 | `action`,`job_budget`,`job_id`,`hook_address?` | 若展示“预算调整” |
| `PaymentReleased(jobId, provider, amount)` | [ ] 是 / [ ] 否 | `action`,`job_budget?`（或 `amount` 新列） | 本期默认不纳入 reputation |
| `Refunded(jobId, client, amount)` | [ ] 是 / [ ] 否 | 同上 | 同上 |
| `EvaluatorFeePaid(jobId, evaluator, amount)` | [ ] 是 / [ ] 否 | 同上 | 同上 |
| `HookWhitelistUpdated(hook, status)` | [ ] 是 / [ ] 否 | 与 job 无关（可能单独页） | 通常不进 agent 时间线 |

---

## 2. Indexer（CommerceProcessor）清单

目标文件：[/Users/pan/work/rido/bas/erc8004/indexer/processor/commerce.go](/Users/pan/work/rido/bas/erc8004/indexer/processor/commerce.go)

- [ ] **事件 topic 订阅**：根据 1.3 的决策，补齐 FilterLogs topics 与 `dealWithEvent` 分发。
- [ ] **字段补全**：为需要 `hook_address` 过滤的行为补齐 hook（必要时 `getJob`）。
- [ ] **幂等写入**：写入走 `ON CONFLICT DO NOTHING`（或等价 upsert），重复回放不产生重复行。
- [ ] **RPC/游标可靠性**：`getJob` 失败不应导致 processor 永久 return；应区分可重试错误并退避重试，成功后再推进游标。
- [ ] **状态推断**：`previous_status` 推断逻辑覆盖 rejected/expired 的常见路径，并有测试。

---

## 3. Model（查询/写入）清单

目标文件：[/Users/pan/work/rido/bas/erc8004/model/commerce.go](/Users/pan/work/rido/bas/erc8004/model/commerce.go)

- [ ] `CreateCommerceAction` 改为 DB 幂等写入（不走 count+insert）。
- [ ] `GetCommerceActionsByAgentUID` 支持：
  - [ ] `action` 多选（或明确仅单值并在 API 层限制）
  - [ ] `hook_address` / 有无 hook 过滤
  - [ ] `min_budget` / `max_budget`
  - [ ] `start_time` / `end_time`
  - [ ] `signal_polarity` 多选
  - [ ] `sort_by`: `timestamp` / `budget` / `signal_weight`

---

## 4. API（参数与返回）清单

目标文件：[/Users/pan/work/rido/bas/erc8004/server/api/router.go](/Users/pan/work/rido/bas/erc8004/server/api/router.go) + `handle/` + `logic/`

- [ ] `GET agent/commerce/actions` 参数对齐设计第 9 章（含 sort/time/budget/hook/counterparty）。
- [ ] `GET agent/commerce/scores` 返回 role 维度数组（或稳定结构），无数据时不报错。

---

## 5. 测试与验证清单

- [ ] `go test ./...`（本地 Postgres 集成测试可运行）
- [ ] 针对典型查询做 `EXPLAIN`（软性）：按 `agent_uid + role + order by block_timestamp desc` 与按 `job_id` 查询命中索引。


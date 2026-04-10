# ERC-8183 Commerce UI API Spec（前端对齐版）

本文面向前端开发，定义 Commerce 展示所需 API 的**参数、默认值、返回结构、枚举与示例**，并保证与当前后端实现能力一致。

设计基线：[/Users/pan/work/rido/bas/erc8004/docs/designs/erc8183-commerce-reputation.md](/Users/pan/work/rido/bas/erc8004/docs/designs/erc8183-commerce-reputation.md) 第 8 章（展示）与第 9 章（API 概览）。

权威性与冗余策略（已对齐）：
- **权威来源**：`GET agent/commerce/scores` 是 commerce score 的权威来源（尤其是 **分段聚合**）。
- **便捷字段**：`GET agent/identity/detail` 中的 `agent.commerce_score` 是**可选便捷字段**（用于首屏减少一次请求）。当 DB 缺表或无数据时允许为空，不影响 detail 主流程。

后端实现参考：
- 路由：[/Users/pan/work/rido/bas/erc8004/server/api/router.go](/Users/pan/work/rido/bas/erc8004/server/api/router.go)
- Handler：[/Users/pan/work/rido/bas/erc8004/server/api/handle/commerceHandle.go](/Users/pan/work/rido/bas/erc8004/server/api/handle/commerceHandle.go)
- Logic：[/Users/pan/work/rido/bas/erc8004/server/api/logic/commerceLogic.go](/Users/pan/work/rido/bas/erc8004/server/api/logic/commerceLogic.go)
- Model Query：[/Users/pan/work/rido/bas/erc8004/model/commerce.go](/Users/pan/work/rido/bas/erc8004/model/commerce.go)
- Indexer 写入（事实来源）：[/Users/pan/work/rido/bas/erc8004/indexer/processor/commerce.go](/Users/pan/work/rido/bas/erc8004/indexer/processor/commerce.go)

---

## 0. 约定与通用规则

### 0.1 角色（role）

- `provider`
- `client`
- `evaluator`

### 0.2 行为类型（action）

来自 `commerce_actions.action`：

- 终态（definitive）：`job_completed` / `job_rejected` / `job_expired`
- 过程（indicative）：`job_created` / `job_funded` / `job_submitted`
- 审计类过程事件（indicative）：`provider_set` / `budget_set`

> 前端展示「行为」列时用该枚举直出（label/颜色由前端映射）。

### 0.3 信号（signal）

- `signal_certainty`: `definitive` / `indicative`
- `signal_polarity`: `positive` / `negative` / `neutral`
- `signal_weight`: number（小数；可用于“按信号强度排序”）

### 0.3.1 枚举与常量表（前端可直接 copy）

| 名称 | 取值 | 传参形式 | 示例 |
|---|---|---|---|
| `role` | `provider` / `client` / `evaluator` | 单值 | `role=provider` |
| `action` | `job_created` / `job_funded` / `job_submitted` / `job_completed` / `job_rejected` / `job_expired` / `provider_set` / `budget_set` | CSV 多选 | `action=job_completed,job_rejected` |
| `certainty` | `definitive` / `indicative` | 单值 | `certainty=definitive` |
| `polarity` | `positive` / `negative` / `neutral` | CSV 多选 | `polarity=positive,neutral` |
| `sort_by` | `timestamp` / `budget` / `signal_weight` | 单值 | `sort_by=budget` |

### 0.4 时间

`block_timestamp` 是 **unix 秒**（uint64）。前端可转为本地时间展示。

### 0.5 分页

- `page`：默认 `1`，小于等于 0 会被后端纠正为 `1`
- `page_size`：默认 `20`，小于等于 0 → `20`，大于 `100` → `100`

返回体包含 `total`（总行数），前端自行计算总页数。

### 0.6 多选参数格式（已对齐确认）

- `action` 多选：**CSV**，例如 `action=job_completed,job_rejected`
- `polarity` 多选：**CSV**，例如 `polarity=positive,neutral`

### 0.7 前端使用方式（推荐加载策略）

推荐前端按“先 detail、后模块化补充”的方式加载，减少首屏额外请求与 UI 抖动：

- **进入 Agent 详情页**
  - 请求 `GET agent/identity/detail?uid={uid}`
  - 使用返回的可选便捷字段 `data.agent.commerce_score` 判断：
    - 是否展示 Provider/Client/Evaluator 三个 tab（例如 `total_jobs > 0` 才展示）
    - Commerce 区域的空状态（全部 role 均无数据 → “暂无商业行为记录”）

- **进入某个 Role Tab**
  - 请求 `GET agent/commerce/actions?...` 获取行为列表（分页/过滤/排序）

- **需要分段聚合（按 chain+contract）或刷新权威 score**
  - 请求 `GET agent/commerce/scores?...`（本接口是权威来源，尤其分段聚合）

### 0.8 响应封装协议（项目通用约定）

本项目接口统一使用 JSON 包装层（成功/失败都会返回 HTTP 200）：

- **成功响应**
  - `code = 0`
  - `status = "success"`
  - 业务数据在 `data` 字段内
- **失败响应**
  - `code = 1`
  - `status = "failed"`
  - 错误信息在 `msg` 字段内（例如 `"Invalid Request"` / `"Internal Error"`）

---

## 1. GET `agent/commerce/scores`

用于渲染每个 role tab 的 **Score 卡片 + 统计摘要**。返回数组，包含 0~3 个 role 的 score；某 role 没数据则不返回该 role 记录。

> 注意：如果前端使用了 `agent detail` 的 `agent.commerce_score` 作为首屏便捷字段，仍应以本接口作为后续刷新/分段聚合的权威来源。

### 1.1 Query 参数

- `uid`（必填，uint64）：agent uid
- `chain_id`（可选，string）：指定链（与 `commerce_contract` 同时提供时生效）
- `commerce_contract`（可选，string）：指定合约（与 `chain_id` 同时提供时生效）

**行为：**
- 若同时提供 `chain_id` + `commerce_contract`：返回该段（segmented）聚合数据（读 `commerce_scores`）
- 否则：返回全局 rollup 聚合数据（读 `commerce_scores_global`）

### 1.1.1 segmented 模式口径（精确定义）

- 仅当 **同时提供** `chain_id` 与 `commerce_contract` 时，才进入 segmented 模式。
- segmented 的含义是“该链 + 该合约地址范围内”的聚合（按 `commerce_scores` 分段主键）。
- 前端需要 segmented 时，不应依赖 `agent detail` 的 `commerce_score`（该字段为全局便捷摘要）。

### 1.2 返回结构

```json
{
  "scores": [
    {
      "role": "provider",
      "completed_count": 0,
      "rejected_count": 0,
      "expired_responsible_count": 0,
      "success_rate": 0,
      "weighted_score": 0,
      "created_count": 0,
      "funded_count": 0,
      "funded_rate": 0,
      "completion_rate": 0,
      "evaluated_count": 0,
      "expired_from_submitted_count": 0,
      "responsiveness": 0,
      "total_jobs": 0,
      "total_volume": 0,
      "unique_counterparties": 0,
      "confidence": 0
    }
  ]
}
```

字段含义与展示建议：见设计文档第 8.4~8.6 章。

### 1.3 示例

- 全局：`GET agent/commerce/scores?uid=123`
- 分段：`GET agent/commerce/scores?uid=123&chain_id=84532&commerce_contract=0xabc...`

### 1.4 错误与校验规则

- **参数校验**
  - `uid` 必须是 uint64
  - `chain_id` 与 `commerce_contract` 必须同时提供才会进入“分段聚合”模式；否则返回全局 rollup
- **错误返回**
  - `uid` 非法：`code=1`，`msg="Invalid Request"`
  - DB/内部异常：`code=1`，`msg="Internal Error"`

---

## 2. GET `agent/commerce/actions`

用于渲染 role tab 的 **行为列表（表格 + 分页 + 过滤 + 排序）**。

### 2.1 Query 参数（全部为可选，除 uid 外）

- **必填**
  - `uid`（uint64）

- **过滤（与设计第 8 章一致）**
  - `role`（string）：`provider|client|evaluator`
  - `action`（CSV，多选）：见 0.2
  - `certainty`（string）：`definitive|indicative`
  - `polarity`（CSV，多选）：`positive|negative|neutral`
  - `chain_id`（string）：链 id
  - `commerce_contract`（string）：合约地址
  - `counterparty`（string）：对手方地址（精确匹配）
  - `hook_address`（string）：hook 地址（精确匹配；空 hook 由 `has_hook` 控制）
  - `has_hook`（bool）：`true`=仅有 hook（`hook_address <> ''`），`false`=仅无 hook（`hook_address = ''`）
  - `min_budget`（float64）
  - `max_budget`（float64）
  - `start_time`（uint64）：最小 `block_timestamp`
  - `end_time`（uint64）：最大 `block_timestamp`

- **排序**
  - `sort_by`（string，可选）：`timestamp`（默认） / `budget` / `signal_weight`

- **分页**
  - `page`（int，默认 1）
  - `page_size`（int，默认 20，最大 100）

### 2.2 返回结构

```json
{
  "actions": [
    {
      "uid": 0,
      "chain_id": "84532",
      "commerce_contract": "0x...",
      "job_id": 0,
      "agent_uid": 123,
      "agent_address": "0x...",
      "role": "provider",
      "action": "job_completed",
      "signal_polarity": "positive",
      "signal_weight": 1,
      "signal_certainty": "definitive",
      "job_budget": 0,
      "counterparty": "0x...",
      "reason": "",
      "deliverable": "",
      "previous_status": "submitted",
      "hook_address": "0x...",
      "block_number": 0,
      "tx_hash": "0x...",
      "log_index": 0,
      "block_timestamp": 0
    }
  ],
  "total": 0
}
```

### 2.3 与 UI 列对齐（最小映射）

- 时间：`block_timestamp`
- Job：`job_id` + `chain_id` + `commerce_contract`
- 行为：`action`
- 信号：`signal_polarity` + `signal_weight`
- 金额：`job_budget`
- 对手方：`counterparty`（主要在终态扇出记录里有值）
- 前序状态：`previous_status`（仅终态）
- Hook：`hook_address`（非空字符串即认为有 hook）
- Tx：`tx_hash`

### 2.3.1 字段语义（前端展示注意）

| 字段 | 语义 | 可能为空/默认值 | 前端建议 |
|---|---|---|---|
| `hook_address` | hook 地址 | **空字符串**表示无 hook | 用 `has_hook` 做开关过滤；展示时空则不显示标签 |
| `previous_status` | 终态前序状态（from open/funded/submitted） | 仅终态行有意义，其它可能为空 | 非终态行展示为 `-` 或隐藏该列值 |
| `counterparty` | 对手方地址 | 主要终态 fan-out 会填充；过程事件可能为空 | 为空则展示 `-` |
| `job_budget` | 预算金额 | 可能为 0（未设置或未写入） | 0 值可展示为 0；若要区分未知可在 UI 侧加占位规则 |

### 2.3.2 排序与分页稳定性

- `sort_by=timestamp`：按 `block_timestamp DESC`
- `sort_by=budget`：按 `job_budget DESC`，并以 `block_timestamp DESC` 作为二级排序（保证翻页稳定）
- `sort_by=signal_weight`：按 `signal_weight DESC`，并以 `block_timestamp DESC` 作为二级排序（保证翻页稳定）

### 2.4 示例

- Provider tab 默认（最新在前）：
  - `GET agent/commerce/actions?uid=123&role=provider&page=1&page_size=20`

- Evaluator tab（只看终态）：
  - `GET agent/commerce/actions?uid=123&role=evaluator&certainty=definitive`

- Client tab（created+funded+终态，且只看有 hook）：
  - `GET agent/commerce/actions?uid=123&role=client&action=job_created,job_funded,job_completed,job_rejected,job_expired&has_hook=true`

- 按金额排序 + 金额区间：
  - `GET agent/commerce/actions?uid=123&role=provider&sort_by=budget&min_budget=0.1&max_budget=10`

---

## 2.5 前端默认查询（Role Tab 推荐口径）

以下是前端“首次进入某个 role tab”时推荐使用的默认 query（与设计第 8 章一致，且后端可实现）。

### 2.5.1 Provider Tab（默认）

- **默认列表（含 submitted + 终态）**

`GET agent/commerce/actions?uid={uid}&role=provider&action=job_submitted,job_completed,job_rejected,job_expired&page=1&page_size=20&sort_by=timestamp`

> 说明：provider 的 `job_submitted` 为 indicative，有助于时间线完整；其余为 definitive。

### 2.5.2 Client Tab（默认）

- **默认列表（含 created/funded + 终态）**

`GET agent/commerce/actions?uid={uid}&role=client&action=job_created,job_funded,job_completed,job_rejected,job_expired&page=1&page_size=20&sort_by=timestamp`

> 说明：设计第 8.5 明确 client tab 默认包含 indicative（created/funded）。

### 2.5.3 Evaluator Tab（默认）

- **默认列表（仅 definitive）**

`GET agent/commerce/actions?uid={uid}&role=evaluator&certainty=definitive&action=job_completed,job_rejected,job_expired&page=1&page_size=20&sort_by=timestamp`

---

## 2.6 前端过滤器 → Query 参数映射（可直接落 UI 交互）

### 2.6.1 行为类型（多选）

- UI 多选行为 → `action` CSV
- 示例：选 completed + rejected
  - `action=job_completed,job_rejected`

### 2.6.2 信号极性（多选）

- UI 多选极性 → `polarity` CSV
- 示例：选 positive + neutral
  - `polarity=positive,neutral`

### 2.6.3 确定性（单选）

- 全部：不传 `certainty`
- 仅终态：`certainty=definitive`
- 仅过程：`certainty=indicative`

> 注意：如果 UI 同时传了 `certainty=definitive` 且 `action` 里包含 `job_created/job_funded/job_submitted`，后端会以 WHERE 逻辑自然过滤掉这些 indicative 行（不报错）。

### 2.6.4 Hook 过滤（开关）

- 仅有 hook：`has_hook=true`
- 仅无 hook：`has_hook=false`
- 精确指定某 hook：`hook_address=0x...`（此时不必再传 `has_hook`）

### 2.6.5 时间范围（日期区间）

- UI 日期区间需转为 unix 秒
- **包含整天口径（已对齐）**：
  - `start_time`：起始日期的 `00:00:00`
  - `end_time`：结束日期的 `23:59:59`
- 示例（按本地时区选择日期范围）：
  - `start_time=1735689600`（2025-01-01 00:00:00）
  - `end_time=1735862399`（2025-01-02 23:59:59）

### 2.6.6 金额范围（数值区间）

- `min_budget` / `max_budget`（float）

### 2.6.7 排序切换

- 最新：`sort_by=timestamp`（或不传）
- 按金额：`sort_by=budget`
- 按信号强度：`sort_by=signal_weight`

### 2.6.8 前端 state <-> query（可逆性建议）

- **多选序列化**：`action` / `polarity` 使用 CSV；提交前做 `trim` + 去空值 + 去重（保持稳定顺序方便缓存与对比）。
- **清空筛选**：不要传空字符串（例如不要传 `action=`），而是直接不带该参数。
- **时间区间**：按 2.6.5 的“包含整天”口径在前端完成转换，再传 unix 秒给后端。

### 2.7 错误与校验规则

- **必填校验**
  - `uid` 必须是 uint64

- **枚举校验（传了但非法 → Invalid Request）**
  - `role`：`provider|client|evaluator`
  - `certainty`：`definitive|indicative`
  - `sort_by`：`timestamp|budget|signal_weight`
  - `action`：CSV 中每个值必须在允许枚举内（见 0.2）
  - `polarity`：CSV 中每个值必须在允许枚举内（`positive|negative|neutral`）

- **类型/范围校验**
  - `has_hook`：bool
  - `min_budget` / `max_budget`：number（不允许 NaN/Inf）
  - `start_time` / `end_time`：uint64（unix 秒）

- **错误返回（注意：HTTP 200 + code/msg）**
  - 参数非法：`code=1`，`msg="Invalid Request"`
  - DB/内部异常：`code=1`，`msg="Internal Error"`

---

## 3. GET `agent/identity/detail`（`commerce_score` 便捷字段）

本接口本身是 agent 详情页的主接口；Commerce 模块在其中提供一个**可选便捷字段**，用于首屏减少一次 `agent/commerce/scores` 请求。

- **字段路径**：`data.agent.commerce_score`
- **字段形状**：按 role 组织的对象（map），key 为 `provider|client|evaluator`
- **字段值**：与 `GET agent/commerce/scores` 的单条 score **字段完全一致**（见 1.2），仅组织方式不同
- **权威性声明**：本字段仅为便捷摘要；分段聚合与后续刷新以 `GET agent/commerce/scores` 为准
- **空/缺表行为**：
  - 若 DB 未初始化 commerce 表或该 agent 无 commerce 数据：`commerce_score` 可能缺失或为空对象
  - 前端应按“无数据”处理，不应阻塞详情页渲染

### 3.1 示例响应片段（最小）

```json
{
  "code": 0,
  "status": "success",
  "msg": "",
  "data": {
    "agent": {
      "uid": 123,
      "commerce_score": {
        "provider": {
          "role": "provider",
          "completed_count": 0,
          "rejected_count": 0,
          "expired_responsible_count": 0,
          "success_rate": 0,
          "weighted_score": 0,
          "created_count": 0,
          "funded_count": 0,
          "funded_rate": 0,
          "completion_rate": 0,
          "evaluated_count": 0,
          "expired_from_submitted_count": 0,
          "responsiveness": 0,
          "total_jobs": 0,
          "total_volume": 0,
          "unique_counterparties": 0,
          "confidence": 0
        }
      }
    }
  }
}
```

### 3.2 前端渲染规则（建议）

- 若 `commerce_score` 缺失或为空对象：视为 Commerce 无数据，展示“暂无商业行为记录”。
- 若某个 role 的 `commerce_score[role].total_jobs <= 0`：该 role tab 不展示（与设计第 8 章一致）。
- tab 详情列表统一通过 `GET agent/commerce/actions` 获取，不应使用 `commerce_score` 替代列表查询。
- 分段聚合与刷新以 `GET agent/commerce/scores` 为准（`agent detail` 仅为便捷摘要）。

---

## 4. 错误返回（约定）

当前 handler 对以下参数会返回 `Invalid Request`：
- `uid` 非法
- `has_hook` 非法（非 bool）
- `min_budget` / `max_budget` 非法（NaN/Inf/非数字）
- `start_time` / `end_time` 非法（非 uint64）

其余查询错误统一返回 `Internal Error`。

---

## 5. 后端可实现性声明（前端联调须知）

### 5.1 数据来源与一致性

- `actions` 的所有行均来源于 indexer 写入的 `commerce_actions`；不依赖实时 RPC 查询。
- `scores` 来源于 DB 触发器维护的 `commerce_scores` / `commerce_scores_global`（需先完成建表与 trigger 初始化）。

### 5.2 必要的 DB 初始化

若数据库尚无 commerce 表，请先执行安全建表脚本：

- [/Users/pan/work/rido/bas/erc8004/migrations/202604071600_commerce_init_safe.psql](/Users/pan/work/rido/bas/erc8004/migrations/202604071600_commerce_init_safe.psql)

否则即使 indexer 正常启动，写入会失败，前端查询也会为空/报错。

---

## 6. 联调 Checklist（前端验收用）

- [ ] DB 已初始化 commerce 表与 trigger（见 `migrations/202604071600_commerce_init_safe.psql`）
- [ ] indexer 已运行且 DB 内有 `commerce_actions` 数据（可用 actions 接口验证 `total > 0`）
- [ ] detail 首屏：`GET agent/identity/detail?uid={uid}` 返回 `data.agent`，且 `commerce_score` 缺失/为空时页面不报错
- [ ] Provider Tab 默认 query 可用（见 2.5.1）
- [ ] Client Tab 默认 query 可用（见 2.5.2）
- [ ] Evaluator Tab 默认 query 可用（见 2.5.3）
- [ ] 过滤器组合回归：
  - [ ] `has_hook=true/false`
  - [ ] `action` 多选（CSV）
  - [ ] `polarity` 多选（CSV）
  - [ ] `min_budget/max_budget`
  - [ ] 时间区间（包含整天：end_time=23:59:59）
  - [ ] `sort_by` 切换（timestamp/budget/signal_weight）


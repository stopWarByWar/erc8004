# ERC-8183：`GET /agent/commerce/jobs` 列表项返回字段增补说明

本文**仅**描述 Job List 接口在单条 `job` 对象上**新增或变更的返回字段**，不重复整页 UI / Filters / 其它接口规范。

**接口**：`GET /agent/commerce/jobs` → `data.jobs[]` 中每一项。

---

## 1. 字段一览

| JSON 字段 | 类型（JSON） | 是否必有 | 说明 |
| --- | --- | --- | --- |
| `job_uid` | number | 是 | 库内主键（`commerce_jobs` 主键）。建议作为表格行 `key`，与 `chain_id`+`job_id` 不同，全局唯一。 |
| `chain_name` | string | 否 | 由后端按 `chain_id` 解析链配置；命中则有值，与 `GET .../jobs/detail` 中 `job` 同源。 |
| `chain_logo` | string（URL） | 否 | 同上；用于列表展示链图标。缺省或加载失败时前端回退为仅展示 `chain_id`（见母 spec §2.2）。 |
| `payment_decimals` | number | 是 | ERC20 decimals；库内为 0 时后端**归一为 18**。与 `token_symbol`、字符串金额字段一起用于展示/校验。 |
| `budget` | string | 是 | 原币预算金额，**字符串**（`formatAmount`：最多保留 8 位小数再 trim，无科学计数法）。 |
| `paid_amount` | string | 是 | 同上（非 `completed` 时 UI 仍可按母 spec 对 Paid 列显示 `—`，与字段是否为 `"0"` 无关）。 |
| `platform_fee_amount` | string | 是 | 同上。 |
| `evaluator_fee_amount` | string | 是 | 同上。 |

**未改名的 USD 字段**（`budget_usd`、`paid_amount_usd` 等）仍为 number，语义不变；仅上表中原币金额从「历史上可能被误当作 number」统一为 **string**，与 detail 对齐。

---

## 2. 响应示例片段（仅突出增补字段）

```json
{
  "job_uid": 900001,
  "chain_id": "8453",
  "chain_name": "Base",
  "chain_logo": "https://...",
  "commerce_contract": "0x...",
  "job_id": 12,
  "payment_token": "0x...",
  "payment_decimals": 18,
  "token_symbol": "ETH",
  "budget": "0.5",
  "paid_amount": "0.45",
  "platform_fee_amount": "0.03",
  "evaluator_fee_amount": "0.02"
}
```

---

## 3. 前端类型提示（节选）

```ts
/** 与 GET /agent/commerce/jobs 单条对齐；仅标出本次增补/变更 */
type CommerceJobListItemDelta = {
  job_uid: number;
  chain_name?: string;
  chain_logo?: string;
  payment_decimals: number;
  budget: string;
  paid_amount: string;
  platform_fee_amount: string;
  evaluator_fee_amount: string;
};
```

完整 `CommerceJob` 仍以母 spec §8 为准，实现时请将上列字段与既有类型**合并**（勿再假定上述金额为 `number`）。

# Commerce Jobs Filters API：`payment_tokens` 返回结构变更说明

**日期**：2026-05-12  
**接口**：`GET /agent/commerce/jobs/filters`（经网关后常见路径为 `/agent/api/commerce/jobs/filters`，以实际部署为准）  
**影响范围**：响应体中 `data.filters.payment_tokens` 字段的 **元素类型**（破坏性变更）。

主规格仍见：[202604131905_erc8183-job-api.md](./202604131905_erc8183-job-api.md)。

---

## 1. 未变更部分

以下字段与此前一致，**无结构变化**：

| 路径 | 说明 |
|------|------|
| `filters.chains` | `chain_id`、`chain_name`、`chain_logo` |
| `filters.commerce_contracts` | `string[]` |
| `filters.last_updated` | `number`（Unix 秒） |

列表接口 `GET /agent/commerce/jobs` 的查询参数 **未改名**：仍使用 `payment_token`、`chain_id` 等，与本次变更无关。

---

## 2. 变更内容：`filters.payment_tokens`

### 2.1 以前（已废弃）

`payment_tokens` 为 **字符串数组**，每个元素为代币合约地址（与列表查询参数 `payment_token` 取值相同），**不包含链维度**。

```json
{
  "payment_tokens": [
    "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
    "0x0000000000000000000000000000000000000000"
  ]
}
```

### 2.2 现在（当前行为）

`payment_tokens` 为 **对象数组**。每条记录对应快照表 `commerce_jobs` 中一个 **去重键** `(chain_id, payment_token)`，用于在 UI 上展示「某链上的某支付代币」，并在筛选时与 `chain_id` + `payment_token` 查询参数对齐。

| JSON 字段 | 类型 | 说明 |
|-----------|------|------|
| `chain_id` | string | 链 ID，来自 `commerce_jobs.chain_id` |
| `chain_logo` | string | 链图标 URL；来自运行时配置 `GetChainInfo(chain_id).ChainLogo`，与 `filters.chains[].chain_logo` 同源 |
| `symbol` | string | 代币符号；来自同组行上聚合后的 `token_symbol`（若库中均为空则可能省略或为空） |
| `contract` | string | 代币合约地址，等于原字符串元素，对应列表查询参数 **`payment_token`** |
| `logo` | string | 代币图标；当前无库表/索引字段填充，一般为空。使用 `json:"logo,omitempty"`，**空时通常不出现在 JSON 中** |

示例：

```json
{
  "payment_tokens": [
    {
      "chain_id": "8453",
      "chain_logo": "https://example.com/base.png",
      "symbol": "USDC",
      "contract": "0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"
    },
    {
      "chain_id": "1",
      "chain_logo": "https://example.com/eth.png",
      "symbol": "USDC",
      "contract": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
    }
  ]
}
```

---

## 3. 前端 / 调用方迁移建议

1. **类型**：将 `payment_tokens: string[]` 改为 `Array<{ chain_id: string; chain_logo?: string; symbol?: string; logo?: string; contract: string }>`（或等价命名）。
2. **下拉展示**：使用 `symbol`、`logo`、`contract`、`chain_logo`、`chain_id` 组合渲染；无 `symbol` 时可回退为缩短的 `contract`。
3. **选中后请求 Job List**：  
   - 设置 `payment_token=<所选项的 contract>`（与旧行为一致）；  
   - **建议同时**设置 `chain_id=<所选项的 chain_id>`，避免跨链同形地址（若存在）时的歧义。  
   列表接口已同时支持 `chain_id` 与 `payment_token` 过滤。

---

## 4. 实现参考（仓库内）

| 说明 | 路径 |
|------|------|
| 响应 DTO | `server/api/types/types.go`（`CommerceJobsFilterPaymentToken`、`CommerceJobsFilters`） |
| 聚合查询 | `model/commerce.go`（`GetCommerceJobsDistinctPaymentTokenFacets`） |
| 缓存组装 | `server/api/utils/jobs_filters_cache.go`（`buildJobsFilters`） |
| 逻辑层（与缓存逻辑对齐） | `server/api/logic/commerceLogic.go`（`BuildCommerceJobsFilters`） |

---

## 5. 版本与兼容

- 本次为 **破坏性变更**：仅依赖 `payment_tokens` 为 `string[]` 的客户端在升级后端后必须同步改解析逻辑。  
- 若需长期双形态兼容，应在网关或 BFF 做适配层（当前后端实现未提供双写/开关）。

# ERC-8183 Commerce Reputation 数据展示提案

> 创建日期：2026-04-10
> 状态：Draft

---

## 一、背景

ERC-8183 Commerce Reputation 数据包含三个层次：
- **Job 层**（`commerce_jobs`）：完整状态机、Budget/Token/资金状态
- **Action 层**（`commerce_actions`）：链上行为明细、全链路可溯源
- **Score 层**（`commerce_scores_global`）：三角色聚合指标

Volume 是核心维度——`total_volume_usd`、`weighted_score_usd`、`token_symbol` 分布都是关键信号。

---

## 二、功能分类

### 第一类：Agent 个人主页（单 Agent 视角）

| 功能 | 说明 | 核心数据 |
|------|------|---------|
| **Provider Tab** | 服务交付能力 | success_rate、weighted_score、volume、completed/rejected/expired 分布 |
| **Client Tab** | 委托方信用 | funded_rate、completion_rate、volume、created/funded 行为 |
| **Evaluator Tab** | 评审方信用 | responsiveness、evaluated_count、expired_from_submitted |
| **综合 Score 卡片** | 三角色汇总 | commerce_scores_global |

### 第二类：Job 浏览器（列表视角）

| 功能 | 说明 | 核心数据 |
|------|------|---------|
| **Job 列表页** | 全局 Job 搜索 + 过滤 | commerce_jobs 分页 |
| **Escrow 看板** | Open / Funded / Submitted 状态分布 | status 统计 |
| **按状态过滤** | Open / Funded / Submitted / Completed / Rejected / Expired | status |
| **按 Token 过滤** | 按 payment_token / token_symbol 筛选 | token_symbol |
| **按金额过滤** | min_budget / max_budget | budget / budget_usd |
| **按时间过滤** | start_time / end_time | updated_at |

### 第三类：Job 详情页（单 Job 视角）

| 功能 | 说明 | 核心数据 |
|------|------|---------|
| **状态机时间线** | 完整生命周期可视化 | action 链 + block_timestamp |
| **资金状态** | budget / paid_amount / paid_amount_usd | budget + payment_token |
| **参与者卡片** | Client / Provider / Evaluator 地址 | client / provider / evaluator |
| **行为明细** | 该 Job 所有事件的完整记录 | commerce_actions by job_id |
| **信号强度** | 终态对应的 polarity/weight | signal_polarity / signal_weight |

### 第四类：统计图表（聚合分析）

| 功能 | 说明 | 核心数据 |
|------|------|---------|
| **Action 分布环形图** | 各 action 类型占比 | stats.action_breakdown |
| **成功率折线图** | 24h / 7d / 30d 趋势 | stats.time_series |
| **Budget 分布直方图** | small/medium/large 按 token 分组 | stats.budget_distribution |
| **终态分布环形图** | completed/rejected/expired 三段比例 | scores.completed/rejected/expired_count |

### 第五类：Volume 分析（核心维度）

| 功能 | 说明 | 核心数据 |
|------|------|---------|
| **USD 总量仪表板** | weighted_score_usd + total_volume_usd | volume 加权分 |
| **Token 分布条形图** | 各 token（ETH/USDT/BNB）占比 | GetTokenVolumeBreakdown |
| **Volume 时序图** | 累计 USD 量随时间变化 | budget_usd 时间聚合 |
| **大额 Job 高亮** | 高 USD 量 Job 单独标注 | budget_usd 超阈值 |
| **跨 Token 换算展示** | 前端做汇率换算 + symbol 显示 | token_symbol + budget |

### 第六类：关系网络（Agent 间交互）

| 功能 | 说明 | 核心数据 |
|------|------|---------|
| **对手方列表** | 该 Agent 合作过的所有对手方 | counterparty 聚合 |
| **合作频次** | 各对手方合作次数 | counterparty count |
| **合作 Volume** | 与某对手方累计 USD 量 | counterparty + budget_usd |

### 第七类：跨维度对比

| 功能 | 说明 | 核心数据 |
|------|------|---------|
| **ERC-8004 vs ERC-8183** | 主观评价 vs 客观行为并排 | feedbacks + commerce_actions |
| **多 Agent 对比** | 最多选 3 个 Agent 的 score 对比 | commerce_scores_global |

---

## 三、优先级建议

**P0（核心，必须）**
- 第一类：Agent 个人主页三 Tab（Provider / Client / Evaluator）
- 第二类：Job 列表页 + 基础过滤
- 第五类：USD 总量仪表板 + Token 分布

**P1（重要）**
- 第三类：Job 详情页（状态机时间线）
- 第四类：统计图表（action 分布 + 成功率折线）

**P2（增强）**
- 第六类：对手方关系网络
- 第七类：跨维度对比

---

## 四、待确认事项

1. Job 详情页的"状态机时间线"UI 形态——是横向步骤条还是垂直时间轴？
2. Token 分布的 USD 换算汇率来源——前端自己查还是后端提供？
3. "大额 Job"阈值——是固定值（如 $1000 USD）还是百分位数？
4. 多 Agent 对比——是否需要 API 支持批量查询？

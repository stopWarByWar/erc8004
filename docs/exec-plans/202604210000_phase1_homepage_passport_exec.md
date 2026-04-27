# Phase 1 Exec Plan: Homepage + Passport（ERC-8183 证据链）

来源（上游真相）：`docs/designs/202604201900_homepage_erc8183_passport_prd.md`
分支上下文（实现位置）：前端 `erc8004-1.1-web`（Next.js）
目标：把“可验证证据链”做成可点击、可分享、可匿名打开的 MVP 闭环。

---

## 0. 固定项目上下文（不讨论方向）

### 0.1 核心用户与任务
- **用户**：Agent 开发者 / 项目方老板 / 增长（卖方）
- **任务**：把一个 agent 的“终态结算证据”压缩成可分享链接，让对方无需登录也能验证

### 0.2 Phase 1 的输出物（用户可见）
- Homepage：Hero + 3 个 proof block（证据入口），并隐藏 commerce trending
- Agent detail：新增 `?view=passport` 的 Passport 视图
- Jobs：作为证据入口，不改逻辑，但需要确保首页/Passport 可稳定跳转并下钻链上 tx

---

## 1. Phase 1 范围（钉死）

### 1.1 本阶段做什么（IN SCOPE）

#### A) Homepage（证据导流）
- Hero 文案按 PRD 4.1.1
- CTA1：`/jobs`
- CTA2：跳转到 Passport（但 **不绑定真实 demo agent**，见 “决策已锁定”）
- Proof blocks（3 张卡）：
  - Outcome-based reputation → `/jobs`（终态过滤预设，若无过滤能力则仅跳 `/jobs`）
  - Sybil-resistant by design → `/jobs`
  - Portable proof → **说明文档链接**（不跳真实 agent）
- Trending：保留现有 Feedback trending，**隐藏 Trending by Commerce**

#### B) Passport view（`/agents/{uid}?view=passport`）
- 展示 SR/WS/CF（来自 `CommerceScore.success_rate / weighted_score / confidence`）
- **按 role 切换展示**（provider / client / evaluator）
- 最近 3 条终态记录摘要（可用 `commerce/actions` 拉取并筛终态）
- View jobs / View onchain tx（证据下钻）
- 弱一致：30s stale（由 react-query 及现有 QueryClient 设置实现即可），失败时：
  - 有缓存显示缓存
  - 顶部红色 banner + Retry

#### C) 最小测试（必须）
- 增加可运行的前端单测（Vitest + RTL）
- 覆盖至少：
  - role tabs 切换（`CommerceReputation` 已有测试基建与样例）
  - Passport 的空态 / 不足态 / 失败态（如相关组件实现后补）

### 1.2 本阶段不做什么（OUT OF SCOPE）
- 独立 `/passport/{uid}` 路由（本期只做 `?view=passport`）
- Commerce Trending / curated samples / proof-feed 首页大改
- 平台排序接入与商业交付
- 动态 OG：**本期降级为站点级默认 OG**（不按 agent 动态生成）
- 隐私/争议仲裁/系统化反作弊（只做最小护栏：不使用 HTML 注入）

---

## 2. 决策已锁定（来自本次评审与用户选择）

- **Passport role**：允许在 Passport 内切换 role（provider/client/evaluator）
- **第三张 Proof 卡**：先链到说明文档，不绑定真实示例 agent
- **OG**：本期仅站点级默认 OG，不做按 agent 动态生成

---

## 3. 交付清单（Deliverables）

### 3.1 代码改动（预期模块）
- `erc8004-1.1-web/src/app/(app)/(home)/`：Hero/Proof blocks/Trending（隐藏 commerce 叙事）
- `erc8004-1.1-web/src/app/(app)/agents/[id]/`：增加 Passport 视图模式与 UI 组件
- 复用 `src/http/agent.ts` 的：
  - `useGetCommerceScores`
  - `useGetCommerceActions`
  - `useGetCommerceStats`（若 Passport 用图表则用，否则可不拉）

### 3.2 文档
- PRD 已内嵌 `GSTACK REVIEW REPORT`（无需再写 design）
- 本文件作为 Phase 1 执行文档（执行层照此做）

---

## 4. 验收（Given / When / Then，对齐 PRD）

1) Passport 无数据
- Given：terminal=0
- When：打开 `/agents/{uid}?view=passport`
- Then：显示“暂无终态 job…”，允许显示 CF=0/5，且给出引导入口（Jobs）

2) Passport 数据不足
- Given：terminal=3
- Then：显示“数据不足 3/5”，不得出现榜单相关信息

3) 接口失败可见
- When：scores 接口失败或 code!=0
- Then：显示缓存（如有）+ 红色 banner + Retry

4) 匿名访问
- When：无登录/无钱包连接打开 Passport 链接
- Then：可渲染，证据链接可点开

5) 首页证据链可走通
- When：从 homepage 点 CTA/三张 proof cards
- Then：到 Jobs 或 Passport，并可下钻到链上 tx

---

## 5. 冲突清单（如出现，先停下来对齐）

当前已知“上游 PRD”与“本期决定”存在一处显式调整：
- PRD 原 4.2.4 提到动态 OG（按 agent 渲染），但已决策 **本期降级为站点级默认 OG**。
  - 处理：PRD 已同步写明降级原因，本 exec plan 以此为准。

---

## 6. 风险与护栏（Phase 1 最小）

- 风险：分享页易被爬，接口压力上升
  - 护栏：前端不做轮询；Retry 需用户主动点击
- 风险：数据延迟导致“刚成交却显示无数据”
  - 护栏：Passport 明示“可能延迟 1-2 分钟，可刷新”


# PRD: Homepage（ERC-8004 Reputation）+ ERC-8183 Passport（可验证证据链）

创建日期：2026-04-20 19:00  
状态：Draft  
范围：前端 homepage + agent detail 的 passport view（不含平台接入）  

## 1. 背景与问题

### 1.1 背景

我们的核心工作是基于 **ERC-8004** 构建 agent reputation（综合画像），其中包含两类互补信号：

- **主观反馈（ERC-8004 Feedback）**：第三方提交的显式评价，擅长表达“体验/偏好/质量感知”。  
- **客观证据（ERC-8183 Commerce）**：链上可验证的作业终态与结算事实，擅长回答“是否真的完成、是否真的结算、是否可追溯”。  

在客观证据维度里，**ERC-8183 的链上经济行为与终态结算记录**是最关键、也最可信的信号来源之一。它比主观评价更难被女巫/刷分操纵，且具备可追溯、可验证、可结算、可迁移的特性。

### 1.2 问题

- 现有 homepage 难以在 10 秒内说明 “为什么这套 reputation 值得信”。
- Agent 开发者缺少一个**可分享、可匿名打开**的证明链接，用于降低客户/平台的第一次信任成本。
- 当前处于 pre-product 阶段，需求证据不足，过强口径（例如“唯一/主来源”）容易反噬。

## 2. 目标与非目标

### 2.1 目标（Goals）

- **G1 10 秒理解**：陌生人 10 秒内理解“这是 agent reputation（综合画像），其中 ERC-8183 提供可验证的客观证据维度，而不是主观评分替代品”。  
- **G2 证据链下钻**：从 homepage 可一键下钻到 Jobs / 链上 Tx / Passport view，完成“我自己验证”的路径。
- **G3 可分享 Passport**：Agent 开发者可分享公开链接，对方无需登录/连接钱包即可打开并验证（把证据压缩成一页）。
- **G4 不误导**：当 8183 数据不足时，不做 commerce 榜单叙事，不制造“top/trending”的错觉。

### 2.2 非目标（Non-goals，本期不做）

- 独立 `/passport/{uid}` 产品化页面（本期先做 `?view=passport`）。
- 平台真实接入排序系统、合作交付。
- curated samples、Proof-feed 首页大改（本期不做）。
- 隐私、争议仲裁、反作弊系统化方案（只做最小护栏）。

## 3. 用户与场景

### 3.1 目标用户（Primary persona）

**Agent 开发者 / 项目方老板 / 增长（卖方）**

- 想要：名声不被抄、能抗恶意差评、平台愿意给流量、客户敢用。
- 最怕：reputation 被刷、被恶评、平台不信、客户不敢试。

### 3.2 核心场景（Top scenarios）

- **S1 分享证明**：开发者把 Passport 链接发给客户/平台，对方点开就能验证。
- **S2 自己验证**：陌生用户从 homepage 点入 Jobs，看见真实终态记录与溯源，从“听你说”变成“我自己验证”。
- **S3 数据不足**：新 agent 没有足够终态 job 时，明确提示“数据不足 N/5”，避免误导。

## 4. 产品形态（本期交付）

### 4.1 Homepage（信息架构调整）

核心原则：**证据优先**，少自夸，多下钻。

#### 4.1.1 Hero（第一屏）

- 标题（中英双语，短句）：
  - 中文：**“Agent Reputation：主观反馈 + 可验证证据”**
  - 英文：**“Agent reputation: opinions + verifiable proof”**
- 副标题（中英双语）：
  - 中文：**“ERC-8183 把终态与结算压缩成 Passport：不是口碑分，是可点开验证的结果。”**
  - 英文：**“ERC-8183 turns outcomes into a portable passport.”**
- CTA（两个）：
  - CTA1：**View ERC-8183 Jobs** → `/jobs`
  - CTA2：**Create / Share Passport** → 引导到某个 agent 的 Passport view（见 4.2）

口径要求：本期采用**保守口径**，不写“唯一/主来源”，写“更可信的来源之一 / 可验证证据入口”。

#### 4.1.2 Proof blocks（第二屏，三张卡）

三张卡必须满足：**每张都有可点击验证落点**。

1. **Outcome-based reputation**

- 作用（证明什么）：把“是否真的完成/是否真的结算”从一句话，变成可验证的终态事实（Completed/Rejected/Expired）。
- 如何验证（点进去看到什么）：进入 Jobs，默认先看终态记录；每条 job 可继续下钻到 job 详情与链上 tx（浏览器链接）。
- 如何服务主线 reputation：与 ERC-8004 的主观 feedback 互补——feedback 解释体验与偏好，终态证据解释“是否履约/是否结算”。
- 点击：跳 `/jobs`（预设过滤：只看终态 action）

1. **Sybil-resistant by design**

- 作用（证明什么）：说明“主观评分容易被刷”，而 ERC-8183 的终态与资金事实需要走状态机与结算流程，具备天然抗女巫属性。
- 如何验证（点进去看到什么）：进入 Jobs/Job detail，通过终态路径与资金事件（结算 tx）验证“结果确实发生过”。  
- 如何服务主线 reputation：给 ERC-8004 feedback 一个可信锚点——当主观口碑与客观终态冲突时，用户至少能独立核验客观部分。
- 最小验收：必须有一个极简对比图（主观评分 vs 终态结算）
- 点击：跳 `/jobs`（预设过滤）

1. **Portable proof**

- 作用（证明什么）：把复杂证据链压缩成一页 Passport，适合“发给客户/平台”的低摩擦验证入口。
- 如何验证（点进去看到什么）：打开 Passport view 即可看到 SR/WS/CF 与最近终态摘要，并可一键跳 Jobs 与链上 tx 进一步核验。  
- 如何服务主线 reputation：Passport 是“综合画像”的可分享切片——主观 feedback 仍然存在，但这里强调客观可验证部分，降低第一次信任成本。
- 点击：先跳到说明文档（本期不绑定真实示例 agent，避免环境不稳定/无数据导致误导）

#### 4.1.3 Trending（第三屏）

本期固定策略：

- 仅保留现有 Feedback trending（或现状 leaderboard）
- **隐藏 Trending by Commerce**

### 4.2 Passport view（本期核心交付）

实现方式：在现有 agent detail 路由上增加视图模式。

#### 4.2.1 URL

`/agents/{uid}?view=passport`

#### 4.2.2 展示内容（MVP，≤ 8 项）

1. Agent 名称 + 头像 + chain name + chain logo + agent_id
2. **SR**（success rate）  
3. **WS**（weighted score）  
4. **CF**（confidence，显示 N/5）  
5. 最近 3 条终态记录摘要（Completed/Rejected/Expired）  
6. View jobs（跳 `/jobs` 或 job detail）  
7. View onchain tx（区块浏览器链接）  
8)（可选）跳到 Feedback 概览（提示“主观反馈仍在，但这里是客观证据”）

#### 4.2.2.1 指标释义（SR / WS / CF）

这三个指标来自 ERC-8183 的客观证据维度（Commerce Score），用于让用户在 10 秒内读懂“履约表现 + 金额权重 + 样本可靠性”。

- **SR（Success Rate，成功率）**：终态结果里成功（Completed）的占比。  
  - 直觉：回答“这类角色的 job，最终成功收尾的比例有多高”。  
  - 口径（Provider 视角）：`completed_count / terminal_count`，其中 `terminal_count = completed_count + rejected_count + expired_responsible_count`。
- **WS（Weighted Score，金额加权得分）**：把每个终态按金额加权后的结果分。金额越大的 job 对得分影响越大。  
  - 直觉：区分“做过大单且完成”与“小单刷量”。  
  - 口径（USD 版本，跨 token 可比）：`Σ(budget_usd × outcome_score) / Σ(budget_usd)`（若无 USD，则退化为 token 原生金额加权）。  
- **CF（Confidence，置信度，显示 N/5）**：用终态样本量提示分数可靠性。  
  - 直觉：回答“这套 SR/WS 是基于多少笔已结算/已终态的 job 得出的”。  
  - 展示：`N = min(terminal_job_count, 5)`，显示为 **N/5**；当 N < 5 时，必须文案提示“数据不足”。  

补充（已拍板）：

- Passport 内允许 **按 role 切换展示**（provider / client / evaluator），每个 role 的 SR/WS/CF 均来自对应的 `CommerceScore`。

#### 4.2.3 数据策略（弱一致）

- stale time：**30 秒**
- 接口失败策略：若有缓存先展示缓存，顶部红色 banner 提示“数据暂不可用/可能过期”，提供 Retry。
- 数据延迟提示：当显示“暂无终态 job”或“数据不足”时，提示“数据可能延迟 1-2 分钟，可手动刷新”。

#### 4.2.4 OG（分享卡片）

本期降级为站点级默认 OG（不按 agent 动态生成）。
原因：当前 Passport 以 `?view=passport` 方式挂在 client page 上，实现动态 OG 需要额外的服务端 metadata 层，本期先不引入。

### 4.3 Jobs（证据入口）

本期不改 Jobs 逻辑，但必须保证：

- homepage CTA 与 proof cards 能稳定导到 Jobs
- Jobs 页面能让用户看到终态证据并继续下钻

## 5. 规则与阈值（避免误导）

### 5.1 Confidence 规则

- **定义**：`confidence = min(terminal_job_count / 5, 1.0)`（阈值 THRESHOLD=5）。  
- **UI 映射**：`N = min(terminal_job_count, 5)`，展示为 `N/5`。  
- **约束**：当 `terminal_job_count < 5` 时，必须显示“数据不足 N/5”，且不得宣称 top/trending 或暗示排名意义。

### 5.2 Commerce Trending（本期固定）

- 本期默认隐藏 commerce trending，不做 curated samples。

## 6. 验收标准（Acceptance Criteria）

1. Passport 无数据

- Given：terminal=0  
- When：打开 passport view  
- Then：显示“暂无终态 job…”，不显示榜单位次，可显示 CF=0/5，并提供引导入口

1. Passport 数据不足

- Given：terminal=3  
- Then：显示“数据不足 3/5”，不得出现榜单相关信息

1. 接口失败可见

- When：scores 接口失败或 code!=0  
- Then：显示缓存（如有）+ 红色 banner + Retry

1. 匿名访问

- When：无登录/无钱包连接打开 passport 链接  
- Then：可渲染，证据链接可点开

1. 首页证据链可走通

- When：从 homepage 点 CTA/三张 proof cards  
- Then：到 Jobs 或 passport，并可下钻到链上 tx

## 7. 指标（可先定义）

- Hero CTA 点击率（Jobs / Passport）
- Proof cards 点击率（3 张分别统计）
- Jobs 停留时长、滚动深度
- Passport 复制/分享次数（至少复制按钮点击）
- Passport 错误率（banner 触发比例）

## 8. 风险与护栏

### 8.1 风险

- 公开分享页必然被爬，可能造成接口压力。
- 数据延迟可能造成“刚成交就分享却显示无数据”的挫败体验。

### 8.2 护栏（本期必须）

- Passport view 中所有用户/agent 文本按纯文本渲染，禁止 HTML 注入。
- 最小限流/反爬策略（网关或接口层）。

## 9. 发布策略（Rollout）

- 先上线 passport view（内部可访问），验证匿名打开与错误态。
- 再上线 homepage 新 Hero + proof cards，把流量导入证据链。
- 一周内收集 5 个开发者访谈原话，迭代文案与证据入口。


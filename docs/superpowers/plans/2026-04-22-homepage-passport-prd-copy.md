# Homepage + Passport（ERC8004 Reputation + ERC8183 Proof）Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 按最新 PRD 文案口径，把 Homepage 与 Passport 的叙事从“只讲 ERC-8183”对齐为“ERC-8004 reputation 主线 + ERC-8183 可验证证据维度”，并让 Proof blocks 与 SR/WS/CF 在 UI 中可解释、可验证、可分享。

**Architecture:** 不改路由形态：Homepage 仍在 `/`，Passport 仍是 `/agents/{uid}?view=passport`。实现以“文案+可解释信息”为核心：更新 Hero/Proof cards 的文案；在 Commerce 分数卡上增加 SR/WS/CF/Confidence 的 tooltip 解释；保持点击下钻到 `/jobs` 不变。

**Tech Stack:** Next.js App Router, React 19, Tailwind, MUI, TanStack React Query, Vitest + React Testing Library

---

## Scope lock (DO / DON'T)

**DO (本计划范围内):**
- Homepage Hero 文案对齐 PRD（突出“ERC-8004 主线 + ERC-8183 证据”）
- Proof blocks 三卡文案对齐 PRD（每卡说明“证明什么/如何验证”）
- Passport view：确保用户能读懂 **SR / WS / CF(N/5)** 的含义（tooltip/文案）
- 更新对应 Vitest 用例，确保 `yarn test` 通过

**DON'T (明确不做):**
- 不新增独立 `/passport/{uid}` 路由
- 不引入站内 docs 路由（第三张 proof card 仍按 Phase1 走 `/agents` 作为入口）
- 不做动态 OG
- 不做 commerce trending

---

## File map (create/modify)

**Modify (PRD 已完成，无需再动):**
- `erc8004/docs/designs/202604201900_homepage_erc8183_passport_prd.md`

**Modify (Homepage):**
- `erc8004-1.1-web/src/app/(app)/(home)/components/TopMain.tsx`（Hero 标题/副标题）
- `erc8004-1.1-web/src/app/(app)/(home)/page.tsx`（Proof blocks 三卡说明文案）

**Modify (Passport / Commerce score copy):**
- `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.tsx`（空态提示文案）
- `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceScoreCard.tsx`（SR/WS/CF tooltip 与 CF(N/5) 文案）

**Tests:**
- `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.test.tsx`
- `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceScoreCard.test.tsx`

---

### Task 1: Homepage Hero 文案对齐“ERC-8004 主线 + ERC-8183 证据”

**Files:**
- Modify: `erc8004-1.1-web/src/app/(app)/(home)/components/TopMain.tsx:39-44`

- [ ] **Step 1: 更新 Hero 标题/副标题（中英双语）**

将：

```tsx
<div className="text-[46px] ...">
  8183：用可验证成交记录，构建 Agent 商业信用
</div>
<div className="mt-2 text-[16px] ...">
  不是评分。是结算过的结果，所有人都能点开验证。
</div>
```

替换为（保持现有 UI 结构，新增一行英文的小字作为“bilingual”最小实现）：  

```tsx
<div className="text-[46px] tracking-[0.92px] leading-[1.3] md-max:text-[26px] md-max:leading-[1.35] md-max:tracking-[0.5px]">
  Agent Reputation：主观反馈 + 可验证证据
</div>
<div className="mt-2 text-[16px] tracking-[0.44px] leading-[1.8] md-max:text-[14px] md-max:leading-[1.5] md-max:tracking-[0.2px]">
  ERC-8183 把终态与结算压缩成 Passport：不是口碑分，是可点开验证的结果。
</div>
<div className="mt-1 text-[12px] tracking-[0.2px] opacity-90 md-max:text-[12px]">
  Agent reputation: opinions + verifiable proof
</div>
```

- [ ] **Step 2: 运行本地开发查看首屏是否换行合理**

Run:

```bash
cd erc8004-1.1-web && yarn dev
```

Expected:
- Hero 标题/副标题更新
- CTA 不变：`View ERC-8183 Jobs` 仍跳 `/jobs`，`Create / Share Passport` 仍跳 `/agents`

- [ ] **Step 3: Commit**

```bash
git add erc8004-1.1-web/src/app/(app)/(home)/components/TopMain.tsx
git commit -m "docs(ui): align homepage hero with ERC8004+ERC8183 narrative"
```

---

### Task 2: Proof blocks 三卡文案对齐 PRD（证明点 + 可验证落点）

**Files:**
- Modify: `erc8004-1.1-web/src/app/(app)/(home)/page.tsx:45-80`

- [ ] **Step 1: 更新三张卡的描述文案为“可验证/可下钻”的明确表述**

将三张卡的 `text-bas-text-secondary` 描述替换为更贴近 PRD 的中文（保留标题英文不动，降低改动面）。

示例（卡 1）：  

```tsx
<div className="mt-2 text-[13px] text-bas-text-secondary">
  看终态（Completed/Rejected/Expired）与结算事实，把“听你说”变成“我自己验证”。
</div>
```

卡 2：  

```tsx
<div className="mt-2 text-[13px] text-bas-text-secondary">
  主观评分可被刷；终态与结算需走状态机与资金事实，更难伪造。
</div>
```

卡 3（仍跳 `/agents`，但文案解释“可分享 Passport”）：  

```tsx
<div className="mt-2 text-[13px] text-bas-text-secondary">
  把证据压缩成一页 Passport（匿名打开），一键跳 Jobs 与链上 Tx 继续核验。
</div>
```

- [ ] **Step 2: 运行 dev，确认三卡仍可点击**

Run:

```bash
cd erc8004-1.1-web && yarn dev
```

Expected:
- 卡 1/2 仍跳 `/jobs`
- 卡 3 仍跳 `/agents`（Phase 1 降级策略不变）

- [ ] **Step 3: Commit**

```bash
git add erc8004-1.1-web/src/app/(app)/(home)/page.tsx
git commit -m "docs(ui): clarify proof blocks value and verification path"
```

---

### Task 3: Passport/Commerce 分数解释（SR/WS/CF(N/5)）在 UI 可读

**Files:**
- Modify: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.tsx:38-48`
- Modify: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceScoreCard.tsx:17-45,48-70`
- Test: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.test.tsx`
- Test: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceScoreCard.test.tsx`

- [ ] **Step 1: Passport 空态文案对齐 PRD（提示可下钻 + 可延迟）**

把 `PassportView.tsx` 空态从：

```tsx
<div>暂无终态 job</div>
```

改为：

```tsx
<div>暂无终态 job（数据可能延迟 1-2 分钟，可手动刷新）</div>
```

- [ ] **Step 2: 为 SR / WS 增加 tooltip（Provider 卡）**

在 `CommerceScoreCard.tsx` 的 ProviderCard，把标题从纯文本换为带 tooltip 的 `StatItem` 风格（或直接在现有 span 上包 `BaseToolTip`）。保持改动最小，推荐直接替换 label 行为：

```tsx
<div className="flex items-center gap-1 text-[13px] text-bas-text-secondary">
  <span>SR（Success Rate）</span>
  <BaseToolTip content="终态结果里 Completed 的占比。直觉：回答“最终成功收尾的比例有多高”。">
    <span className="cursor-help text-[10px] text-bas-text-secondary">?</span>
  </BaseToolTip>
</div>
```

对 WS 类似：

```tsx
<BaseToolTip content="金额加权得分：Σ(budget_usd × outcome_score) / Σ(budget_usd)。金额越大的单影响越大，用于区分大单履约与小单刷量。">
  <span className="cursor-help text-[10px] text-bas-text-secondary">?</span>
</BaseToolTip>
```

并把显示名称从 `Success Rate` / `Weighted Score` 更新为 `SR（Success Rate）` / `WS（Weighted Score）` 以对齐 PRD。

- [ ] **Step 3: CF(N/5) 展示与 tooltip**

在 `ConfidenceBadge` 中，把英文 badge 改为中文 + N/5（逻辑沿用 `Math.round(confidence*5)`，与 `confidence = min(terminal/5,1)` 一致）：

```tsx
return (
  <BaseToolTip content="置信度：基于终态 job 样本量的可靠性提示。N=终态 job 数（上限 5），显示为 N/5。">
    <span className="text-[11px] text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full cursor-help">
      数据不足（{done}/{totalNeeded}）
    </span>
  </BaseToolTip>
);
```

并同步更新 `CommerceScoreCard.test.tsx` 断言文本。

- [ ] **Step 4: 更新测试用例（先红后绿）**

1) 更新 `CommerceScoreCard.test.tsx`：  

```tsx
expect(screen.getByText('数据不足（3/5）')).toBeInTheDocument();
```

2) 更新 `PassportView.test.tsx`：空态断言改为包含延迟提示即可：  

```tsx
expect(screen.getByText(/暂无终态 job/i)).toBeInTheDocument();
```

- [ ] **Step 5: 运行测试**

Run:

```bash
cd erc8004-1.1-web && yarn test
```

Expected:
- PASS（所有现有 tests + 本次改动相关 tests）

- [ ] **Step 6: Commit**

```bash
git add \
  erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.tsx \
  erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceScoreCard.tsx \
  erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.test.tsx \
  erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceScoreCard.test.tsx
git commit -m "docs(ui): explain SR/WS/CF and show confidence as N/5"
```

---

## Self-Review

**1) Spec coverage（对齐 PRD）**
- Homepage：Hero 与 Proof blocks 已明确“ERC-8004 主线 + ERC-8183 证据”叙事，并且每卡都指向可验证入口（`/jobs` 或 `/agents`）。
- Passport：SR/WS/CF 的定义已通过 tooltip/文案露出；CF 显示为 N/5。

**2) Placeholder scan**
- 全文无 TBD/TODO/“自行补充”。

**3) Type consistency**
- 只改 UI 文案与 tooltip，不改 API 字段命名：`success_rate` / `weighted_score` / `confidence` 沿用现有 `CommerceScore`。

---

## Execution Handoff

Plan complete and saved to `erc8004/docs/superpowers/plans/2026-04-22-homepage-passport-prd-copy.md`. Two execution options:

1. **Subagent-Driven (recommended)** - per-task dispatch + review between tasks  
2. **Inline Execution** - execute tasks in this session with checkpoints

Which approach?


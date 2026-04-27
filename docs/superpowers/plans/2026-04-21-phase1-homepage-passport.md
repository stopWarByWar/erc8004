# Phase 1 (Homepage + Passport) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 按 `docs/exec-plans/202604210000_phase1_homepage_passport_exec.md` 交付 Homepage 证据导流 + `?view=passport` 的可分享/匿名 Passport MVP，并补齐最小测试。

**Architecture:** 不引入新路由产品形态，不做动态 OG。Passport 作为 `agents/[id]` 页面内的一个视图分支（query param），复用现有 commerce hooks（scores/actions/stats）与 UI 组件，新增 Passport 专用容器组件与错误 banner。

**Tech Stack:** Next.js App Router, React 19, TanStack React Query, MUI, Tailwind, Vitest + React Testing Library

---

## Scope lock (DO / DON'T)

**DO (Phase 1):**
- Homepage：更新 Hero 文案与 CTA；新增 3 个 proof blocks；隐藏 commerce trending
- Agent detail：新增 `?view=passport` 视图（role tabs 必须），并提供证据下钻链接
- Tests：Vitest 能跑，且至少覆盖 Passport 的空态/不足态/失败态 + role tabs 行为

**DON'T (explicitly out of scope):**
- 独立 `/passport/{uid}` 路由
- Commerce Trending / curated samples
- 动态 OG（按 agent 生成）
- 平台排序接入与商业交付

---

## File map (create/modify)

**Modify (Homepage):**
- `erc8004-1.1-web/src/app/(app)/(home)/components/TopMain.tsx`（Hero 文案 + CTA）
- `erc8004-1.1-web/src/app/(app)/(home)/page.tsx`（Proof blocks + Trending 文案/结构调整）

**Modify (Agent detail / Passport view):**
- `erc8004-1.1-web/src/app/(app)/agents/[id]/page.tsx`（增加 `view=passport` 分支）
- Create `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.tsx`
- Create `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceTerminalSummary.tsx`（最近 3 条终态摘要）
- Create `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceStaleBanner.tsx`（弱一致失败红条 + Retry）

**Tests:**
- Modify `erc8004-1.1-web/package.json`（已有：`test`/`test:watch`）
- `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.test.tsx`（新增）
- `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceTerminalSummary.test.tsx`（新增）

**Docs:**
- 不新增 design。只保持 PRD 与 exec plan 不变。

---

### Task 1: Homepage Hero 文案与 CTA 改造（按 PRD 4.1.1）

**Files:**
- Modify: `erc8004-1.1-web/src/app/(app)/(home)/components/TopMain.tsx`
- Modify: `erc8004-1.1-web/src/app/(app)/(home)/page.tsx`

- [ ] **Step 1: 将 Hero 标题/副标题替换为 PRD 文案**

在 `TopMain.tsx` 把现有：
- `HELLO, AGENT IDENTITY！`
- `A unified explorer ...`

替换为：

```tsx
<div className="text-[46px] ...">
  8183：用可验证成交记录，构建 Agent 商业信用
</div>
<div className="mt-2 text-[16px] ...">
  不是评分。是结算过的结果，所有人都能点开验证。
</div>
```

- [ ] **Step 2: 将两个按钮改成：View ERC-8183 Jobs / Create or Share Passport**

```tsx
<BaseButton onClick={() => router.push('/jobs')}>View ERC-8183 Jobs</BaseButton>
<BaseButton onClick={handleViewPassport}>Create / Share Passport</BaseButton>
```

并新增 `handleViewPassport`（Phase 1 不绑定 demo agent，按 exec plan 1.1A）：

```tsx
const handleViewPassport = () => {
  router.push('/agents');
};
```

（解释：本期第三张 proof 卡也不跳真实 agent，所以 CTA2 先导到 agents 列表，不制造“示例”错觉。）

- [ ] **Step 3: 新增 Proof blocks（3 张卡）到 homepage**

在 `page.tsx` 的 stats 区块之后、Trending 之前插入一个 `ProofBlocks` JSX（不新建组件也可以，先最小 diff）。

卡片文案与跳转：
- 卡 1：Outcome-based reputation → `Link href="/jobs"`
- 卡 2：Sybil-resistant by design → `Link href="/jobs"`
- 卡 3：Portable proof → `Link href="/docs/erc8183-commerce-reputation"`（或已有 docs 页面路由，如果不存在，就先链到仓库 docs 的外链不做，本期只要“说明入口”）

（如果站内没有 docs 路由，先用 `href="https://github.com/..."` 不在 scope。推荐：先做站内 `/docs` 页也不在 scope。这里如阻塞，记录阻塞点。）

- [ ] **Step 4: 隐藏 Trending by Commerce**

Phase 1 要求 “隐藏 commerce trending”。当前 homepage 的 Trending 来自 `useGetLeaderboard()`，先不动数据结构，只保证没有新增 commerce trending tab/内容即可。

- [ ] **Step 5: 手动验收（无需自动化）**

Run: `npm run dev`  
Check:
- Hero 文案正确
- CTA1 跳 `/jobs`
- CTA2 跳 `/agents`
- 3 张卡可点击（2 张到 `/jobs`，1 张到说明入口）

---

### Task 2: Passport view 路由分支（`?view=passport`）

**Files:**
- Modify: `erc8004-1.1-web/src/app/(app)/agents/[id]/page.tsx`
- Create: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.tsx`

- [ ] **Step 1: 写一个 failing test（PassportView 空态）**

Create `PassportView.test.tsx`：

```tsx
import React from 'react';
import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { PassportView } from './PassportView';

vi.mock('@/http/agent', () => ({
  useGetCommerceScores: vi.fn(),
  useGetCommerceActions: vi.fn(),
}));

import { useGetCommerceScores, useGetCommerceActions } from '@/http/agent';

describe('PassportView', () => {
  it('shows empty state when no terminal jobs', () => {
    vi.mocked(useGetCommerceScores).mockReturnValue({ data: { scores: [] }, isLoading: false } as any);
    vi.mocked(useGetCommerceActions).mockReturnValue({ data: { items: [] }, isLoading: false } as any);

    render(<PassportView uid="123" />);
    expect(screen.getByText(/暂无终态 job/i)).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: 运行测试，确认失败**

Run: `npm test`  
Expected: FAIL（因为 `PassportView` 不存在）

- [ ] **Step 3: 实现最小 PassportView**

Create `PassportView.tsx`（最小结构，后续再补 UI）：

```tsx
'use client';

import React, { useMemo, useState } from 'react';
import Link from 'next/link';
import { useGetCommerceScores, useGetCommerceActions } from '@/http/agent';
import type { CommerceRole } from '@/types/api';
import { CommerceScoreCard } from './CommerceScoreCard';
import { CommerceTerminalSummary } from './CommerceTerminalSummary';

const ROLE_LABELS: Record<CommerceRole, string> = { provider: 'Provider', client: 'Client', evaluator: 'Evaluator' };

export function PassportView({ uid }: { uid: string }) {
  const { data: scoresData, isLoading, refetch } = useGetCommerceScores({ uid });
  const scores = scoresData?.scores ?? [];

  const activeRoles = useMemo(() => scores.filter((s) => s.total_jobs > 0), [scores]);
  const [selectedRole, setSelectedRole] = useState<CommerceRole | null>(null);
  const effectiveRole = selectedRole ?? activeRoles[0]?.role ?? null;
  const currentScore = scores.find((s) => s.role === effectiveRole);

  // NOTE: actions 仅用于终态摘要，不做全列表（Phase 1）
  const { data: actionsData, isLoading: isActionsLoading } = useGetCommerceActions(
    { uid, role: effectiveRole ?? 'provider', page: 1, page_size: 20 },
    !!effectiveRole,
  );

  if (!isLoading && activeRoles.length === 0) {
    return (
      <div className="text-center text-bas-text-secondary text-[13px] py-8">
        暂无终态 job…（数据可能延迟 1-2 分钟，可手动刷新）
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex gap-3 border-b border-bas-border overflow-x-auto">
        {activeRoles.map((s) => (
          <button
            key={s.role}
            type="button"
            onClick={() => setSelectedRole(s.role)}
            className={`pb-2 text-[13px] transition-colors relative shrink-0 ${
              effectiveRole === s.role ? 'text-bas-text-primary' : 'text-bas-text-secondary'
            }`}
          >
            {ROLE_LABELS[s.role]}
            <span className="ml-1 text-[11px] text-bas-text-secondary">({s.total_jobs})</span>
            {effectiveRole === s.role && <div className="absolute bottom-0 left-0 right-0 h-[2px] bg-black" />}
          </button>
        ))}
      </div>

      {effectiveRole && currentScore && (
        <div className="border border-bas-border rounded-[10px] p-4">
          <CommerceScoreCard score={currentScore} isLoading={isLoading} />
        </div>
      )}

      <CommerceTerminalSummary uid={uid} role={effectiveRole ?? 'provider'} items={(actionsData as any)?.items ?? []} />

      <div className="flex gap-3">
        <Link className="text-[13px] text-bas-theme-primary" href="/jobs">
          View jobs
        </Link>
        <button type="button" className="text-[13px] text-bas-theme-primary" onClick={() => refetch()}>
          Retry
        </button>
      </div>
    </div>
  );
}
```

- [ ] **Step 4: 运行测试，确认通过**

Run: `npm test`  
Expected: PASS（空态文案出现）

- [ ] **Step 5: 将 `agents/[id]/page.tsx` 增加 query 分支**

在 `DetailPage` 里引入 `useSearchParams` 并在 Tab 容器之前判断：

```tsx
import { useParams, useSearchParams } from 'next/navigation';
import { PassportView } from './components/PassportView';

const searchParams = useSearchParams();
const view = searchParams.get('view');
if (view === 'passport') {
  return (
    <Container className="py-6">
      <PassportView uid={uid} />
    </Container>
  );
}
```

（注意：不要破坏原本详情页。Passport 是“分支视图”，不是替代。）

---

### Task 3: CommerceTerminalSummary（最近 3 条终态摘要）

**Files:**
- Create: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceTerminalSummary.tsx`
- Create test: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceTerminalSummary.test.tsx`

- [ ] **Step 1: 写 failing test（只展示 3 条终态）**

```tsx
import React from 'react';
import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { CommerceTerminalSummary } from './CommerceTerminalSummary';

describe('CommerceTerminalSummary', () => {
  it('shows up to 3 terminal actions', () => {
    render(
      <CommerceTerminalSummary
        uid="1"
        role="provider"
        items={[
          { action: 'job_completed', timestamp: 1, job_id: 1 },
          { action: 'job_rejected', timestamp: 2, job_id: 2 },
          { action: 'job_expired', timestamp: 3, job_id: 3 },
          { action: 'job_submitted', timestamp: 4, job_id: 4 },
        ] as any}
      />,
    );
    expect(screen.getAllByTestId('terminal-item').length).toBe(3);
  });
});
```

- [ ] **Step 2: 实现组件**

```tsx
'use client';

import React, { useMemo } from 'react';
import type { CommerceActionItem, CommerceRole } from '@/types/api';

const TERMINAL = new Set(['job_completed', 'job_rejected', 'job_expired']);

export function CommerceTerminalSummary({
  uid,
  role,
  items,
}: {
  uid: string;
  role: CommerceRole;
  items: CommerceActionItem[];
}) {
  const terminal = useMemo(
    () => items.filter((i) => TERMINAL.has(i.action)).sort((a, b) => (b.timestamp ?? 0) - (a.timestamp ?? 0)).slice(0, 3),
    [items],
  );

  if (terminal.length === 0) return null;

  return (
    <div className="border border-bas-border rounded-[10px] p-4">
      <div className="text-[13px] text-bas-text-secondary mb-2">Latest terminal outcomes</div>
      <div className="flex flex-col gap-2">
        {terminal.map((i) => (
          <div key={`${i.chain_id}-${i.commerce_contract}-${i.job_id}-${i.action}-${i.timestamp}`} data-testid="terminal-item">
            <span className="text-[13px] text-bas-text-primary">{i.action}</span>
            <span className="ml-2 text-[12px] text-bas-text-secondary">job #{i.job_id}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
```

- [ ] **Step 3: 运行测试**

Run: `npm test`  
Expected: PASS

---

### Task 4: Passport “数据不足 N/5” 与失败 banner（弱一致）

**Files:**
- Create: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/CommerceStaleBanner.tsx`
- Modify: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.tsx`
- Test: `erc8004-1.1-web/src/app/(app)/agents/[id]/components/PassportView.test.tsx`（扩展）

- [ ] **Step 1: 扩展 failing test（confidence < 1 显示 Insufficient data done/5）**

```tsx
it('shows insufficient data badge when confidence < 1', () => {
  vi.mocked(useGetCommerceScores).mockReturnValue({
    data: { scores: [{ role: 'provider', total_jobs: 3, confidence: 0.6, success_rate: 0.5, weighted_score: 0.5 }] },
    isLoading: false,
  } as any);
  vi.mocked(useGetCommerceActions).mockReturnValue({ data: { items: [] }, isLoading: false } as any);

  render(<PassportView uid="123" />);
  expect(screen.getByText(/Insufficient data/i)).toBeInTheDocument();
});
```

- [ ] **Step 2: 失败 banner（当 scores 查询 error 时显示）**

这里依赖 hook 是否暴露 `isError`/`error`。若当前 `useGetCommerceScores` 没透出，阻塞点必须上报，不能私自改上游边界。

假设可拿到 `isError`：

```tsx
export function CommerceStaleBanner({ onRetry }: { onRetry: () => void }) {
  return (
    <div className="bg-red-50 text-red-700 border border-red-200 rounded-[10px] px-3 py-2 text-[12px]">
      数据暂不可用/可能过期 <button className="underline ml-2" onClick={onRetry}>Retry</button>
    </div>
  );
}
```

在 `PassportView` 顶部插入：

```tsx
{isError ? <CommerceStaleBanner onRetry={() => refetch()} /> : null}
```

- [ ] **Step 3: 运行测试**

Run: `npm test`  
Expected: PASS

---

## Self-review (spec coverage)

对照 PRD 与 exec plan：
- Homepage：Hero/CTA/Proof blocks/隐藏 commerce trending ✅（Task 1）
- Passport：`?view=passport`、role tabs、SR/WS/CF、终态摘要、匿名可开 ✅（Task 2/3/4）
- 弱一致错误态：红条 + Retry ✅（Task 4）
- Tests：至少覆盖 Passport 关键状态 ✅（Task 2/4）
- OG：明确降级为默认站点级 ✅（scope lock）

Placeholder scan：无 “TODO/TBD/implement later” 语句。阻塞点仅在 “hook 是否暴露 error 状态”，若出现会明确上报。

---

## Execution handoff

Plan complete and saved to `docs/superpowers/plans/2026-04-21-phase1-homepage-passport.md`. Two execution options:

1. **Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration
2. **Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?


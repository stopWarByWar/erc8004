---
description: 更新 ERC-8183 AgenticCommerce 合约源码（仅在有更新时同步并生成变更记录）
---

# 更新 ERC-8183 AgenticCommerce 合约源码

检查 ERC-8183 AgenticCommerce 合约源码是否有更新，有更新则同步到本地并生成变更记录。

## 数据源

- **合约源码**：https://github.com/erc-8183/base-contracts/blob/main/contracts/AgenticCommerce.sol

## 清理旧数据

在同步前，先清理 `docs/erc-8183/` 目录下除 `.sync-sha` 之外的所有旧文件（changelog 和旧版源码），确保目录只保留最新的合约和当前的 changelog。

## 执行步骤

### Step 0：清理旧数据

1. 读取 `docs/erc-8183/.sync-sha`（如果存在），记录其内容备用。
2. 删除 `docs/erc-8183/` 目录下除 `.sync-sha` 之外的所有文件。

### Step 1：获取最新 commit SHA（轻量检查，无需下载内容）

运行以下 shell 命令：

```bash
curl -s "https://api.github.com/repos/erc-8183/base-contracts/commits?path=contracts/AgenticCommerce.sol&per_page=1" | python3 -c "import sys,json; print(json.load(sys.stdin)[0]['sha'])"
```

将结果记为 `REMOTE_SHA`。

### Step 2：读取本地已记录的 SHA

读取文件 `docs/erc-8183/.sync-sha`：
- 如果文件不存在，视为首次同步，继续执行 Step 3。
- 如果文件存在，读取内容记为 `LOCAL_SHA`。
  - 若 `REMOTE_SHA == LOCAL_SHA`，**直接输出"AgenticCommerce.sol 无更新"并停止，不做任何其他操作**。
  - 若不同，继续执行 Step 3。

### Step 3：下载最新文件内容

```bash
curl -s "https://raw.githubusercontent.com/erc-8183/base-contracts/main/contracts/AgenticCommerce.sol"
```

将内容记为 `NEW_CONTENT`。

### Step 4：生成 diff

- 读取本地现有文件 `docs/erc-8183/AgenticCommerce.sol`（若不存在则视为空文件）。
- 将旧内容与 `NEW_CONTENT` 进行 diff，提取**变更部分**（新增、删除、修改的段落），记为 `DIFF_CONTENT`。

### Step 5：生成 changelog 文件

**仅基于 `DIFF_CONTENT`**（不读全文），生成一份变更记录文档：

- **文件名**格式：`YYYYMMDD-<更新大纲>.md`
  - 日期取今天的日期（本地时间）
  - 更新大纲：用 3-8 个中文字概括本次核心变更内容，用连字符连接（例如：`20260408-新增预算审计事件.md`）
- **文件内容**格式：

```markdown
# AgenticCommerce.sol 更新记录 - YYYY-MM-DD

## 变更摘要
（1-2 句话概括本次更新的核心内容）

## 详细变更

### 新增内容
（列出新增的函数、事件、修改器等）

### 修改内容
（列出有实质性变化的代码段）

### 删除内容
（列出被移除的内容，若无则省略此节）

## 对应 commit
`<REMOTE_SHA>`
```

- 将此文件保存到 `docs/erc-8183/` 目录。

### Step 6：更新本地文件

1. 确保目录 `docs/erc-8183/` 存在（若不存在则创建）。
2. 将 `NEW_CONTENT` 写入 `docs/erc-8183/AgenticCommerce.sol`，覆盖旧内容。
3. 将 `REMOTE_SHA` 写入 `docs/erc-8183/.sync-sha`，覆盖旧内容。

### Step 7：输出结果

向用户报告：
- 已同步的 commit SHA
- 新建的 changelog 文件名
- 关键变更摘要（一句话）
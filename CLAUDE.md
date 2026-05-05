# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go backend service implementing the ERC-8004 standard for AI agent identity registration and validation. It indexes blockchain events from smart contracts, stores agent identity/reputation data in PostgreSQL, and serves a REST API. Two separate processes: **server** (REST API) and **indexer** (blockchain event listener).

## Commands

```bash
# Build
go build ./server/main.go
go build ./indexer/main.go

# Run
go run ./server/main.go -f ./server/api/logic/config.yaml
go run ./indexer/main.go -f ./config/conf.yaml
./start_indexer.sh   # interactive multi-chain indexer launcher

# CLI tool - batch update vector embeddings
go run ./cmd/update_desc_vector.go -f ./config/conf.yaml

# Test
go test ./...
go test -v -run TestName ./path/to/package/
go test -cover ./...
go test -race ./...
```

## Architecture

### Layered Request Flow
```
HTTP Request → handle/ → logic/ → model/ → PostgreSQL
```
- `server/api/handle/` — HTTP parsing, validation, calling logic functions
- `server/api/logic/` — business logic, data transformation, response shaping
- `model/` — GORM-based data access; queries split by domain (`identity.go`, `reputation.go`, `validation.go`, `embedding.go`)

**Server** (`server/main.go`): Gin-based REST API. Routes defined in `server/api/router.go`. All handler logic lives in `handle/` + `logic/`.

**Indexer** (`indexer/main.go`): Polls blockchain events from smart contracts. Processors in `indexer/processor/` handle different contract types (Identity, Reputation, Validation, Comments). Uses a fan-out channel pattern to distribute events.

### Key Packages
| Package | Purpose |
|---------|---------|
| `model/` | GORM models (`types.go`) + domain-specific query files |
| `config/` | Loads YAML config; `config.ChainMap` and `config.RegisterMap` drive multi-chain support |
| `abi/` | Generated Ethereum smart contract ABIs |
| `agentCard/` | Agent profile card data structures |
| `helper/` | AWS S3 file upload helper |
| `logger/` | Structured logging wrapper |

### Multi-Chain Support
All agents are scoped by `chain_id` + `identity_registry` address. Chain configs live in `config/testnet/` and `config/mainnet/` as individual YAML files, aggregated by `config/config.yaml`.

### Semantic Search
Agent descriptions are embedded via OpenAI API and stored in `agent_vectors` table (`vector(1536)` type via pgvector). The `cmd/update_desc_vector.go` tool batch-processes embeddings. Similarity queries run in `model/embedding.go`.

### Database
PostgreSQL with GORM. Schema versioned in `migrations/`. Initialize with `model.InitDB()`. All models defined in `model/types.go`; key tables: `agents`, `feedbacks`, `validations`, `agent_vectors`, `attestation`.

## Configuration

- Server runtime config: `server/api/logic/config.yaml` (DB, OpenAI key, AWS, etc.)
- Indexer/chain config: `config/conf.yaml` + per-chain files in `config/testnet/` or `config/mainnet/`
- Secrets (DB password, OpenAI key, AWS keys) are read from the YAML config — never hardcode them

## Work Flow
- 对于任何新的功能，先写设计文档，再写spec，最后写exec-plan
- 按照exec-plan执行，每执行完成一项任务，在对应的文档中将对应的任务勾选掉

## Skill routing

When the user's request matches an available skill, ALWAYS invoke it using the Skill
tool as your FIRST action. Do NOT answer directly, do NOT use other tools first.
The skill has specialized workflows that produce better results than ad-hoc answers.

Key routing rules:
- Product ideas, "is this worth building", brainstorming → invoke office-hours
- Bugs, errors, "why is this broken", 500 errors → invoke investigate
- Ship, deploy, push, create PR → invoke ship
- QA, test the site, find bugs → invoke qa
- Code review, check my diff → invoke review
- Update docs after shipping → invoke document-release
- Weekly retro → invoke retro
- Design system, brand → invoke design-consultation
- Visual audit, design polish → invoke design-review
- Architecture review → invoke plan-eng-review
- Save progress, checkpoint, resume → invoke checkpoint
- Code quality, health check → invoke health

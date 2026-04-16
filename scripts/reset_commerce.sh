#!/usr/bin/env bash
set -euo pipefail

# Reset Commerce tables, then re-apply Commerce migrations.
#
# Usage:
#   DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=disable" ./scripts/reset_commerce.sh
#
# Or with libpq env vars:
#   PGHOST=... PGPORT=... PGUSER=... PGPASSWORD=... PGDATABASE=... ./scripts/reset_commerce.sh
#
# Notes:
# - This script is destructive: it DROPs commerce tables.
# - It relies on psql being installed and reachable in PATH.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIG_DIR="$ROOT_DIR/migrations"

PSQL_ARGS=(-v ON_ERROR_STOP=1)

if [[ -n "${DATABASE_URL:-}" ]]; then
  PSQL_CONN=("$DATABASE_URL")
else
  PSQL_CONN=()
fi

echo "[1/3] Dropping commerce tables (CASCADE)..."
psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" <<'SQL'
BEGIN;

DROP TABLE IF EXISTS commerce_actions CASCADE;
DROP TABLE IF EXISTS commerce_scores CASCADE;
DROP TABLE IF EXISTS commerce_scores_global CASCADE;
DROP TABLE IF EXISTS commerce_jobs CASCADE;

COMMIT;
SQL

echo "[2/3] Re-applying commerce migrations..."

# Base schema + triggers (safe init, no DROP)
psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" -f "$MIG_DIR/202604071600_commerce_init_safe.psql"

# Jobs snapshot table
psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" -f "$MIG_DIR/202604111000_commerce_jobs.sql"

# Token-aware fields + indexes for commerce_actions/scores
psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" -f "$MIG_DIR/202604101430_commerce_token_aware.sql"

# Jobs USD/payment breakdown columns
psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" -f "$MIG_DIR/202604131300_commerce_jobs_payment_breakdown.sql"

# payment_decimals columns
psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" -f "$MIG_DIR/202604131600_commerce_payment_decimals.sql"

echo "[2/3] Applying commerce query indexes (extra)..."
# NOTE: we intentionally do NOT run the merged optimization migration here, because it also
# adds non-commerce constraints (e.g. agents unique keys) which may fail if the DB has
# pre-existing duplicate data. This script is commerce-only.
psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" <<'SQL'
DO $$
BEGIN
  -- commerce_jobs: stable pagination + common filters
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_cj_chain_contract_updated_job') THEN
    CREATE INDEX idx_cj_chain_contract_updated_job
      ON commerce_jobs(chain_id, commerce_contract, updated_at DESC, job_id DESC);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_cj_chain_contract_status_updated_job') THEN
    CREATE INDEX idx_cj_chain_contract_status_updated_job
      ON commerce_jobs(chain_id, commerce_contract, status, updated_at DESC, job_id DESC);
  END IF;

  -- Role-specific filters in Job Browser (client/provider views).
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_cj_chain_contract_client_updated_job') THEN
    CREATE INDEX idx_cj_chain_contract_client_updated_job
      ON commerce_jobs(chain_id, commerce_contract, client, updated_at DESC, job_id DESC);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_cj_chain_contract_provider_updated_job') THEN
    CREATE INDEX idx_cj_chain_contract_provider_updated_job
      ON commerce_jobs(chain_id, commerce_contract, provider, updated_at DESC, job_id DESC);
  END IF;
END $$;

DO $$
BEGIN
  -- commerce_actions: agent timeline + sortable lists + job detail timeline
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_ca_job_timestamp_log') THEN
    CREATE INDEX idx_ca_job_timestamp_log
      ON commerce_actions(chain_id, commerce_contract, job_id, block_timestamp DESC, log_index DESC);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_ca_agent_chain_contract_timestamp') THEN
    CREATE INDEX idx_ca_agent_chain_contract_timestamp
      ON commerce_actions(agent_uid, chain_id, commerce_contract, block_timestamp DESC);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_ca_agent_budget_timestamp') THEN
    CREATE INDEX idx_ca_agent_budget_timestamp
      ON commerce_actions(agent_uid, job_budget DESC, block_timestamp DESC);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_ca_agent_signal_weight_timestamp') THEN
    CREATE INDEX idx_ca_agent_signal_weight_timestamp
      ON commerce_actions(agent_uid, signal_weight DESC, block_timestamp DESC);
  END IF;
END $$;
SQL

echo "[3/3] Done."


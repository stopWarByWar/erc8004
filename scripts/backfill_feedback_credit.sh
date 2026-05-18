#!/usr/bin/env bash
set -euo pipefail

# Backfill feedback_tag_scores_v2 from existing feedbacks rows.
#
# Why: the trigger feedback_credit_insert_trigger only fires on NEW INSERTs.
# If the migrations/202605120000_feedback_credit_score.psql was applied AFTER
# feedbacks already had data, feedback_tag_scores_v2 will be empty for those
# legacy rows. This script catches it up by calling the idempotent helper
# fc_refresh_tag_stats(agent_uid, tag) for every distinct (agent_uid, tag1)
# combination present in feedbacks.
#
# Safe properties:
#   - Read-only on `feedbacks` (never mutates feedback rows)
#   - Writes only to feedback_tag_scores_v2 (INSERT … ON CONFLICT DO UPDATE)
#   - Fully idempotent — re-running is a no-op if data hasn't changed
#   - feedback_credit_scores / feedback_reviewer_weights NOT touched:
#       feedback_credit_scores  → fills itself via the Go cron (5min ticker)
#       feedback_reviewer_weights → fills on-demand from commerce_scores_global
#
# Usage:
#   DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=disable" \
#     ./scripts/backfill_feedback_credit.sh
#
# Or with libpq env vars:
#   PGHOST=... PGPORT=... PGUSER=... PGPASSWORD=... PGDATABASE=... \
#     ./scripts/backfill_feedback_credit.sh
#
# Or with --dry-run to see the row count without writing:
#   ./scripts/backfill_feedback_credit.sh --dry-run

DRY_RUN=0
for arg in "$@"; do
  case "$arg" in
    -n|--dry-run) DRY_RUN=1 ;;
    -h|--help)
      sed -n '3,30p' "$0"
      exit 0 ;;
    *) echo "unknown arg: $arg" >&2; exit 2 ;;
  esac
done

PSQL_ARGS=(-v ON_ERROR_STOP=1 -q)
if [[ -n "${DATABASE_URL:-}" ]]; then
  PSQL_CONN=("$DATABASE_URL")
else
  PSQL_CONN=()
fi

run_psql() {
  psql "${PSQL_ARGS[@]}" "${PSQL_CONN[@]}" "$@"
}

# ─── 0. preflight ──────────────────────────────────────────────────────────
echo "▸ Preflight: checking required objects exist…"
PREFLIGHT="$(run_psql -tAc "
  SELECT
    (SELECT to_regclass('public.feedbacks')                IS NOT NULL) AS has_feedbacks,
    (SELECT to_regclass('public.feedback_tag_scores_v2')   IS NOT NULL) AS has_v2,
    (SELECT EXISTS (SELECT 1 FROM pg_proc
                    WHERE proname = 'fc_refresh_tag_stats')) AS has_func;
" | tr -d '[:space:]')"

case "$PREFLIGHT" in
  t\|t\|t) ;;  # all three exist
  *)
    echo "  ✗ Missing required objects (feedbacks / feedback_tag_scores_v2 / fc_refresh_tag_stats)" >&2
    echo "    Apply migrations/202605120000_feedback_credit_score.psql first." >&2
    exit 1 ;;
esac
echo "  ✓ feedbacks, feedback_tag_scores_v2, fc_refresh_tag_stats — all present"

# ─── 1. counts before ──────────────────────────────────────────────────────
echo
echo "▸ Before:"
run_psql <<'SQL'
SELECT
  (SELECT COUNT(*) FROM feedbacks
     WHERE tag1 IS NOT NULL AND tag1 <> '')                AS feedbacks_with_tag1,
  (SELECT COUNT(DISTINCT (agent_uid, tag1)) FROM feedbacks
     WHERE tag1 IS NOT NULL AND tag1 <> '')                AS distinct_agent_tag_pairs,
  (SELECT COUNT(*) FROM feedback_tag_scores_v2)            AS v2_rows_now;
SQL

if [[ $DRY_RUN -eq 1 ]]; then
  echo
  echo "▸ Dry-run: skipping backfill. Re-run without --dry-run to apply."
  exit 0
fi

# ─── 2. backfill ───────────────────────────────────────────────────────────
echo
echo "▸ Running fc_refresh_tag_stats() for every distinct (agent_uid, tag1)…"
START_TS="$(date +%s)"

# Wrap in a transaction so a failure mid-batch doesn't leave a partial state.
# fc_refresh_tag_stats is itself transactional per call; this BEGIN/COMMIT
# is purely an "all or nothing" wrapper for the loop.
run_psql <<'SQL'
BEGIN;
SELECT fc_refresh_tag_stats(agent_uid, tag1)
FROM (
  SELECT DISTINCT agent_uid, tag1
  FROM feedbacks
  WHERE tag1 IS NOT NULL AND tag1 <> ''
) AS t;
COMMIT;
SQL

ELAPSED=$(( $(date +%s) - START_TS ))
echo "  ✓ done in ${ELAPSED}s"

# ─── 3. counts after + sanity ─────────────────────────────────────────────
echo
echo "▸ After:"
run_psql <<'SQL'
SELECT
  (SELECT COUNT(*) FROM feedback_tag_scores_v2)                     AS v2_rows,
  (SELECT COUNT(DISTINCT agent_uid) FROM feedback_tag_scores_v2)    AS distinct_agents,
  (SELECT SUM(feedback_count) FROM feedback_tag_scores_v2)          AS total_fb_in_v2,
  (SELECT SUM(active_count)   FROM feedback_tag_scores_v2)          AS active_fb_in_v2,
  (SELECT SUM(revoked_count)  FROM feedback_tag_scores_v2)          AS revoked_in_v2;
SQL

echo
echo "▸ Sample (top 10 by active_count):"
run_psql <<'SQL'
SELECT agent_uid, tag, active_count, revoked_count, unique_reviewer_count,
       value_min, value_p50, value_p95, value_max,
       first_active_ts, last_active_ts
FROM feedback_tag_scores_v2
ORDER BY active_count DESC
LIMIT 10;
SQL

# ─── 4. cross-check ───────────────────────────────────────────────────────
echo
echo "▸ Cross-check: does SUM(feedback_count) in v2 match feedbacks-with-tag1?"
DIFF="$(run_psql -tAc "
  SELECT
    (SELECT COALESCE(SUM(feedback_count),0) FROM feedback_tag_scores_v2) -
    (SELECT COUNT(*) FROM feedbacks WHERE tag1 IS NOT NULL AND tag1 <> '');
" | tr -d '[:space:]')"

if [[ "$DIFF" == "0" ]]; then
  echo "  ✓ counts match exactly"
else
  echo "  ⚠ delta = $DIFF row(s). Investigate before relying on the cache."
fi

echo
echo "▸ Done."
echo "  feedback_credit_scores (per-agent rollup cache) will be filled by the Go cron"
echo "  on its next tick (default 5min). Or hit /agent/feedback/credit?uid=N to"
echo "  recompute on-demand."

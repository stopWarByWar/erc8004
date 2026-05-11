#!/bin/bash
# Re-run all PSQL migrations in chronological order
# Usage: ./rerun_all.sh [psql_conn_string]
#   e.g.: ./rerun_all.sh "host=localhost dbname=erc8004 user=postgres password=secret"
#   If no arg, reads from DATABASE_URL env var or uses PGPASSWORD+PG* vars.

set -e

CONN="${1:-${DATABASE_URL:-}}"
if [[ -z "$CONN" ]]; then
  # Fallback to individual PG* env vars
  CONN="host=${PGHOST:-localhost} port=${PGPORT:-5432} dbname=${PGDATABASE:-erc8004} user=${PGUSER:-postgres} password=${PGPASSWORD:-}"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS=$(ls -1 "$SCRIPT_DIR"/*.psql | sort)

echo "=== Re-running all PSQL migrations ==="
echo "Connection: ${CONN%% @*}"  # redact password in output
echo ""

for f in $MIGRATIONS; do
  name=$(basename "$f")
  echo "--- Executing: $name ---"
  psql "$CONN" -f "$f"
  echo "--- Done: $name ---"
  echo ""
done

echo "=== All migrations complete ==="

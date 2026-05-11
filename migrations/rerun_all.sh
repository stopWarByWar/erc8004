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
MIGRATIONS=$(ls -1 "$SCRIPT_DIR"/{*.psql,*.sql} 2>/dev/null | sort)

echo "=== PSQL Migration Runner ==="
echo "Connection: ${CONN%% @*}"
echo ""

if [[ ! -t 0 ]]; then
  echo "Error: not a terminal, cannot ask for input. Run with '-y' to skip all prompts."
  echo "Usage: $0 [-y] [conn_string]"
  exit 1
fi

echo -n "Run ALL migrations now? [y/N] "
read -r reply
echo ""
if [[ ! "$reply" =~ ^[Yy]$ ]]; then
  echo "Aborted."
  exit 0
fi

for f in $MIGRATIONS; do
  name=$(basename "$f")
  while true; do
    echo -n "--- Execute: $name? [y/s(=skip file)/q(=quit)] "
    read -r reply
    echo ""
    case "$reply" in
      [Yy] )
        psql "$CONN" -f "$f" && echo "--- Done: $name ---" && echo "" && break
        ;;
      [Ss] )
        echo "--- Skipped: $name ---" && echo "" && break
        ;;
      [Qq] )
        echo "Quit."
        exit 0
        ;;
      * )
        echo "Please answer y, s, or q."
        ;;
    esac
  done
done

echo "=== All done ==="
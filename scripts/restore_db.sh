#!/usr/bin/env bash
set -euo pipefail
: "${DB_HOST?}"; : "${DB_PORT?}"; : "${DB_USER?}"; : "${DB_PASS?}"; : "${DB_NAME?}"
export PGPASSWORD="$DB_PASS"
if [ -z "${1-}" ]; then echo "Usage: $0 path/to/backup.dump"; exit 1; fi
pg_restore -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" --clean "$1"
echo "Restore completed."

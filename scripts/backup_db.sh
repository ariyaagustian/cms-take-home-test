#!/usr/bin/env bash
set -euo pipefail
: "${DB_HOST?}"; : "${DB_PORT?}"; : "${DB_USER?}"; : "${DB_PASS?}"; : "${DB_NAME?}"
export PGPASSWORD="$DB_PASS"
pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -Fc > backup_$(date +%Y%m%d_%H%M%S).dump
echo "Backup created."

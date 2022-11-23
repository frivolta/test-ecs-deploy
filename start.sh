#!/bin/sh
set -e

export GIN_MODE=release
source /app/app.env
echo "run db migration $DB_SOURCE"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
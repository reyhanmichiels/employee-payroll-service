#!/bin/bash

set -euxo pipefail

HOST="${PSQL_HOST:-127.0.0.1}"
PORT="${PSQL_PORT:-5432}"
DB_NAME="${PSQL_DB:-employee_payroll_db}"
USERNAME="${PSQL_USER:-root}"
PASSWORD="${PSQL_PASSWORD:-password}"
PROTOCOL="${PSQL_PROTOCOL:-tcp}"

rm temp.sql || true
echo "DROP DATABASE IF EXISTS ${DB_NAME};" >> temp.sql
echo "CREATE DATABASE ${DB_NAME};" >> temp.sql
echo "\\c ${DB_NAME}" >> temp.sql
echo "SET session_replication_role = 'replica';" >> temp.sql
cat **/[!temp]*.sql >> temp.sql || true
echo "SET session_replication_role = 'origin';" >> temp.sql

PGPASSWORD=${PASSWORD} psql -h ${HOST} -p ${PORT} -U ${USERNAME} -d ${DB_NAME} -f temp.sql
rm temp.sql

#!/bin/sh

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<EOF
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
EOF

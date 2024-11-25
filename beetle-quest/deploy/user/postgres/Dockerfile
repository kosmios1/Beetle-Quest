#FROM postgres:17.0-alpine3.20
FROM cybertecpostgresql/postgresql-ee-demo:16

ENV POSTGRES_DB ''
ENV POSTGRES_USER ''
ENV POSTGRES_PASSWORD ''

COPY 1.load-extensions.sh /docker-entrypoint-initdb.d/
COPY user-table.sql /docker-entrypoint-initdb.d/

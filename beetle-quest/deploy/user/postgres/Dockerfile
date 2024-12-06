FROM cybertecpostgresql/postgresql-ee-demo:17.2ee1.4.0-1

ENV POSTGRES_DB ''
ENV POSTGRES_USER ''
ENV POSTGRES_PASSWORD ''

COPY 1.load-extensions.sh /docker-entrypoint-initdb.d/
COPY user-table.sql /docker-entrypoint-initdb.d/

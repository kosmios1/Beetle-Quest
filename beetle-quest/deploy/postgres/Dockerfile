FROM postgres:17.0-alpine3.20

ENV POSTGRES_DB ''
ENV POSTGRES_USER ''
ENV POSTGRES_PASSWORD ''

COPY load-extensions.sh /docker-entrypoint-initdb.d/
COPY load-gacha.sql     /docker-entrypoint-initdb.d/

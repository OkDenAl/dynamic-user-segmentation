FROM postgres:latest

RUN apt-get update && apt-get -y install postgresql-15-cron
COPY init-db /docker-entrypoint-initdb.d

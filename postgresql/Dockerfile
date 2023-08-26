FROM postgres:latest

RUN apt-get update && apt-get -y install postgresql-15-cron
RUN echo "shared_preload_libraries = 'pg_cron'" >> /var/lib/postgresql/data/postgresql.conf
RUN echo "cron.database_name = 'segmentation'" >> /var/lib/postgresql/data/postgresql.conf
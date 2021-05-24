#!/usr/bin/dumb-init /bin/bash
set -eax

POSTGRES_PORT=${POSTGRES_PORT:-5432}
mkdir -p /var/pg_dlogs
touch /var/pg_dlogs/postgresql.log
chmod -R 777 /var/pg_dlogs
pg-entrypoint.sh "postgres -p $POSTGRES_PORT -c logging_collector=on -c log_directory=/var/pg_dlogs -c log_filename=postgresql.log -c log_statement=all" &

while ! pg_isready -U$POSTGRES_USER -hlocalhost -p$POSTGRES_PORT > /dev/null 2> /dev/null; do
    echo "postgres not ready yet"
    sleep 3
done

horizon-entrypoint.sh $@

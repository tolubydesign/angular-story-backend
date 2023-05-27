#!/bin/bash
# set -e
echo "Initalising database bash"
# psql -h localhost -p 5432 -U postgres -c  "$(cat $PWD/init_database/init-table.sql)"
psql -h localhost -p 5432 -U postgres -c "$(cat $PWD/create_tables.sql)"
psql -h localhost -p 5432 -U postgres -c  "$(cat $PWD/fill_tables.sql)"
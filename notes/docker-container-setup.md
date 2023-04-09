# Handling docker compose and docker builds

This document describes my finding on how to spin up a postgres docker compose container. But this can be applied to a NoSQL or SQLite database.

### Spin Up a Container 
Commands:

`$ psql -h localhost -p 5432 -U postgres`

`$ psql -h localhost -p 5432 -U postgres -c "$(cat $PWD/init_database/create_tables.sql)"`

`$ psql -h localhost -p 5432 -U postgres -c  "$(cat $PWD/init_database/fill_tables.sql)"`

`$ psql -h localhost -p 5432 -U postgres -c  "$(cat $PWD/init_database/init-table.sql)"`


#### Relevant Resources

_Docker Installation:_

[docker arch linux](https://linuxhint.com/docker_arch_linux/)\
[Setting up docker with sudo](https://www.baeldung.com/linux/docker-run-without-sudo)\
[StackOverflow: Setting up user groups with docker](https://stackoverflow.com/questions/56305613/cant-add-user-to-docker-group/66297855)

_Backup and Restoring Database_

[StackOverflow: Backup and Restore Dockerized PostgreSQL Database](https://stackoverflow.com/questions/24718706/backup-restore-a-dockerized-postgresql-database)\
[How to Dump and Restore a postgresql Database from a Docker Container](https://davejansen.com/how-to-dump-and-restore-a-postgresql-database-from-a-docker-container/)\
```bash
$ docker exec -t your-db-container pg_dumpall -c -U postgres > dump_`date +%d-%m-%Y"_"%H_%M_%S`.sql
```
```bash
$ cat your_dump.sql | docker exec -i your-db-container psql -U postgres
```

_Get into database_

[StackOverflow: Connecting to PostgreSQL Docker Container](https://stackoverflow.com/questions/37694987/connecting-to-postgresql-in-a-docker-container-from-outside)

```bash
$ docker exec -it {container_id} bash
```
```bash
$ root@{container_id}:/# psql -U postgres
OR
$ sudo -iu postgres 
$ su -> $ su -l postgres
```
```bash
$ psql -U postgres OR $ psql -h localhost -p 5432 -U postgres
```

_Other_
[GitHub: Docker Cookbook](https://github.com/robcowart/docker_compose_cookbook)

[Postgres with Docker](https://geshan.com.np/blog/2021/12/docker-postgres/)


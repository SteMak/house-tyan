#!/bin/bash

docker rm -fv anihouse-postgres
docker volume create anihouse-pgdata
docker create --name anihouse-postgres -e POSTGRES_USER=bot -e POSTGRES_PASSWORD=password -e POSTGRES_DB=anihouse -p 54320:5432 postgres:11.5
docker start anihouse-postgres
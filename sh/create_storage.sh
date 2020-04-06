#!/bin/bash

docker rm -fv anihouse-postgres
# docker volume create anihouse-pgdata
docker create --name anihouse-postgres -e POSTGRES_PASSWORD=password -p 54320:5432 postgres:11.5
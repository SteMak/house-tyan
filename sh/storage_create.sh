#!/bin/bash

docker rm -fv pgdata-storage
docker volume create pgdata
docker create --name pgdata-storage -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=monitor -p 54323:5432 -v pgdata:/var/lib/postgresql/data postgres:11.5
docker start pgdata-storage
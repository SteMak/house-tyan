#!/bin/bash

docker rm -fv pgdata-storage
docker create --name pgdata-storage -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=bot -p 54320:5432 postgres:11.5
docker start pgdata-storage
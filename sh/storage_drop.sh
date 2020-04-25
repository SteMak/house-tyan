#!/bin/bash

DIR=../cli/bot/migrations/postgres
DATABASE=postgresql://bot:password@localhost:54320/anihouse?sslmode=disable

migrate -database $DATABASE -path $DIR drop
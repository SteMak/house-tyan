#!/bin/bash

DIR=../cli/bot/migrations/postgres
DATABASE=postgresql://admin:admin@localhost:54320/bot?sslmode=disable

migrate -database $DATABASE -path $DIR up
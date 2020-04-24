#!/bin/bash

migrate -path ../cli/bot/migrations/postgres -database postgres://bot:password@localhost:54320/anihouse?sslmode=disable down
migrate -path ../cli/bot/migrations/postgres -database postgres://bot:password@localhost:54320/anihouse?sslmode=disable up
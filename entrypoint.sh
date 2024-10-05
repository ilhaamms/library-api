#!/bin/sh

migrate -database "sqlite3://db/library.db" -path db/migrations up

./main

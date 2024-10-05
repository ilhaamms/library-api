#!/bin/sh

migrate -database "sqlite://db/library.db" -path db/migrations up

./main

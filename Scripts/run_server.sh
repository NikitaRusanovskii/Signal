#!/bin/bash

source .env

# "$GOLANG_MIGRATE" -path "$MIGRATIONS_DIR" -database "$DATABASE_URL?sslmode=disable" down
"$GOLANG_MIGRATE" -path "$MIGRATIONS_DIR" -database "$DATABASE_URL?sslmode=disable" up

cd ../
sudo docker compose up


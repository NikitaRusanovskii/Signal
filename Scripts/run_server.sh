#!/bin/bash

source .env


cd ../
sudo docker compose up

# "$GOLANG_MIGRATE" -path "$MIGRATIONS_DIR" -database "$DATABASE_URL?sslmode=disable" down
"$GOLANG_MIGRATE" -path "$MIGRATIONS_DIR" -database "$DATABASE_URL?sslmode=disable" up


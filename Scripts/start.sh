#!/bin/bash

source ../.env

cd $ROOT

sudo docker compose up -d;
sleep 2
gnome-terminal -- bash -c "\
    $GOLANG_MIGRATE -path ./migrations -database \"$DATABASE_URL?sslmode=disable\" up;\
    go run ./cmd/app;\
    $GOLANG_MIGRATE -path db/migrations -database \"$DATABASE_URL?sslmode=disable\" down; bash exec\
"
sleep 2
gnome-terminal -- bash -c ".$Root/Scripts/test_api.sh; bash exec"
read
sudo docker compose down;
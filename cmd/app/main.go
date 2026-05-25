package main

import (
	"context"
	"fmt"
	"os"
	"signal/internal/api"
	"signal/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	var dbURL string
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("godotenv.Loat() error: ", err)
	}
	dbURL = os.Getenv("DATABASE_URL")

	ctx := context.Background()
	cm := repository.NewConnectionManager(nil)
	cm.Connect(ctx, dbURL)
	defer cm.Disconnect()

	pool, err := cm.GetPool()
	if err != nil {
		panic(err)
	}
	pr, err := repository.NewPeerManager(pool)
	if err != nil {
		panic(err)
	}

	srv := api.InitServer(pr)
	srv.RunInactivePeerKiller(ctx, 30)
	srv.Run()
}

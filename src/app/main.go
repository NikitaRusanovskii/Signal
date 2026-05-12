package main

import (
	"context"
	"fmt"
	"os"
	"signal/internal"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("godotenv.Load() error")
	}
	dbURL := os.Getenv("DATABASE_URL")

	ctx := context.Background()

	cm := internal.NewConnectionManager(nil)
	cm.Connect(ctx, dbURL)
	defer cm.Disconnect()

	pool, err := cm.GetPool()
	if err != nil {
		panic("cm.GetPool() error")
	}
	pm, err := internal.NewPeerManager(pool)
	if err != nil {
		panic("NewPeerManager() error")
	}

	test_peer := internal.NewPeer(uuid.New(), internal.MasterRole,
		true, "192.0.0.51", time.Now())

	err_insert := pm.Insert(ctx, test_peer)
	if err_insert != nil {
		fmt.Println(err_insert)
		panic("insertion in db error")
	}

	peer, err_getByID := pm.GetByID(ctx, test_peer.ID)
	if err_getByID != nil {
		panic("getByID from db error")
	}
	fmt.Println("addr:port: ", peer.AddrPort)

	isExists, err_ExistsByID := pm.ExistsByID(ctx, test_peer.ID)
	if err_ExistsByID != nil {
		panic("existsByID from db error")
	}
	fmt.Println("test_peer is exists: ", isExists)

	err_deleteByID := pm.DeleteByID(ctx, test_peer.ID)
	if err_deleteByID != nil {
		panic("deleteByID from db error")
	}

	isExists, err_ExistsByID = pm.ExistsByID(ctx, test_peer.ID)
	if err_ExistsByID != nil {
		panic("existsByID from db error")
	}
	fmt.Println("test_peer is exists: ", isExists)

}

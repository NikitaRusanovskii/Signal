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
		panic(err)
	}
	pr, err := internal.NewPeerManager(pool)
	if err != nil {
		panic(err)
	}
	ps := internal.NewPeerService(pr)

	test_master_peer := internal.NewPeer(
		uuid.New(),
		internal.MasterRole,
		true,
		"192.0.0.1:1458",
		time.Now(),
	)

	test_slave_peer := internal.NewPeer(
		uuid.New(),
		internal.SlaveRole,
		true,
		"102.1.23.1:1321",
		time.Now(),
	)

	ps.RegisterPeer(ctx, test_master_peer)
	ps.RegisterPeer(ctx, test_slave_peer)

	ap_master, err := ps.GetMasterIP(ctx)
	if err != nil {
		panic(err)
	}

	ap_slave, err := ps.GetLastPeerIP(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Master IP:port : ", ap_master.String())
	fmt.Println("Slave IP:port : ", ap_slave.String())

	ps.UpdatePeerAddrPort(ctx, test_slave_peer.ID, "102.1.23.1:1322")
	ap_slave2, err := ps.GetLastPeerIP(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Slave IP:port : ", ap_slave2.String())

	ps.RemovePeer(ctx, test_slave_peer.ID)

	ap_slave3, err := ps.GetLastPeerIP(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Slave IP:port : ", ap_slave3.String())
}

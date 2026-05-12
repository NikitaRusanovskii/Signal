package internal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionManager struct {
	db *pgxpool.Pool
}

func NewConnectionManager(db *pgxpool.Pool) *ConnectionManager {
	if db == nil {
		return &ConnectionManager{db: nil}
	}
	return &ConnectionManager{db: db}
}

func (c *ConnectionManager) Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (c *ConnectionManager) Disconnect() {
	c.db.Close()
}

func (c *ConnectionManager) GetPool() (*pgxpool.Pool, error) {
	if c.db == nil {
		return nil, errors.New("db pgxpool.Pool is nil. ConnectionManager:GetPool()")
	}
	return c.db, nil
}

type PeerManager struct {
	db *pgxpool.Pool
}

func NewPeerManager(db *pgxpool.Pool) (*PeerManager, error) {
	if db == nil {
		return nil, errors.New("db pgxpool.Pool is nil. NewPeerManager()")
	}
	return &PeerManager{db: db}, nil
}

func (p *PeerManager) Insert(ctx context.Context, peer Peer) error {
	query := `
		INSERT INTO peers (id, role, is_online, addr_port, connection_time)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id)
			DO UPDATE SET
			role = EXCLUDED.role,
			is_online = EXCLUDED.is_online,
			addr_port = EXCLUDED.addr_port,
			connection_time = EXCLUDED.connection_time
	`

	_, err := p.db.Exec(ctx, query,
		peer.ID, peer.Role, peer.IsOnline,
		peer.AddrPort, peer.ConnectionTime)

	return err
}

func (p *PeerManager) GetByID(ctx context.Context, id uuid.UUID) (*Peer, error) {
	query := `
	SELECT id, role, is_online, addr_port, connection_time
	FROM peers
	WHERE id = $1
	`

	res := p.db.QueryRow(ctx, query, id)
	peer := &Peer{}
	err := res.Scan(&peer.ID, &peer.Role,
		&peer.IsOnline, &peer.AddrPort, &peer.ConnectionTime)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("Unknown peer id")
	}
	return peer, nil
}

func (p *PeerManager) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM peers
	WHERE id = $1
	`
	_, err := p.db.Exec(ctx, query, id)
	return err
}

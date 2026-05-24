package repository

import (
	"context"
	"errors"
	"net/netip"
	"signal/internal/domain"

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
	c.db = pool
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

func (p *PeerManager) Insert(ctx context.Context, peer domain.Peer) error {
	query := `
		INSERT INTO peers (role, is_online, addr_port, connection_time)
		SELECT ($1, $2, $3, $4)
		WHERE NOT EXISTS (
    		SELECT 1 FROM peers WHERE addr_port = $3
		)
	`

	_, err := p.db.Exec(ctx, query, peer.Role, peer.IsOnline,
		peer.AddrPort, peer.LastSeen)

	return err
}

func (p *PeerManager) Save(ctx context.Context, peer domain.Peer) error {
	query := `
		INSERT INTO peers (role, is_online, addr_port, connection_time)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (addr_port)
			DO UPDATE SET
			role = EXCLUDED.role,
			is_online = EXCLUDED.is_online,
			connection_time = EXCLUDED.connection_time
	`

	_, err := p.db.Exec(ctx, query, peer.Role, peer.IsOnline,
		peer.AddrPort, peer.LastSeen)

	return err
}

func (p *PeerManager) Delete(ctx context.Context, addrPort netip.AddrPort) error {
	query := `
	DELETE FROM peers
	WHERE addr_port = $1
	`
	_, err := p.db.Exec(ctx, query, addrPort.String())
	return err
}

func (p *PeerManager) GetByRole(ctx context.Context, role domain.Role) ([]*domain.Peer, error) {
	query := `
	SELECT role, is_online, addr_port, connection_time
	FROM peers
	WHERE role = $1
	`

	rows, err := p.db.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var peers []*domain.Peer

	for rows.Next() {

		peer := &domain.Peer{}
		err := rows.Scan(&peer.Role, &peer.IsOnline, &peer.AddrPort, &peer.LastSeen)
		if err != nil {
			return nil, err
		}

		peers = append(peers, peer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return peers, nil
}

func (p *PeerManager) GetByAddrPort(ctx context.Context, addrPort netip.AddrPort) (*domain.Peer, error) {
	query := `
	SELECT role, is_online, addr_port, connection_time
	FROM peers
	WHERE addr_port = $1
	`

	res := p.db.QueryRow(ctx, query, addrPort.String())
	peer := &domain.Peer{}
	err := res.Scan(&peer.Role,
		&peer.IsOnline, &peer.AddrPort, &peer.LastSeen)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("unknown peer role")
	}
	return peer, nil
}

func (p *PeerManager) SetOffline(ctx context.Context, periodInSeconds uint) error {
	query := `
	UPDATE peers
	SET is_online = FALSE
	WHERE is_online = TRUE AND last_seen < NOW() - ($1 * INTERVAL '1 second')
	`

	_, err := p.db.Exec(ctx, query, int(periodInSeconds))
	return err
}

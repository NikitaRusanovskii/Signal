package internal

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	MasterRole Role = "master"
	SlaveRole  Role = "slave"
)

type Peer struct {
	ID             uuid.UUID `db:"id"`
	Role           Role      `db:"role"`
	IsOnline       bool      `db:"is_online"`
	AddrPort       string    `db:"addr_port"`
	ConnectionTime time.Time `db:"connection_time"`
}

func NewPeer(ID uuid.UUID, role Role, isOnline bool,
	addrPort string, connectionTime time.Time) Peer {
	return Peer{
		ID:             ID,
		Role:           role,
		IsOnline:       isOnline,
		AddrPort:       addrPort,
		ConnectionTime: connectionTime,
	}
}

func (p *Peer) SetRole(newRole Role) {
	p.Role = newRole
}

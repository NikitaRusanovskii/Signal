package domain

import (
	"time"
)

type Role string

const (
	MasterRole Role = "master"
	SlaveRole  Role = "slave"
)

type Peer struct {
	Role     Role      `db:"role"`
	IsOnline bool      `db:"is_online"`
	AddrPort string    `db:"addr_port"`
	LastSeen time.Time `db:"last_seen"`
}

func NewPeer(role Role, isOnline bool,
	addrPort string, lastSeen time.Time) Peer {
	return Peer{
		Role:     role,
		IsOnline: isOnline,
		AddrPort: addrPort,
		LastSeen: lastSeen,
	}
}

func (p *Peer) SetRole(newRole Role) {
	p.Role = newRole
}

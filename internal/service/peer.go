package service

import (
	"context"
	"errors"
	"net/netip"
	"signal/internal/domain"
	"signal/internal/repository"

	"github.com/google/uuid"
)

type PeerService struct {
	pr *repository.PeerManager
}

func NewPeerService(pr *repository.PeerManager) *PeerService {
	return &PeerService{pr: pr}
}

func (ps *PeerService) RegisterPeer(ctx context.Context, p domain.Peer) error {
	isExists, err := ps.pr.ExistsByID(ctx, p.ID)
	if err != nil {
		return err
	}
	if isExists {
		return errors.New("Peer already registered")
	}
	err = ps.pr.Insert(ctx, p)
	return err
}

func (ps *PeerService) UpdatePeerRole(ctx context.Context, id uuid.UUID, role domain.Role) error {
	isExists, err := ps.pr.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !isExists {
		return errors.New("Peer does not exists")
	}
	peer, err := ps.pr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	peer.Role = role
	err = ps.pr.Insert(ctx, *peer)
	return err
}

func (ps *PeerService) UpdatePeerOnlineStatus(ctx context.Context, id uuid.UUID, onlineStatus bool) error {
	isExists, err := ps.pr.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !isExists {
		return errors.New("Peer does not exists")
	}
	peer, err := ps.pr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	peer.IsOnline = onlineStatus
	err = ps.pr.Insert(ctx, *peer)
	return err
}

func (ps *PeerService) UpdatePeerAddrPort(ctx context.Context, id uuid.UUID, AddrPort string) error {
	isExists, err := ps.pr.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !isExists {
		return errors.New("Peer does not exists")
	}
	peer, err := ps.pr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	peer.AddrPort = AddrPort
	err = ps.pr.Insert(ctx, *peer)
	return err
}

func (ps *PeerService) RemovePeer(ctx context.Context, id uuid.UUID) error {
	err := ps.pr.DeleteByID(ctx, id)
	return err
}

func (ps *PeerService) GetLastSlavePeerIP(ctx context.Context) (*netip.AddrPort, error) {
	peer, err := ps.pr.GetLastSlaveByTime(ctx)
	if err != nil {
		return nil, err
	}
	AddrPort, err := netip.ParseAddrPort(peer.AddrPort)
	if err != nil {
		return nil, err
	}
	return &AddrPort, nil
}

func (ps *PeerService) GetMasterIP(ctx context.Context) (*netip.AddrPort, error) {
	peer, err := ps.pr.GetByRole(ctx, domain.MasterRole)
	if err != nil {
		return nil, err
	}
	AddrPort, err := netip.ParseAddrPort(peer.AddrPort)
	if err != nil {
		return nil, err
	}
	return &AddrPort, nil
}

func (ps *PeerService) UpdatePeerOnlineStatusByIP(ctx context.Context, onlineStatus bool, addrPort netip.AddrPort) error {
	peer, err := ps.pr.GetPeerByAddrPort(ctx, addrPort)
	if err != nil {
		return err
	}
	peer.IsOnline = onlineStatus
	err = ps.pr.Insert(ctx, *peer)
	return err
}

func (ps *PeerService) IsEmpty(ctx context.Context) bool {
	exists, _ := ps.pr.IsEmpty(ctx)
	return exists
}

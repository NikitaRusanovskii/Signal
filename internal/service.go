package internal

import (
	"context"
	"errors"
	"net/netip"

	"github.com/google/uuid"
)

type PeerService struct {
	pr *PeerManager
}

func NewPeerService(pr *PeerManager) *PeerService {
	return &PeerService{pr: pr}
}

func (ps *PeerService) RegisterPeer(ctx context.Context, p Peer) error {
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

func (ps *PeerService) UpdatePeerRole(ctx context.Context, id uuid.UUID, role Role) error {
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
	peer, err := ps.pr.GetByRole(ctx, MasterRole)
	if err != nil {
		return nil, err
	}
	AddrPort, err := netip.ParseAddrPort(peer.AddrPort)
	if err != nil {
		return nil, err
	}
	return &AddrPort, nil
}

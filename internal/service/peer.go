package service

import (
	"context"
	"net/netip"
	"signal/internal/domain"
	"signal/internal/repository"
	"time"
)

type PeerService struct {
	pr *repository.PeerManager
}

func NewPeerService(pr *repository.PeerManager) *PeerService {
	return &PeerService{pr: pr}
}

func (ps *PeerService) AddPeer(ctx context.Context, peer domain.Peer) error {
	err := ps.pr.Insert(ctx, peer)
	return err
}

func (ps *PeerService) DeletePeer(ctx context.Context, addrPort netip.AddrPort) error {
	err := ps.pr.Delete(ctx, addrPort)
	return err
}

func (ps *PeerService) SetRole(ctx context.Context, addrPort netip.AddrPort, role domain.Role) error {
	peer, err := ps.pr.GetByAddrPort(ctx, addrPort)
	if err != nil {
		return nil
	}
	peer.Role = role
	err = ps.pr.Save(ctx, *peer)
	return err
}

func (ps *PeerService) SetOnline(ctx context.Context, addrPort netip.AddrPort, isOnline bool) error {
	peer, err := ps.pr.GetByAddrPort(ctx, addrPort)
	if err != nil {
		return err
	}

	peer.IsOnline = isOnline
	err = ps.pr.Save(ctx, *peer)
	return err
}

func (ps *PeerService) GetMastersAddrPort(ctx context.Context) ([]netip.AddrPort, error) {
	peers, err := ps.pr.GetByRole(ctx, domain.MasterRole)
	if err != nil {
		return nil, err
	}
	var addrs []netip.AddrPort
	for _, peer := range peers {
		addrs = append(addrs, netip.MustParseAddrPort(peer.AddrPort))
	}
	return addrs, nil
}

func (ps *PeerService) GetSlavesAddrPort(ctx context.Context) ([]netip.AddrPort, error) {
	peers, err := ps.pr.GetByRole(ctx, domain.SlaveRole)
	if err != nil {
		return nil, err
	}
	var addrs []netip.AddrPort
	for _, peer := range peers {
		addrs = append(addrs, netip.MustParseAddrPort(peer.AddrPort))
	}
	return addrs, nil
}

func (ps *PeerService) Killer(ctx context.Context, periodInSeconds uint) {
	ticker := time.NewTicker(time.Duration(periodInSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := ps.pr.SetOffline(ctx, periodInSeconds)
			if err != nil {
				return
			}
		}
	}
}

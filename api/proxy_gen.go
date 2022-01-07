// Code generated by github.com/lyswifter/dbridge/gen/api. DO NOT EDIT.

package api

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"golang.org/x/xerrors"
)

var ErrNotSupported = xerrors.New("method not supported")

type CommonStruct struct {
	Internal struct {
		AuthNew func(p0 context.Context, p1 []auth.Permission) ([]byte, error) `perm:"admin"`

		AuthVerify func(p0 context.Context, p1 string) ([]auth.Permission, error) `perm:"read"`

		Closing func(p0 context.Context) (<-chan struct{}, error) `perm:"read"`

		Session func(p0 context.Context) (uuid.UUID, error) `perm:"read"`

		Shutdown func(p0 context.Context) error `perm:"admin"`
	}
}

type CommonStub struct {
}

type CommonNetStruct struct {
	CommonStruct

	NetStruct

	Internal struct {
	}
}

type CommonNetStub struct {
	CommonStub

	NetStub
}

type FullNodeStruct struct {
	CommonStruct

	NetStruct

	Internal struct {
	}
}

type FullNodeStub struct {
	CommonStub

	NetStub
}

type NetStruct struct {
	Internal struct {
		ID func(p0 context.Context) (peer.ID, error) `perm:"read"`

		NetAddrsListen func(p0 context.Context) (peer.AddrInfo, error) `perm:"read"`

		NetAgentVersion func(p0 context.Context, p1 peer.ID) (string, error) `perm:"read"`

		NetAutoNatStatus func(p0 context.Context) (NatInfo, error) `perm:"read"`

		NetBandwidthStats func(p0 context.Context) (metrics.Stats, error) `perm:"read"`

		NetBandwidthStatsByPeer func(p0 context.Context) (map[string]metrics.Stats, error) `perm:"read"`

		NetBandwidthStatsByProtocol func(p0 context.Context) (map[protocol.ID]metrics.Stats, error) `perm:"read"`

		NetBlockAdd func(p0 context.Context, p1 NetBlockList) error `perm:"admin"`

		NetBlockList func(p0 context.Context) (NetBlockList, error) `perm:"read"`

		NetBlockRemove func(p0 context.Context, p1 NetBlockList) error `perm:"admin"`

		NetConnect func(p0 context.Context, p1 peer.AddrInfo) error `perm:"write"`

		NetConnectedness func(p0 context.Context, p1 peer.ID) (network.Connectedness, error) `perm:"read"`

		NetDisconnect func(p0 context.Context, p1 peer.ID) error `perm:"write"`

		NetFindPeer func(p0 context.Context, p1 peer.ID) (peer.AddrInfo, error) `perm:"read"`

		NetPeerInfo func(p0 context.Context, p1 peer.ID) (*ExtendedPeerInfo, error) `perm:"read"`

		NetPeers func(p0 context.Context) ([]peer.AddrInfo, error) `perm:"read"`

		NetPubsubScores func(p0 context.Context) ([]PubsubScore, error) `perm:"read"`
	}
}

type NetStub struct {
}

func (s *CommonStruct) AuthNew(p0 context.Context, p1 []auth.Permission) ([]byte, error) {
	if s.Internal.AuthNew == nil {
		return *new([]byte), ErrNotSupported
	}
	return s.Internal.AuthNew(p0, p1)
}

func (s *CommonStub) AuthNew(p0 context.Context, p1 []auth.Permission) ([]byte, error) {
	return *new([]byte), ErrNotSupported
}

func (s *CommonStruct) AuthVerify(p0 context.Context, p1 string) ([]auth.Permission, error) {
	if s.Internal.AuthVerify == nil {
		return *new([]auth.Permission), ErrNotSupported
	}
	return s.Internal.AuthVerify(p0, p1)
}

func (s *CommonStub) AuthVerify(p0 context.Context, p1 string) ([]auth.Permission, error) {
	return *new([]auth.Permission), ErrNotSupported
}

func (s *CommonStruct) Closing(p0 context.Context) (<-chan struct{}, error) {
	if s.Internal.Closing == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.Closing(p0)
}

func (s *CommonStub) Closing(p0 context.Context) (<-chan struct{}, error) {
	return nil, ErrNotSupported
}

func (s *CommonStruct) Session(p0 context.Context) (uuid.UUID, error) {
	if s.Internal.Session == nil {
		return *new(uuid.UUID), ErrNotSupported
	}
	return s.Internal.Session(p0)
}

func (s *CommonStub) Session(p0 context.Context) (uuid.UUID, error) {
	return *new(uuid.UUID), ErrNotSupported
}

func (s *CommonStruct) Shutdown(p0 context.Context) error {
	if s.Internal.Shutdown == nil {
		return ErrNotSupported
	}
	return s.Internal.Shutdown(p0)
}

func (s *CommonStub) Shutdown(p0 context.Context) error {
	return ErrNotSupported
}

func (s *NetStruct) ID(p0 context.Context) (peer.ID, error) {
	if s.Internal.ID == nil {
		return *new(peer.ID), ErrNotSupported
	}
	return s.Internal.ID(p0)
}

func (s *NetStub) ID(p0 context.Context) (peer.ID, error) {
	return *new(peer.ID), ErrNotSupported
}

func (s *NetStruct) NetAddrsListen(p0 context.Context) (peer.AddrInfo, error) {
	if s.Internal.NetAddrsListen == nil {
		return *new(peer.AddrInfo), ErrNotSupported
	}
	return s.Internal.NetAddrsListen(p0)
}

func (s *NetStub) NetAddrsListen(p0 context.Context) (peer.AddrInfo, error) {
	return *new(peer.AddrInfo), ErrNotSupported
}

func (s *NetStruct) NetAgentVersion(p0 context.Context, p1 peer.ID) (string, error) {
	if s.Internal.NetAgentVersion == nil {
		return "", ErrNotSupported
	}
	return s.Internal.NetAgentVersion(p0, p1)
}

func (s *NetStub) NetAgentVersion(p0 context.Context, p1 peer.ID) (string, error) {
	return "", ErrNotSupported
}

func (s *NetStruct) NetAutoNatStatus(p0 context.Context) (NatInfo, error) {
	if s.Internal.NetAutoNatStatus == nil {
		return *new(NatInfo), ErrNotSupported
	}
	return s.Internal.NetAutoNatStatus(p0)
}

func (s *NetStub) NetAutoNatStatus(p0 context.Context) (NatInfo, error) {
	return *new(NatInfo), ErrNotSupported
}

func (s *NetStruct) NetBandwidthStats(p0 context.Context) (metrics.Stats, error) {
	if s.Internal.NetBandwidthStats == nil {
		return *new(metrics.Stats), ErrNotSupported
	}
	return s.Internal.NetBandwidthStats(p0)
}

func (s *NetStub) NetBandwidthStats(p0 context.Context) (metrics.Stats, error) {
	return *new(metrics.Stats), ErrNotSupported
}

func (s *NetStruct) NetBandwidthStatsByPeer(p0 context.Context) (map[string]metrics.Stats, error) {
	if s.Internal.NetBandwidthStatsByPeer == nil {
		return *new(map[string]metrics.Stats), ErrNotSupported
	}
	return s.Internal.NetBandwidthStatsByPeer(p0)
}

func (s *NetStub) NetBandwidthStatsByPeer(p0 context.Context) (map[string]metrics.Stats, error) {
	return *new(map[string]metrics.Stats), ErrNotSupported
}

func (s *NetStruct) NetBandwidthStatsByProtocol(p0 context.Context) (map[protocol.ID]metrics.Stats, error) {
	if s.Internal.NetBandwidthStatsByProtocol == nil {
		return *new(map[protocol.ID]metrics.Stats), ErrNotSupported
	}
	return s.Internal.NetBandwidthStatsByProtocol(p0)
}

func (s *NetStub) NetBandwidthStatsByProtocol(p0 context.Context) (map[protocol.ID]metrics.Stats, error) {
	return *new(map[protocol.ID]metrics.Stats), ErrNotSupported
}

func (s *NetStruct) NetBlockAdd(p0 context.Context, p1 NetBlockList) error {
	if s.Internal.NetBlockAdd == nil {
		return ErrNotSupported
	}
	return s.Internal.NetBlockAdd(p0, p1)
}

func (s *NetStub) NetBlockAdd(p0 context.Context, p1 NetBlockList) error {
	return ErrNotSupported
}

func (s *NetStruct) NetBlockList(p0 context.Context) (NetBlockList, error) {
	if s.Internal.NetBlockList == nil {
		return *new(NetBlockList), ErrNotSupported
	}
	return s.Internal.NetBlockList(p0)
}

func (s *NetStub) NetBlockList(p0 context.Context) (NetBlockList, error) {
	return *new(NetBlockList), ErrNotSupported
}

func (s *NetStruct) NetBlockRemove(p0 context.Context, p1 NetBlockList) error {
	if s.Internal.NetBlockRemove == nil {
		return ErrNotSupported
	}
	return s.Internal.NetBlockRemove(p0, p1)
}

func (s *NetStub) NetBlockRemove(p0 context.Context, p1 NetBlockList) error {
	return ErrNotSupported
}

func (s *NetStruct) NetConnect(p0 context.Context, p1 peer.AddrInfo) error {
	if s.Internal.NetConnect == nil {
		return ErrNotSupported
	}
	return s.Internal.NetConnect(p0, p1)
}

func (s *NetStub) NetConnect(p0 context.Context, p1 peer.AddrInfo) error {
	return ErrNotSupported
}

func (s *NetStruct) NetConnectedness(p0 context.Context, p1 peer.ID) (network.Connectedness, error) {
	if s.Internal.NetConnectedness == nil {
		return *new(network.Connectedness), ErrNotSupported
	}
	return s.Internal.NetConnectedness(p0, p1)
}

func (s *NetStub) NetConnectedness(p0 context.Context, p1 peer.ID) (network.Connectedness, error) {
	return *new(network.Connectedness), ErrNotSupported
}

func (s *NetStruct) NetDisconnect(p0 context.Context, p1 peer.ID) error {
	if s.Internal.NetDisconnect == nil {
		return ErrNotSupported
	}
	return s.Internal.NetDisconnect(p0, p1)
}

func (s *NetStub) NetDisconnect(p0 context.Context, p1 peer.ID) error {
	return ErrNotSupported
}

func (s *NetStruct) NetFindPeer(p0 context.Context, p1 peer.ID) (peer.AddrInfo, error) {
	if s.Internal.NetFindPeer == nil {
		return *new(peer.AddrInfo), ErrNotSupported
	}
	return s.Internal.NetFindPeer(p0, p1)
}

func (s *NetStub) NetFindPeer(p0 context.Context, p1 peer.ID) (peer.AddrInfo, error) {
	return *new(peer.AddrInfo), ErrNotSupported
}

func (s *NetStruct) NetPeerInfo(p0 context.Context, p1 peer.ID) (*ExtendedPeerInfo, error) {
	if s.Internal.NetPeerInfo == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.NetPeerInfo(p0, p1)
}

func (s *NetStub) NetPeerInfo(p0 context.Context, p1 peer.ID) (*ExtendedPeerInfo, error) {
	return nil, ErrNotSupported
}

func (s *NetStruct) NetPeers(p0 context.Context) ([]peer.AddrInfo, error) {
	if s.Internal.NetPeers == nil {
		return *new([]peer.AddrInfo), ErrNotSupported
	}
	return s.Internal.NetPeers(p0)
}

func (s *NetStub) NetPeers(p0 context.Context) ([]peer.AddrInfo, error) {
	return *new([]peer.AddrInfo), ErrNotSupported
}

func (s *NetStruct) NetPubsubScores(p0 context.Context) ([]PubsubScore, error) {
	if s.Internal.NetPubsubScores == nil {
		return *new([]PubsubScore), ErrNotSupported
	}
	return s.Internal.NetPubsubScores(p0)
}

func (s *NetStub) NetPubsubScores(p0 context.Context) ([]PubsubScore, error) {
	return *new([]PubsubScore), ErrNotSupported
}

var _ Common = new(CommonStruct)
var _ CommonNet = new(CommonNetStruct)
var _ FullNode = new(FullNodeStruct)
var _ Net = new(NetStruct)

package api

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

type Net interface {
	// ID returns peerID of libp2p node backing this API
	ID(context.Context) (peer.ID, error) //perm:read
}

type CommonNet interface {
	Common
	Net
}

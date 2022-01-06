package node

import (
	"time"

	"github.com/cskr/pubsub"
	"github.com/lyswifter/dbridge/api"
	"github.com/lyswifter/dbridge/node/config"
	"github.com/lyswifter/dbridge/node/impl/common"
	"github.com/lyswifter/dbridge/node/impl/net"
	"github.com/lyswifter/dbridge/node/modules"
	"github.com/lyswifter/dbridge/node/modules/dtypes"
	"github.com/lyswifter/dbridge/node/modules/lp2p"
	"github.com/lyswifter/dbridge/node/repo"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/xerrors"
)

var (
	AddrsFactoryKey = special{3} // Libp2p option
	NatPortMapKey   = special{8} // Libp2p option
)

func ConfigFullNode(c interface{}) Option {
	cfg, ok := c.(*config.BdridgeNode)
	if !ok {
		return Error(xerrors.Errorf("invalid config from repo, got: %T", c))
	}

	enableLibp2pNode := true // always enable libp2p for full nodes

	return Option(
		ConfigCommon(&cfg.Common, enableLibp2pNode),
	)
}

// Config sets up constructors based on the provided Config
func ConfigCommon(cfg *config.Common, enableLibp2pNode bool) Option {
	return Options(
		func(s *Settings) error { s.Config = true; return nil },
		Override(new(dtypes.APIEndpoint), func() (dtypes.APIEndpoint, error) {
			return multiaddr.NewMultiaddr(cfg.API.ListenAddress)
		}),
		Override(SetApiEndpointKey, func(lr repo.LockedRepo, e dtypes.APIEndpoint) error {
			return lr.SetAPIEndpoint(e)
		}),
		ApplyIf(func(s *Settings) bool { return s.Base }), // apply only if Base has already been applied
		If(!enableLibp2pNode,
			Override(new(api.Net), new(api.NetStub)),
			Override(new(api.Common), From(new(common.CommonAPI))),
		),
		If(enableLibp2pNode,
			Override(new(api.Net), From(new(net.NetAPI))),
			Override(new(api.Common), From(new(common.CommonAPI))),
			Override(StartListeningKey, lp2p.StartListening(cfg.Libp2p.ListenAddresses)),
			Override(ConnectionManagerKey, lp2p.ConnectionManager(
				cfg.Libp2p.ConnMgrLow,
				cfg.Libp2p.ConnMgrHigh,
				time.Duration(cfg.Libp2p.ConnMgrGrace),
				cfg.Libp2p.ProtectedPeers)),
			Override(new(*pubsub.PubSub), lp2p.GossipSub),
			Override(new(*config.Pubsub), &cfg.Pubsub),

			ApplyIf(func(s *Settings) bool { return len(cfg.Libp2p.BootstrapPeers) > 0 },
				Override(new(dtypes.BootstrapPeers), modules.ConfigBootstrap(cfg.Libp2p.BootstrapPeers)),
			),

			Override(AddrsFactoryKey, lp2p.AddrsFactory(
				cfg.Libp2p.AnnounceAddresses,
				cfg.Libp2p.NoAnnounceAddresses)),

			If(!cfg.Libp2p.DisableNatPortMap, Override(NatPortMapKey, lp2p.NatPortMap)),
		),
		Override(new(dtypes.MetadataDS), modules.Datastore(cfg.Backup.DisableMetadataLog)),
	)
}

package node

import (
	"context"
	"errors"
	"time"

	"github.com/cskr/pubsub"
	logging "github.com/ipfs/go-log/v2"
	metricsi "github.com/ipfs/go-metrics-interface"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/p2p/net/conngater"
	"github.com/lyswifter/dbridge/api"
	"github.com/lyswifter/dbridge/node/config"
	"github.com/lyswifter/dbridge/node/impl"
	"github.com/lyswifter/dbridge/node/modules"
	"github.com/lyswifter/dbridge/node/modules/dtypes"
	"github.com/lyswifter/dbridge/node/modules/helpers"
	"github.com/lyswifter/dbridge/node/modules/lp2p"
	"github.com/lyswifter/dbridge/node/repo"
	"github.com/lyswifter/dbridge/types"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

var log = logging.Logger("builder")

// special is a type used to give keys to modules which
//  can't really be identified by the returned type
type special struct{ id int }

type invoke int

// Invokes are called in the order they are defined.
//nolint:golint
const (
	// InitJournal at position 0 initializes the journal global var as soon as
	// the system starts, so that it's available for all other components.
	InitJournalKey = invoke(iota)

	StartListeningKey
	ConnectionManagerKey

	DefaultTransportsKey
	PstoreAddSelfKeysKey
	SmuxTransportKey
	RelayKey
	SecurityKey
	DiscoveryHandlerKey
	BaseRoutingKey
	BandwidthReporterKey
	AutoNATSvcKey

	ConnGaterKey
	RunPeerMgrKey

	ExtractApiKey

	SetApiEndpointKey

	_nInvokes // keep this last
)

type Settings struct {
	// modules is a map of constructors for DI
	//
	// In most cases the index will be a reflect. Type of element returned by
	// the constructor, but for some 'constructors' it's hard to specify what's
	// the return type should be (or the constructor returns fx group)
	modules map[interface{}]fx.Option

	// invokes are separate from modules as they can't be referenced by return
	// type, and must be applied in correct order
	invokes []fx.Option

	nodeType repo.RepoType

	Base   bool // Base option applied
	Config bool // Config option applied
	Lite   bool // Start node in "lite" mode

	enableLibp2pNode bool
}

func IsType(t repo.RepoType) func(s *Settings) bool {
	return func(s *Settings) bool { return s.nodeType == t }
}

// Basic lotus-app services
func defaults() []Option {
	return []Option{
		Override(new(helpers.MetricsCtx), func() context.Context {
			return metricsi.CtxScope(context.Background(), "lorry")
		}),

		Override(new(dtypes.ShutdownChan), make(chan struct{})),

		// // the great context in the sky, otherwise we can't DI build genesis; there has to be a better
		// // solution than this hack.
		Override(new(context.Context), func(lc fx.Lifecycle, mctx helpers.MetricsCtx) context.Context {
			return helpers.LifecycleCtx(mctx, lc)
		}),
	}
}

var LibP2P = Options(
	// Host config
	Override(new(dtypes.Bootstrapper), dtypes.Bootstrapper(false)),

	// Host dependencies
	Override(new(peerstore.Peerstore), lp2p.Peerstore),
	Override(PstoreAddSelfKeysKey, lp2p.PstoreAddSelfKeys),
	Override(StartListeningKey, lp2p.StartListening(config.DefaultDbridgeNode().Libp2p.ListenAddresses)),

	// Host settings
	Override(DefaultTransportsKey, lp2p.DefaultTransports),
	Override(AddrsFactoryKey, lp2p.AddrsFactory(nil, nil)),
	Override(SmuxTransportKey, lp2p.SmuxTransport()),
	Override(RelayKey, lp2p.NoRelay()),
	Override(SecurityKey, lp2p.Security(true, false)),

	// Host
	Override(new(lp2p.RawHost), lp2p.Host),
	Override(new(host.Host), lp2p.RoutedHost),
	Override(new(lp2p.BaseIpfsRouting), lp2p.DHTRouting(dht.ModeAuto)),

	Override(DiscoveryHandlerKey, lp2p.DiscoveryHandler),

	// Routing
	Override(new(record.Validator), modules.RecordValidator),
	Override(BaseRoutingKey, lp2p.BaseRouting),
	Override(new(routing.Routing), lp2p.Routing),

	// Services
	Override(BandwidthReporterKey, lp2p.BandwidthCounter),
	Override(AutoNATSvcKey, lp2p.AutoNATService),

	// Services (pubsub)
	Override(new(*dtypes.ScoreKeeper), lp2p.ScoreKeeper),
	Override(new(*pubsub.PubSub), lp2p.GossipSub),
	Override(new(*config.Pubsub), func(bs dtypes.Bootstrapper) *config.Pubsub {
		return &config.Pubsub{
			Bootstrapper: bool(bs),
		}
	}),

	// Services (connection management)
	Override(ConnectionManagerKey, lp2p.ConnectionManager(50, 200, 20*time.Second, nil)),
	Override(new(*conngater.BasicConnectionGater), lp2p.ConnGater),
	Override(ConnGaterKey, lp2p.ConnGaterOption),
)

func Base() Option {
	return Options(
		func(s *Settings) error { s.Base = true; return nil }, // mark Base as applied
		ApplyIf(func(s *Settings) bool { return s.Config },
			Error(errors.New("the Base() option must be set before Config option")),
		),
		ApplyIf(func(s *Settings) bool { return s.enableLibp2pNode },
			LibP2P,
		),
	)
}

func Repo(r repo.Repo) Option {
	return func(settings *Settings) error {
		lr, err := r.Lock(settings.nodeType)
		if err != nil {
			return err
		}
		c, err := lr.Config()
		if err != nil {
			return err
		}
		return Options(
			Override(new(repo.LockedRepo), modules.LockedRepo(lr)), // module handles closing

			Override(new(ci.PrivKey), lp2p.PrivKey),
			Override(new(ci.PubKey), ci.PrivKey.GetPublic),
			Override(new(peer.ID), peer.IDFromPublicKey),

			Override(new(types.KeyStore), modules.KeyStore),
			Override(new(*dtypes.APIAlg), modules.APISecret),

			ApplyIf(IsType(repo.Dbridge), ConfigFullNode(c)),
		)(settings)
	}
}

type StopFunc func(context.Context) error

// New builds and starts new Filecoin node
func New(ctx context.Context, opts ...Option) (StopFunc, error) {
	settings := Settings{
		modules: map[interface{}]fx.Option{},
		invokes: make([]fx.Option, _nInvokes),
	}

	// apply module options in the right order
	if err := Options(Options(defaults()...), Options(opts...))(&settings); err != nil {
		return nil, xerrors.Errorf("applying node options failed: %w", err)
	}

	// gather constructors for fx.Options
	ctors := make([]fx.Option, 0, len(settings.modules))
	for _, opt := range settings.modules {
		ctors = append(ctors, opt)
	}

	// fill holes in invokes for use in fx.Options
	for i, opt := range settings.invokes {
		if opt == nil {
			settings.invokes[i] = fx.Options()
		}
	}

	app := fx.New(
		fx.Options(ctors...),
		fx.Options(settings.invokes...),

		fx.NopLogger,
	)

	// TODO: we probably should have a 'firewall' for Closing signal
	//  on this context, and implement closing logic through lifecycles
	//  correctly
	if err := app.Start(ctx); err != nil {
		// comment fx.NopLogger few lines above for easier debugging
		return nil, xerrors.Errorf("starting node: %w", err)
	}

	return app.Stop, nil
}

type FullOption = Option

func Lite(enable bool) FullOption {
	return func(s *Settings) error {
		s.Lite = enable
		return nil
	}
}

func FullAPI(out *api.FullNode, fopts ...FullOption) Option {
	return Options(
		func(s *Settings) error {
			s.nodeType = repo.Dbridge
			s.enableLibp2pNode = true
			return nil
		},
		Options(fopts...),
		func(s *Settings) error {
			resAPI := &impl.FullNodeAPI{}
			s.invokes[ExtractApiKey] = fx.Populate(resAPI)
			*out = resAPI
			return nil
		},
	)
}

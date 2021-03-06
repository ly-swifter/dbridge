package main

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/lyswifter/dbridge/api"
	"github.com/lyswifter/dbridge/lib/peermgr"
	"github.com/lyswifter/dbridge/node"
	"github.com/lyswifter/dbridge/node/modules/dtypes"
	"github.com/lyswifter/dbridge/node/repo"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var StopCmd = &cli.Command{
	Name:  "stop",
	Usage: "Stop a running dbridge process",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		// api, closer, err := lcli.GetAPI(cctx)
		// if err != nil {
		// 	return err
		// }
		// defer closer()

		// err = api.Shutdown(lcli.ReqContext(cctx))
		// if err != nil {
		// 	return err
		// }

		return nil
	},
}

var RunCmd = &cli.Command{
	Name:  "run",
	Usage: "Start running a dbridge process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "api",
			Value: "1234",
		},
		&cli.StringFlag{
			Name:    FlagDbridgeRepo,
			EnvVars: []string{"DBRIDGE_NODE_PATH"},
			Hidden:  true,
			Value:   "~/.lorry",
			Usage:   "Specify dbridge node repo path.",
		},
		&cli.IntFlag{
			Name:  "api-max-req-size",
			Usage: "maximum API request size accepted by the JSON RPC server",
		},
		&cli.BoolFlag{
			Name:   "lite",
			Usage:  "start lorry in lite mode",
			Hidden: false,
		},
	},
	Action: func(cctx *cli.Context) error {

		// ctx, _ := tag.New(context.Background(),
		// 	tag.Insert(metrics.NodeType, "dkg-node"),
		// )

		ctx := context.Background()

		isLite := cctx.Bool("lite")

		shutdownChan := make(chan struct{})

		bridgeRepoPath := cctx.String(FlagDbridgeRepo)
		r, err := repo.NewFS(bridgeRepoPath)
		if err != nil {
			return err
		}

		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if !ok {
			return xerrors.Errorf("repo at '%s' is not initialized, run 'lorry init' to set it up", bridgeRepoPath)
		}

		var api api.FullNode
		stop, err := node.New(ctx,
			node.FullAPI(&api, node.Lite(isLite)),

			node.Base(),
			node.Repo(r),

			node.Override(new(dtypes.ShutdownChan), shutdownChan),

			node.ApplyIf(func(s *node.Settings) bool { return cctx.IsSet("api") },
				node.Override(node.SetApiEndpointKey, func(lr repo.LockedRepo) error {
					apima, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/" +
						cctx.String("api"))
					if err != nil {
						return err
					}
					return lr.SetAPIEndpoint(apima)
				})),
			node.ApplyIf(func(s *node.Settings) bool { return !cctx.Bool("bootstrap") },
				node.Unset(node.RunPeerMgrKey),
				node.Unset(new(*peermgr.PeerMgr)),
			),
		)
		if err != nil {
			return xerrors.Errorf("initializing node: %w", err)
		}

		endpoint, err := r.APIEndpoint()
		if err != nil {
			return xerrors.Errorf("getting api endpoint: %w", err)
		}

		//
		// Instantiate JSON-RPC endpoint.
		// ---

		// Populate JSON-RPC options.
		serverOptions := make([]jsonrpc.ServerOption, 0)
		if maxRequestSize := cctx.Int("api-max-req-size"); maxRequestSize != 0 {
			serverOptions = append(serverOptions, jsonrpc.WithMaxRequestSize(int64(maxRequestSize)))
		}

		// Instantiate the full node handler.
		h, err := node.FullNodeHandler(api, true, serverOptions...)
		if err != nil {
			return fmt.Errorf("failed to instantiate rpc handler: %s", err)
		}

		// Serve the RPC.
		rpcStopper, err := node.ServeRPC(h, "lorry-daemon", endpoint)
		if err != nil {
			return fmt.Errorf("failed to start json-rpc endpoint: %s", err)
		}

		// Monitor for shutdown.
		finishCh := node.MonitorShutdown(shutdownChan,
			node.ShutdownHandler{Component: "rpc server", StopFunc: rpcStopper},
			node.ShutdownHandler{Component: "node", StopFunc: stop},
		)
		<-finishCh // fires when shutdown is complete.

		return nil
	},
	Subcommands: []*cli.Command{
		StopCmd,
	},
}

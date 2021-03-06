package main

import (
	"context"

	logging "github.com/ipfs/go-log/v2"
	lcli "github.com/lyswifter/dbridge/cli"
	cliutil "github.com/lyswifter/dbridge/cli/util"
	"github.com/lyswifter/dbridge/node/repo"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/trace"
)

var log = logging.Logger("main")

const (
	FlagDbridgeRepo = "repo"
)

var AdvanceBlockCmd *cli.Command

func main() {
	local := []*cli.Command{
		sampleCmd,
		initCmd,
		RunCmd,
	}

	if AdvanceBlockCmd != nil {
		local = append(local, AdvanceBlockCmd)
	}

	ctx, span := trace.StartSpan(context.Background(), "/cli")
	defer span.End()

	app := &cli.App{
		Name:                 "dbridge",
		Usage:                "Dbridge decentralized bridge network client",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    FlagDbridgeRepo,
				EnvVars: []string{"DBRIDGE_NODE_PATH"},
				Hidden:  true,
				Value:   "~/.lorry",
				Usage:   "Specify dbridge node repo path.",
			},
			cliutil.FlagVeryVerbose,
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Commands: append(local, lcli.Commands...),
	}

	app.Setup()
	app.Metadata["traceContext"] = ctx
	app.Metadata["repoType"] = repo.Dbridge

	lcli.RunApp(app)
}

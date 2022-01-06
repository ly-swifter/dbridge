package main

import (
	"github.com/ipfs/go-datastore"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"

	lcli "github.com/lyswifter/dbridge/cli/util"
	"github.com/lyswifter/dbridge/node/repo"
	"golang.org/x/xerrors"
)

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "Initialize a dbridge node repo",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		log.Info("Initializing dbridge node")

		ctx := lcli.ReqContext(c)

		repoPath := c.String("repo")

		log.Infof("repo name: %s", repoPath)

		{
			dir, err := homedir.Expand(repoPath)
			if err != nil {
				log.Warnw("could not expand repo location", "error", err)
			} else {
				log.Infof("dbridge repo: %s", dir)
			}
		}

		r, err := repo.NewFS(repoPath)
		if err != nil {
			return xerrors.Errorf("opening fs repo: %w", err)
		}

		err = r.Init(repo.Dbridge)
		if err != nil && err != repo.ErrRepoExists {
			return xerrors.Errorf("repo init error: %w", err)
		}

		lr, err := r.Lock(repo.Dbridge)
		if err != nil {
			return err
		}
		defer lr.Close()

		mds, err := lr.Datastore(ctx, "/metadata")
		if err != nil {
			return err
		}

		nodeName := "node-1"

		if err := mds.Put(ctx, datastore.NewKey("node-address"), []byte(nodeName)); err != nil {
			return err
		}

		log.Infof("Initializing dbridge node(%s) at %s, you can now running node with: ./dbridge daemon", nodeName, FlagDbridgeRepo)

		return nil
	},
}

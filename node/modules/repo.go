package modules

import (
	"context"
	"path/filepath"

	"github.com/lyswifter/dbridge/lib/backupds"
	"github.com/lyswifter/dbridge/node/modules/dtypes"
	"github.com/lyswifter/dbridge/node/modules/helpers"
	"github.com/lyswifter/dbridge/node/repo"
	"github.com/lyswifter/dbridge/types"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

func LockedRepo(lr repo.LockedRepo) func(lc fx.Lifecycle) repo.LockedRepo {
	return func(lc fx.Lifecycle) repo.LockedRepo {
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return lr.Close()
			},
		})

		return lr
	}
}

func KeyStore(lr repo.LockedRepo) (types.KeyStore, error) {
	return lr.KeyStore()
}

func Datastore(disableLog bool) func(lc fx.Lifecycle, mctx helpers.MetricsCtx, r repo.LockedRepo) (dtypes.MetadataDS, error) {
	return func(lc fx.Lifecycle, mctx helpers.MetricsCtx, r repo.LockedRepo) (dtypes.MetadataDS, error) {
		ctx := helpers.LifecycleCtx(mctx, lc)
		mds, err := r.Datastore(ctx, "/metadata")
		if err != nil {
			return nil, err
		}

		var logdir string
		if !disableLog {
			logdir = filepath.Join(r.Path(), "kvlog/metadata")
		}

		bds, err := backupds.Wrap(mds, logdir)
		if err != nil {
			return nil, xerrors.Errorf("opening backupds: %w", err)
		}

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return bds.CloseLog()
			},
		})

		return bds, nil
	}
}

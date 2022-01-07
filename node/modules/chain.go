package modules

import (
	"github.com/lyswifter/dbridge/build"
	"github.com/lyswifter/dbridge/node/modules/dtypes"
	"github.com/lyswifter/dbridge/node/modules/helpers"
	"go.uber.org/fx"
)

func NetworkName(mctx helpers.MetricsCtx, lc fx.Lifecycle) (dtypes.NetworkName, error) {
	if !build.Devnet {
		return "mainnet", nil
	}

	return "testnet", nil
}

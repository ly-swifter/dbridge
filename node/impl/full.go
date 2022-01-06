package impl

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/lyswifter/dbridge/api"
	"github.com/lyswifter/dbridge/node/impl/common"
)

var log = logging.Logger("node")

type FullNodeAPI struct {
	common.CommonAPI

	//more
}

var _ api.FullNode = &FullNodeAPI{}

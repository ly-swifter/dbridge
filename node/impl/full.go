package impl

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/lyswifter/dbridge/api"
	"github.com/lyswifter/dbridge/node/impl/common"
	"github.com/lyswifter/dbridge/node/impl/net"
)

var log = logging.Logger("node")

type FullNodeAPI struct {
	common.CommonAPI
	net.NetAPI

	//more
}

var _ api.FullNode = &FullNodeAPI{}

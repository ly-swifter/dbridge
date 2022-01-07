package build

import (
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/lyswifter/dbridge/node/modules/dtypes"
)

func BlocksTopic(netName dtypes.NetworkName) string   { return "/lorry/blocks/" + string(netName) }
func MessagesTopic(netName dtypes.NetworkName) string { return "/lorry/msgs/" + string(netName) }
func DhtProtocolName(netName dtypes.NetworkName) protocol.ID {
	return protocol.ID("/lorry/kad/" + string(netName))
}

// /////
// Devnet settings

var Devnet = true

package client

import (
	"context"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/lyswifter/dbridge/api"
)

// NewFullNodeRPCV0 creates a new http jsonrpc client.
func NewFullNodeRPCV0(ctx context.Context, addr string, requestHeader http.Header) (api.FullNode, jsonrpc.ClientCloser, error) {
	var res api.FullNodeStruct

	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Dbridge",
		api.GetInternalStructs(&res), requestHeader)

	return &res, closer, err
}

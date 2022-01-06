package cli

import (
	logging "github.com/ipfs/go-log/v2"

	cliutil "github.com/lyswifter/dbridge/cli/util"
)

var log = logging.Logger("cli")

var GetAPIInfo = cliutil.GetAPIInfo
var GetRawAPI = cliutil.GetRawAPI
var GetAPI = cliutil.GetCommonAPI

var DaemonContext = cliutil.DaemonContext
var ReqContext = cliutil.ReqContext

var GetFullNodeAPI = cliutil.GetFullNodeAPI

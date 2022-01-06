package node

import (
	"context"
	"net"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"github.com/lyswifter/dbridge/api"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"golang.org/x/xerrors"
)

var rpclog = logging.Logger("rpc")

// ServeRPC serves an HTTP handler over the supplied listen multiaddr.
//
// This function spawns a goroutine to run the server, and returns immediately.
// It returns the stop function to be called to terminate the endpoint.
//
// The supplied ID is used in tracing, by inserting a tag in the context.
func ServeRPC(h http.Handler, id string, addr multiaddr.Multiaddr) (StopFunc, error) {
	// Start listening to the addr; if invalid or occupied, we will fail early.
	lst, err := manet.Listen(addr)
	if err != nil {
		return nil, xerrors.Errorf("could not listen: %w", err)
	}

	// Instantiate the server and start listening.
	srv := &http.Server{
		Handler: h,
		BaseContext: func(listener net.Listener) context.Context {
			// ctx, _ := tag.New(context.Background(), tag.Upsert(metrics.APIInterface, id))
			ctx := context.Background()
			return ctx
		},
	}

	go func() {
		err = srv.Serve(manet.NetListener(lst))
		if err != http.ErrServerClosed {
			rpclog.Warnf("rpc server failed: %s", err)
		}
	}()

	return srv.Shutdown, err
}

func metricedFullAPI(a api.FullNode) api.FullNode {
	var out api.FullNodeStruct
	//Fixme: something miss here
	return &out
}

// FullNodeHandler returns a full node handler, to be mounted as-is on the server.
func FullNodeHandler(a api.FullNode, permissioned bool, opts ...jsonrpc.ServerOption) (http.Handler, error) {
	m := mux.NewRouter()

	serveRpc := func(path string, hnd interface{}) {
		rpcServer := jsonrpc.NewServer(opts...)
		rpcServer.Register("Dbridge", hnd)

		var handler http.Handler = rpcServer
		if permissioned {
			handler = &auth.Handler{Verify: a.AuthVerify, Next: rpcServer.ServeHTTP}
		}

		m.Handle(path, handler)
	}

	//Fixme: impl real api methods
	fnapi := metricedFullAPI(a)
	if permissioned {
		fnapi = api.PermissionedFullAPI(a)
	}

	// serveRpc("/rpc/v1", fnapi)
	// serveRpc("/rpc/v0", &v0api.WrapperV1Full{FullNode: fnapi})

	serveRpc("/rpc/v0", &fnapi)

	// Import handler
	// handleImportFunc := handleImport(a.(*impl.FullNodeAPI))
	// handleExportFunc := handleExport(a.(*impl.FullNodeAPI))
	// if permissioned {
	// 	importAH := &auth.Handler{
	// 		Verify: a.AuthVerify,
	// 		Next:   handleImportFunc,
	// 	}
	// 	m.Handle("/rest/v0/import", importAH)

	// 	exportAH := &auth.Handler{
	// 		Verify: a.AuthVerify,
	// 		Next:   handleExportFunc,
	// 	}
	// 	m.Handle("/rest/v0/export", exportAH)
	// } else {
	// 	m.HandleFunc("/rest/v0/import", handleImportFunc)
	// 	m.HandleFunc("/rest/v0/export", handleExportFunc)
	// }

	// debugging
	// m.Handle("/debug/metrics", metrics.Exporter())
	// m.Handle("/debug/pprof-set/block", handleFractionOpt("BlockProfileRate", runtime.SetBlockProfileRate))
	// m.Handle("/debug/pprof-set/mutex", handleFractionOpt("MutexProfileFraction", func(x int) {
	// 	runtime.SetMutexProfileFraction(x)
	// }))
	// m.PathPrefix("/").Handler(http.DefaultServeMux) // pprof

	return m, nil
}

// func handleImport(a *impl.FullNodeAPI) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != "PUT" {
// 			w.WriteHeader(404)
// 			return
// 		}
// 		if !auth.HasPerm(r.Context(), nil, api.PermWrite) {
// 			w.WriteHeader(401)
// 			_ = json.NewEncoder(w).Encode(struct{ Error string }{"unauthorized: missing write permission"})
// 			return
// 		}

// 		c, err := a.ClientImportLocal(r.Context(), r.Body)
// 		if err != nil {
// 			w.WriteHeader(500)
// 			_ = json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
// 			return
// 		}
// 		w.WriteHeader(200)
// 		err = json.NewEncoder(w).Encode(struct{ Cid cid.Cid }{c})
// 		if err != nil {
// 			rpclog.Errorf("/rest/v0/import: Writing response failed: %+v", err)
// 			return
// 		}
// 	}
// }

// func handleExport(a *impl.FullNodeAPI) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != "GET" {
// 			w.WriteHeader(404)
// 			return
// 		}
// 		if !auth.HasPerm(r.Context(), nil, api.PermWrite) {
// 			w.WriteHeader(401)
// 			_ = json.NewEncoder(w).Encode(struct{ Error string }{"unauthorized: missing write permission"})
// 			return
// 		}

// 		var eref api.ExportRef
// 		if err := json.Unmarshal([]byte(r.FormValue("export")), &eref); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		car := r.FormValue("car") == "true"

// 		err := a.ClientExportInto(r.Context(), eref, car, client.ExportDest{Writer: w})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

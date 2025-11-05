package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/0xSplits/pulsar/pkg/runtime"
	"github.com/0xSplits/pulsar/pkg/server/handler"
	"github.com/0xSplits/pulsar/pkg/server/websocket"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twitchtv/twirp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Han are the server specific handlers implementing the actual business
	// logic.
	Han []handler.Interface
	// Int are the Twirp specific interceptors wrapping the endpoint handlers.
	Int []twirp.Interceptor
	// Lis is the main HTTP listener bound to some configured host and port.
	Lis net.Listener
	// Log is the structured logger passed down the stack.
	Log logger.Interface
	// Mid are the protocol specific transport layer middlewares executed before
	// any RPC handler.
	Mid []mux.MiddlewareFunc
	// Web is the websocket handler receiving all the token transfer events from
	// our indexing provider.
	Web *websocket.Handler
}

type Server struct {
	lis net.Listener
	log logger.Interface
	srv *http.Server
}

func New(c Config) *Server {
	if len(c.Han) == 0 {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Han must not be empty", c)))
	}
	if c.Lis == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Lis must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	var rtr *mux.Router
	{
		rtr = mux.NewRouter()
	}

	{
		rtr.Use(c.Mid...)
	}

	// Add a simple health check response to the root.
	{
		rtr.NewRoute().Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(linBrk([]byte("OK")))
		})
	}

	// Add the anubis streaming handler. All GET requests will be upgraded to
	// manage websocket connections.
	{
		rtr.NewRoute().Methods("GET").Path("/indexing").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := c.Web.HandlerFunc(w, r)
			if websocket.IsInvalidWebsocketSecret(err) {
				w.WriteHeader(http.StatusUnauthorized)
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			{
				w.Header().Set("Content-Type", "text/plain")
			}

			if err != nil {
				c.Log.Log(
					"level", "error",
					"message", "request failed",
					"stack", tracer.Json(err),
				)

				{
					_, _ = w.Write(linBrk([]byte(err.Error())))
				}
			}
		})
	}

	// Add the metrics endpoint in Prometehus format.
	{
		rtr.NewRoute().Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	}

	// Add a simple version response for the runtime.
	{
		rtr.NewRoute().Methods("GET").Path("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(linBrk(runtime.JSON()))
		})
	}

	for _, x := range c.Han {
		x.Attach(rtr, twirp.WithServerInterceptors(c.Int...), twirp.WithServerPathPrefix(""))
	}

	return &Server{
		lis: c.Lis,
		log: c.Log,
		srv: &http.Server{
			Handler: rtr,
		},
	}
}

func (s *Server) Daemon() {
	s.log.Log(
		"level", "info",
		"message", "server is accepting calls",
		"address", s.lis.Addr().String(),
	)

	{
		err := s.srv.Serve(s.lis)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}
}

func linBrk(byt []byte) []byte {
	return append(byt, []byte("\n")...)
}

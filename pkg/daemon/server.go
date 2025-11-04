package daemon

import (
	"net"

	"github.com/0xSplits/pulsar/pkg/server"
	"github.com/0xSplits/pulsar/pkg/server/handler"
	"github.com/0xSplits/pulsar/pkg/server/handler/transfers"
	"github.com/0xSplits/pulsar/pkg/server/middleware/cors"
	"github.com/0xSplits/pulsar/pkg/server/websocket"
	"github.com/gorilla/mux"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Server() *server.Server {
	var err error

	var lis net.Listener
	{
		lis, err = net.Listen("tcp", net.JoinHostPort(d.env.HttpHost, d.env.HttpPort))
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var web *websocket.Handler
	{
		web = websocket.New(websocket.Config{
			Env: d.env,
			Log: d.log,
		})
	}

	var ser *server.Server
	{
		ser = server.New(server.Config{
			Han: []handler.Interface{
				transfers.New(transfers.Config{
					Log: d.log,
				}),
			},
			Lis: lis,
			Log: d.log,
			Mid: []mux.MiddlewareFunc{
				cors.New(cors.Config{Log: d.log}).Handler,
			},
			Web: web,
		})
	}

	return ser
}

package websocket

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/coder/websocket"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) HandlerFunc(w http.ResponseWriter, r *http.Request) error {
	var err error

	h.log.Log(
		"level", "info",
		"message", "received websocket request",
	)

	if !h.verify(r.Header.Get("Authorization")) {
		return tracer.Mask(invalidWebsocketSecretError)
	}

	var opt *websocket.AcceptOptions
	{
		opt = &websocket.AcceptOptions{}
	}

	var con *websocket.Conn
	{
		con, err = websocket.Accept(w, r, opt)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for {
		// Read the next incoming message from the connected client.

		_, byt, err := con.Read(context.Background())
		if errors.Is(err, net.ErrClosed) {
			return nil
		} else if err != nil {
			return tracer.Mask(err)
		}

		// Process the incoming message.

		err = h.process(byt)
		if err != nil {
			return tracer.Mask(err)
		}
	}
}

func (h *Handler) process(byt []byte) error {
	// TODO forward to Supabase/Postgres
	fmt.Printf("%s\n", byt)
	return nil
}

func (h *Handler) verify(aut string) bool {
	tok := strings.TrimPrefix(aut, "Bearer ")
	trm := strings.TrimSpace(tok)

	return trm == h.env.WebsocketSecret
}

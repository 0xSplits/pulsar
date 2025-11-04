package websocket

import (
	"fmt"

	"github.com/0xSplits/pulsar/pkg/envvar"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Env envvar.Env
	Log logger.Interface
}

type Handler struct {
	env envvar.Env
	log logger.Interface
}

func New(c Config) *Handler {
	if c.Env.WebsocketSecret == "" {
		tracer.Panic(fmt.Errorf("%T.Env.WebsocketSecret must not be empty", c))
	}
	if c.Log == nil {
		tracer.Panic(fmt.Errorf("%T.Log must not be empty", c))
	}

	return &Handler{
		env: c.Env,
		log: c.Log,
	}
}

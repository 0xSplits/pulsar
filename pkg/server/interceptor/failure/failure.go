package failure

import (
	"context"
	"fmt"

	"github.com/twitchtv/twirp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Log logger.Interface
}

type Interceptor struct {
	log logger.Interface
}

func New(c Config) *Interceptor {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Interceptor{
		log: c.Log,
	}
}

func (i *Interceptor) Method(nex twirp.Method) twirp.Method {
	return func(ctx context.Context, req any) (any, error) {
		res, err := nex(ctx, req)
		if err != nil {
			i.log.Log(
				"level", "error",
				"message", "request failed",
				"stack", tracer.Json(err),
			)
		}

		return res, nil
	}
}

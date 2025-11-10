package accounts

import (
	"fmt"

	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/accounts/client"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Cli client.Interface
}

type Accounts struct {
	cli client.Interface
}

func New(c Config) *Accounts {
	if c.Cli == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Cli must not be empty", c)))
	}

	return &Accounts{
		cli: c.Cli,
	}
}

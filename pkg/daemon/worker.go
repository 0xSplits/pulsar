package daemon

import (
	"fmt"

	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline"
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/accounts"
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/accounts/client"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/registry"
	"github.com/0xSplits/workit/worker/parallel"
)

func (d *Daemon) Worker() *parallel.Worker {
	var cli *client.Client
	{
		cli = client.New(client.Config{
			Url: fmt.Sprintf("https://server.%s.splits.org", d.env.Environment),
		})
	}

	var acc *accounts.Accounts
	{
		acc = accounts.New(accounts.Config{
			Cli: cli,
		})
	}

	var reg *registry.Registry
	{
		reg = registry.New(registry.Config{
			Env: d.env.Environment,
			Log: d.log,
			Met: d.met,
		})
	}

	var par *parallel.Worker
	{
		par = parallel.New(parallel.Config{
			Han: []handler.Cooler{
				pipeline.New(pipeline.Config{Acc: acc, Env: d.env, Fil: d.fil, Log: d.log}),
			},
			Log: d.log,
			Reg: reg,
		})
	}

	return par
}

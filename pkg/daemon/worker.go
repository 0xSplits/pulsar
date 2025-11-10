package daemon

import (
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/registry"
	"github.com/0xSplits/workit/worker/parallel"
)

func (d *Daemon) Worker() *parallel.Worker {
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
				pipeline.New(pipeline.Config{Log: d.log}),
			},
			Log: d.log,
			Reg: reg,
		})
	}

	return par
}

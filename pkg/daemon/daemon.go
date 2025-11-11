package daemon

import (
	"github.com/0xSplits/indexingo/client"
	"github.com/0xSplits/indexingo/filters"
	"github.com/0xSplits/indexingo/pipelines"
	"github.com/0xSplits/indexingo/transformations"
	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/pulsar/pkg/envvar"
	"github.com/0xSplits/pulsar/pkg/runtime"
	"github.com/xh3b4sd/logger"
	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	Env envvar.Env
}

type Daemon struct {
	env envvar.Env
	fil *filters.Filters
	log logger.Interface
	met metric.Meter
	pip *pipelines.Pipelines
	tra *transformations.Transformations
}

func New(c Config) *Daemon {
	var log logger.Interface
	{
		log = logger.New(logger.Config{
			Filter: logger.NewLevelFilter(c.Env.LogLevel),
		})
	}

	var met metric.Meter
	{
		met = recorder.NewMeter(recorder.MeterConfig{
			Env: c.Env.Environment,
			Sco: "pulsar",
			Ver: runtime.Tag(),
		})
	}

	var cli client.Interface
	{
		cli = client.New(client.Config{
			Key: c.Env.IndexingcoApiKey,
		})
	}

	var fil *filters.Filters
	{
		fil = filters.New(filters.Config{
			Cli: cli,
		})
	}

	var tra *transformations.Transformations
	{
		tra = transformations.New(transformations.Config{
			Cli: cli,
		})
	}

	var pip *pipelines.Pipelines
	{
		pip = pipelines.New(pipelines.Config{
			Cli: cli,
		})
	}

	log.Log(
		"level", "info",
		"message", "daemon is launching procs",
		"environment", c.Env.Environment,
	)

	return &Daemon{
		env: c.Env,
		fil: fil,
		log: log,
		met: met,
		pip: pip,
		tra: tra,
	}
}

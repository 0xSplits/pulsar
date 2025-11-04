package daemon

import (
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
	log logger.Interface
	met metric.Meter
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

	log.Log(
		"level", "info",
		"message", "daemon is launching procs",
		"environment", c.Env.Environment,
	)

	return &Daemon{
		env: c.Env,
		log: log,
		met: met,
	}
}

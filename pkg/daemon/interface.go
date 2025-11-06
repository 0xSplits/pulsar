package daemon

import (
	"github.com/0xSplits/pulsar/pkg/server"
	"github.com/0xSplits/workit/worker/parallel"
)

type Interface interface {
	Ensure() error
	Server() *server.Server
	Worker() *parallel.Worker
}

package pipeline

import (
	"fmt"
	"time"

	"github.com/0xSplits/indexingo/filters"
	"github.com/0xSplits/pulsar/pkg/envvar"
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/accounts"
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/cursor"
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/locker"
	"github.com/xh3b4sd/choreo/jitter"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Acc is the accounts client to search iteratively for all relevant wallet
	// addresses known to the system.
	Acc *accounts.Accounts
	// Env is the configuration injected into the process environment.
	Env envvar.Env
	// Fil is the pipeline filter kept up to date by this worker handler.
	Fil *filters.Filters
	// Log is the standard logger interface to emit useful log messages at
	// runtime.
	Log logger.Interface
}

type Handler struct {
	acc *accounts.Accounts
	cur *cursor.Cursor
	env envvar.Env
	fil *filters.Filters
	jit *jitter.Jitter[time.Duration]
	loc *locker.Locker
	log logger.Interface
}

func New(c Config) *Handler {
	if c.Acc == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Acc must not be empty", c)))
	}
	if c.Env.Environment == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Fil == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Fil must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	// Note that the setup below does only relate to the cursor specific database
	// connection and table. This cursor client should ideally not live in this
	// constructor, since it is an antipattern to allow third party dependencies
	// to cause worker handler failure on creation seemingly randomly. For the
	// time being, keeping this code here is just easy and convenient to get the
	// first version going.

	var err error

	var cur *cursor.Cursor
	{
		cur = cursor.New(cursor.Config{
			Dsn: c.Env.PostgresUrl,
			Log: c.Log,
		})
	}

	{
		err = cur.Create()
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	// Note that the setup below does only relate to the locker specific database
	// connection and table. This locker client should ideally not live in this
	// constructor, since it is an antipattern to allow third party dependencies
	// to cause worker handler failure on creation seemingly randomly. For the
	// time being, keeping this code here is just easy and convenient to get the
	// first version going.

	var loc *locker.Locker
	{
		loc = locker.New(locker.Config{
			Dsn: c.Env.PostgresUrl,
			Log: c.Log,
		})
	}

	{
		err = loc.Create()
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return &Handler{
		acc: c.Acc,
		cur: cur,
		env: c.Env,
		fil: c.Fil,
		jit: jitter.New[time.Duration](jitter.Config{Per: 0.20}),
		loc: loc,
		log: c.Log,
	}
}

package locker

import (
	"database/sql"
	"fmt"
	"time"

	"cirello.io/pglock"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Dsn is the data source name for the underlying Postgres instance in order
	// to establish a new database connection.
	Dsn string
	// Log is the standard logger interface to emit useful log messages at
	// runtime.
	Log logger.Interface
}

type Locker struct {
	cli *pglock.Client
	log logger.Interface
}

func New(c Config) *Locker {
	if c.Dsn == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Dsn must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	var err error

	var dat *sql.DB
	{
		dat, err = sql.Open("postgres", c.Dsn)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var opt []pglock.ClientOption
	{
		opt = []pglock.ClientOption{
			pglock.WithLeaseDuration(5 * time.Second),
			pglock.WithHeartbeatFrequency(1 * time.Second),
		}
	}

	var cli *pglock.Client
	{
		cli, err = pglock.New(dat, opt...)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return &Locker{
		cli: cli,
		log: c.Log,
	}
}

package cursor

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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

type Cursor struct {
	log logger.Interface
	poo *pgxpool.Pool
}

func New(c Config) *Cursor {
	if c.Dsn == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Dsn must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	var err error

	// The connection pool set up here will be used for executing Postgres
	// commands related to our actual business logic. In other words, this
	// connection pool is unrelated to the database connection managed for the
	// distributed lock.

	var poo *pgxpool.Pool
	{
		poo, err = pgxpool.New(context.Background(), c.Dsn)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return &Cursor{
		log: c.Log,
		poo: poo,
	}
}

package cursor

import (
	"context"
	"time"

	"github.com/xh3b4sd/tracer"
)

const (
	create = `
		--
	  -- create the database table if it does not exists
		--
		CREATE TABLE IF NOT EXISTS accounts_pagination_cursor (
			id         boolean        PRIMARY KEY DEFAULT true,
			cursor     timestamptz    NOT NULL,
			updated    timestamptz    NOT NULL DEFAULT now(),
			CHECK (id)
		);

		--
	  -- create the default cursor if it does not exists
		--
		INSERT INTO accounts_pagination_cursor (id, cursor)
		VALUES (true, '2020-01-01T00:00:00.000Z')
		ON CONFLICT (id) DO NOTHING;
	`
)

func (c *Cursor) Create() error {
	var err error

	var ctx context.Context
	var can context.CancelFunc
	{
		ctx, can = context.WithTimeout(context.Background(), 3*time.Second)
	}

	{
		defer can()
	}

	{
		_, err = c.poo.Exec(ctx, create)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

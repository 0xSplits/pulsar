package cursor

import (
	"context"
	"time"

	"github.com/xh3b4sd/tracer"
)

const (
	search = `
		--
	  -- search the current pagination cursor
		--
    SELECT cursor
    FROM accounts_pagination_cursor
    WHERE id = true
	`
)

func (c *Cursor) Search() (time.Time, error) {
	var err error

	var ctx context.Context
	var can context.CancelFunc
	{
		ctx, can = context.WithTimeout(context.Background(), 3*time.Second)
	}

	{
		defer can()
	}

	var cur time.Time
	{
		err = c.poo.QueryRow(ctx, search).Scan(&cur)
		if err != nil {
			return time.Time{}, tracer.Mask(err)
		}
	}

	return cur, nil
}

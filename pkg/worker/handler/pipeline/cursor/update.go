package cursor

import (
	"context"
	"time"

	"github.com/xh3b4sd/tracer"
)

const (
	update = `
		--
	  -- update the pagination cursor to set its new version
		--
    UPDATE accounts_pagination_cursor
    SET
      cursor  = $1,
      updated = now()
    WHERE id = true
	`
)

func (c *Cursor) Update(cur time.Time) error {
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
		_, err = c.poo.Exec(ctx, update, cur)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

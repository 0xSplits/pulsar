package locker

import (
	"context"
	"time"

	"cirello.io/pglock"
	"github.com/xh3b4sd/tracer"
)

func (l *Locker) Release(loc *pglock.Lock) error {
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
		err = l.cli.ReleaseContext(ctx, loc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

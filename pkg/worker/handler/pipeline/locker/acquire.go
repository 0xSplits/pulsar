package locker

import (
	"context"
	"errors"
	"time"

	"cirello.io/pglock"
	"github.com/xh3b4sd/tracer"
)

func (l *Locker) Acquire() (*pglock.Lock, error) {
	var err error

	var ctx context.Context
	var can context.CancelFunc
	{
		ctx, can = context.WithTimeout(context.Background(), 3*time.Second)
	}

	{
		defer can()
	}

	var loc *pglock.Lock
	{
		loc, err = l.cli.AcquireContext(ctx, "accounts")
		if errors.Is(err, pglock.ErrNotAcquired) {
			return nil, nil
		} else if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return loc, nil
}

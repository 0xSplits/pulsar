package locker

import (
	"github.com/xh3b4sd/tracer"
)

// Create the database table required for the lock client to work properly, but
// only create the table if it does not already exist. So in case the database
// table is already setup properly, this call should not fail.
func (l *Locker) Create() error {
	var err error

	{
		err = l.cli.TryCreateTable()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

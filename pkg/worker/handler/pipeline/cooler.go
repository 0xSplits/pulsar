package pipeline

import (
	"time"
)

// Cooler is configured to return a dynamically adjusted wait duration for this
// worker handler to sleep before running again. The introduced jitter has the
// purpose of spreading out the same type of work across time, so that we ease
// the load on our dependency APIs, here Postgres, and effectively try to
// prevent unneccessary contention. E.g. a jitter of 20% applied to 10s results
// in execution variation of +-2s.
func (h *Handler) Cooler() time.Duration {
	return h.jit.Percent(10 * time.Second)
}

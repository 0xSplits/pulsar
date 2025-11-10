package pipeline

import (
	"time"

	"cirello.io/pglock"
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/accounts"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	// TODO add logging

	var loc *pglock.Lock
	{
		loc, err = h.loc.Acquire()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if loc == nil {
		return nil // lock not acquired
	}

	{
		defer h.loc.Release(loc) // nolint:errcheck
	}

	var cur time.Time
	{
		cur, err = h.cur.Search()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var res accounts.SearchResponse
	{
		res, err = h.acc.Search(cur.Unix())
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var add []string
	for _, x := range res.Data {
		add = append(add, x.Address)
	}

	// Adding values to the filter is idempotent as per the underlying indexing
	// provider Indexing Co. It is therefore not necessary to de-duplicate the
	// results received from the accounts/search endpoint.

	{
		_, err = h.fil.AddValues("test-filter", add) // TODO fix the filter name, add environment
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// There is potentially a gap issue between consecutive calls of this worker
	// handler. So we should try to ensure a little bit of cursor overlap across
	// consecutive calls in order to prevent any gaps in our observed accounts
	// list. Note that setting back the next cursor by 30 seconds will cause
	// duplicated data to be received from the accounts/search API. Though this
	// should not be an issue as described above.

	var nxt time.Time
	if res.Next != nil {
		nxt = time.Unix(*res.Next-30, 0)
	} else {
		nxt = time.Now().UTC()
	}

	{
		err = h.cur.Update(nxt)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

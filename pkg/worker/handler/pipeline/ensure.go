package pipeline

import (
	"fmt"
	"strconv"
	"time"

	"cirello.io/pglock"
	"github.com/0xSplits/pulsar/pkg/worker/handler/pipeline/accounts"
	"github.com/xh3b4sd/tracer"
)

const (
	// limit is the maximum amount of wallet addresses to search for at once.
	limit = 10
)

func (h *Handler) Ensure() error {
	var err error

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

	h.log.Log(
		"level", "info",
		"message", "acquired distributed lock",
		"owner", loc.Owner(),
		"version", strconv.FormatInt(loc.RecordVersionNumber(), 10),
	)

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

	h.log.Log(
		"level", "info",
		"message", "searched pagination cursor",
		"cursor", strconv.FormatInt(cur.Unix(), 10),
	)

	var res accounts.SearchResponse
	{
		res, err = h.acc.Search(cur.Unix(), limit)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	h.log.Log(
		"level", "info",
		"message", "searched account addresses",
		"amount", strconv.Itoa(len(res.Data)),
	)

	// Note that we do not update the cursor in case of an empty search result.
	// That means we are searching for accounts with the very same cursor until we
	// find an actual non-empty result again.

	if len(res.Data) == 0 {
		return nil
	}

	var add []string
	for _, x := range res.Data {
		add = append(add, x.Address)
	}

	// Adding values to the filter is idempotent as per the underlying indexing
	// provider Indexing Co. It is therefore not necessary to de-duplicate the
	// results received from the accounts/search endpoint.

	{
		_, err = h.fil.AddValues(fmt.Sprintf("%s-filter", h.env.Environment), add)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	h.log.Log(
		"level", "info",
		"message", "updated pipeline filter",
		"amount", strconv.Itoa(len(res.Data)),
	)

	// Given the next page, set the next cursor to this next page. If there is no
	// next page, then we are at the end of the line. Given that we have accounts
	// on that last page, set the next cursor to the created timestamp of this
	// very last result object. If there is no next page, and if there is no data,
	// then return early as to not update the next cursor, and keep the current
	// cursor in place for the next iteration. This last case should never happen
	// if the accounts/search API works properly.

	var nxt time.Time
	if res.Next != nil {
		nxt = time.Unix(*res.Next, 0)
	} else if len(res.Data) != 0 {
		nxt = time.Unix(res.Data[len(res.Data)-1].Created, 0)
	} else {
		return nil
	}

	{
		err = h.cur.Update(nxt)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	h.log.Log(
		"level", "info",
		"message", "updated pagination cursor",
		"next", strconv.FormatInt(nxt.Unix(), 10),
		"delta", strconv.FormatInt(nxt.Unix()-cur.Unix(), 10),
	)

	return nil
}

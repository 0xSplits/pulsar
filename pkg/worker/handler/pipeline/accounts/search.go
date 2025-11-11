package accounts

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/xh3b4sd/tracer"
)

type SearchResponse struct {
	Data []SearchResponseObject `json:"data"`
	Next *int64                 `json:"next"`
}

type SearchResponseObject struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Created int64  `json:"created"`
}

func (a *Accounts) Search(cur int64, lim int64) (SearchResponse, error) {
	var err error

	var qry url.Values
	{
		qry.Set("cursor", strconv.FormatInt(cur, 10))
		qry.Set("limit", strconv.FormatInt(lim, 10))
	}

	var pat string
	{
		pat = a.cli.Endpoint() + "/public/v1/accounts/search?" + qry.Encode()
	}

	var req *http.Request
	{
		req, err = http.NewRequest(http.MethodGet, pat, nil)
		if err != nil {
			return SearchResponse{}, tracer.Mask(err)
		}
	}

	var res SearchResponse
	{
		err = a.cli.Request(req, &res)
		if err != nil {
			return SearchResponse{}, tracer.Mask(err)
		}
	}

	return res, nil
}

package client

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

func IsStatusCode(err error) bool {
	return errors.Is(err, statusCodeError)
}

var statusCodeError = &tracer.Error{
	Description: "The request expects a response with a status code of the 200 family. The response status code was >= 300. Therefore the request failed.",
}

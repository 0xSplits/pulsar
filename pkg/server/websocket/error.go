package websocket

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

func IsInvalidWebsocketSecret(err error) bool {
	return errors.Is(err, invalidWebsocketSecretError)
}

var invalidWebsocketSecretError = &tracer.Error{
	Description: "The request expects a valid websocket secret to be provided with the Authorization header. The provided header value did not match the expected secret. Therefore the request failed.",
}

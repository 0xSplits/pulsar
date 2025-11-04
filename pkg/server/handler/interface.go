package handler

import (
	"github.com/gorilla/mux"
)

type Interface interface {
	Attach(*mux.Router, ...interface{})
}

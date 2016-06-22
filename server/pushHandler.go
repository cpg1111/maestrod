package server

import (
	"net/http"
)

type PushSubHandler struct {
	SubHandler
	Error ErrorHandler
}

func (p PushSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

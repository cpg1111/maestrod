package server

import (
	"net/http"
)

type ProjectSubHandler struct {
	SubHandler
	Error ErrorHandler
}

func (p ProjectSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

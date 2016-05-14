package server

import (
	"net/http"
)

type IndexSubHandler struct {
	SubHandler
}

func (i IndexSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	return
}

func (i IndexSubHandler) Post(res http.ResponseWriter, req *http.Request) {
	return
}

func (i IndexSubHandler) Put(res http.ResponseWriter, req *http.Request) {
	return
}

func (i IndexSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	return
}

func (i IndexSubHandler) Delete(res http.ResponseWriter, req *http.Request) {
	return
}

func NewIndexHandler() *RouteHandler {
	return &RouteHandler{
		Route: "/",
		sub:   IndexSubHandler{},
	}
}

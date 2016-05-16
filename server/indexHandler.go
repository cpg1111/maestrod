package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// IndexSubHandler is the SubHandler for the index route
type IndexSubHandler struct {
	SubHandler
}

// Get is the get method for /
func (i IndexSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	routeResp := RouteResp{Routes: []string{"/push"}}
	log.Println(routeResp)
	encErr := json.NewEncoder(res).Encode(routeResp)
	if encErr != nil {
		log.Fatal(encErr)
	}
}

// Post is the post method for /
func (i IndexSubHandler) Post(res http.ResponseWriter, req *http.Request) {
	return
}

// Put is the put method for /
func (i IndexSubHandler) Put(res http.ResponseWriter, req *http.Request) {
	return
}

// Patch is the Patch method for /
func (i IndexSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	return
}

// Delete is the Delete method for /
func (i IndexSubHandler) Delete(res http.ResponseWriter, req *http.Request) {
	return
}

// NewIndexHandler returns a pointer to a new IndexHandler instance
func NewIndexHandler() *RouteHandler {
	return &RouteHandler{
		Route: "/",
		sub:   IndexSubHandler{},
	}
}

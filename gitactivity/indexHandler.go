/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gitactivity

import (
	"encoding/json"
	"log"
	"net/http"
)

// IndexSubHandler is the SubHandler for the index route
type IndexSubHandler struct {
	SubHandler
	Error ErrorHandler
}

// Get is the get method for /
func (i IndexSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	routeResp := RouteResp{Routes: []string{"/push", "/project", "/pullrequest"}}
	log.Println(routeResp)
	encErr := json.NewEncoder(res).Encode(routeResp)
	if encErr != nil {
		log.Fatal(encErr)
	}
}

// Post is the post method for /
func (i IndexSubHandler) Post(res http.ResponseWriter, req *http.Request) {
	i.Error.handle404(res, req)
}

// Put is the put method for /
func (i IndexSubHandler) Put(res http.ResponseWriter, req *http.Request) {
	i.Error.handle404(res, req)
}

// Patch is the Patch method for /
func (i IndexSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	i.Error.handle404(res, req)
}

// Delete is the Delete method for /
func (i IndexSubHandler) Delete(res http.ResponseWriter, req *http.Request) {
	i.Error.handle404(res, req)
}

// NewIndexHandler returns a pointer to a new IndexHandler instance
func NewIndexHandler() *RouteHandler {
	return &RouteHandler{
		Route: "/",
		sub:   IndexSubHandler{},
	}
}

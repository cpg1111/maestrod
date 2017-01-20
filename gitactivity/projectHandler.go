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
	"net/http"

	"github.com/cpg1111/maestrod/config"
)

// ProjectSubHandler handles creating and deleting projects
type ProjectSubHandler struct {
	SubHandler
	Error ErrorHandler
}

// Get returns metadata on a project
func (p ProjectSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	queries := req.URL.Query()
	doneChan := make(chan bool)
	if len(queries["name"]) > 0 {
		store.Find(queries["name"][0], func(dbRes []byte, err error) {
			if err != nil {
				p.Error.handle500(res, req, err)
				doneChan <- true
			} else {
				res.WriteHeader(http.StatusOK)
				res.Write(dbRes)
				doneChan <- true
			}
		})
	} else {
		go func() {
			res.WriteHeader(http.StatusBadRequest)
			res.Write(([]byte)("Bad Request"))
			doneChan <- true
		}()
	}
	_ = <-doneChan
}

// Post creates a project's metadata
func (p ProjectSubHandler) Post(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	reqProject := &config.Project{}
	dErr := decoder.Decode(reqProject)
	if dErr != nil {
		p.Error.handle500(res, req, dErr)
	}
	doneChan := make(chan bool)
	store.Save(reqProject.Name, reqProject, func(saveErr error) {
		if saveErr != nil {
			p.Error.handle500(res, req, saveErr)
			doneChan <- true
		} else {
			res.WriteHeader(http.StatusCreated)
			res.Write(([]byte)("{status: 201, message: \"created\"}"))
			doneChan <- true
		}
	})
	_ = <-doneChan
}

// Put edits a project's metadata
func (p ProjectSubHandler) Put(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	reqProject := &config.Project{}
	dErr := decoder.Decode(reqProject)
	if dErr != nil {
		p.Error.handle500(res, req, dErr)
	}
	doneChan := make(chan bool)
	store.Update(reqProject.Name, reqProject, func(saveErr error) {
		if saveErr != nil {
			p.Error.handle500(res, req, saveErr)
			doneChan <- true
		} else {
			res.WriteHeader(http.StatusCreated)
			res.Write(([]byte)("{\"status\": 201, \"message\": \"created\"}"))
			doneChan <- true
		}
	})
	_ = <-doneChan
}

// Patch does nothing
func (p ProjectSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

// Delete removes a projects metadata
func (p ProjectSubHandler) Delete(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	queries := req.URL.Query()
	doneChan := make(chan bool)
	store.Remove(queries["name"][0], func(err error) {
		if err != nil {
			p.Error.handle500(res, req, err)
			doneChan <- true
		} else {
			res.WriteHeader(http.StatusOK)
			res.Write(([]byte)("{\"status\": 200, \"message\": \"deleted\"}"))
			doneChan <- true
		}
	})
	_ = <-doneChan
}

// NewProjectHandler returns a pointer to a RouteHandler for ProjectSubHandler
func NewProjectHandler() *RouteHandler {
	return &RouteHandler{
		Route: "/project",
		sub:   ProjectSubHandler{},
	}
}

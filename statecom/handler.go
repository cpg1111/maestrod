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

package statecom

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/lifecycle"
)

// Handler handles HTTP for statecom
type Handler struct {
	http.Handler
	Store   datastore.Datastore
	Running *lifecycle.Running
}

// NewHandler returns a pointer to a Handler struct
func NewHandler(store *datastore.Datastore, running *lifecycle.Running) *Handler {
	return &Handler{
		Store:   *store,
		Running: running,
	}
}

// ServeHTTP serves HTTP for statecom
func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.Get(res, req)
		break
	case "POST":
		h.Create(res, req)
		break
	case "PUT":
		h.HandleUnsupported(res, req)
		break
	case "DELETE":
		h.HandleUnsupported(res, req)
		break
	}
}

// Get gets the state of a project or projects
func (h *Handler) Get(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	project := query.Get("project")
	if len(project) == 0 {
		h.getAll(res, req)
	} else {
		h.getOne(res, req, query, project)
	}
}

func (h Handler) getAll(res http.ResponseWriter, req *http.Request) {
	resChan := make(chan []byte)
	errChan := make(chan error)
	h.Store.Find("/state", func(val []byte, err error) {
		if err != nil {
			errChan <- err
			return
		}
		resChan <- val
	})
	for {
		select {
		case dataErr := <-errChan:
			res.WriteHeader(500)
			res.Write([]byte("500 Internal Error"))
			log.Println("WARNING:", dataErr.Error())
			return
		case dataRes := <-resChan:
			res.WriteHeader(200)
			res.Write(dataRes)
			return
		}
	}
}

func (h Handler) getOne(res http.ResponseWriter, req *http.Request, query url.Values, project string) {
	resChan := make(chan []byte)
	errChan := make(chan error)
	key := fmt.Sprintf("/state/%s", project)
	h.Store.Find(key, func(val []byte, err error) {
		if err != nil {
			errChan <- err
			return
		}
		resChan <- val
	})
	for {
		select {
		case dataErr := <-errChan:
			res.WriteHeader(500)
			res.Write([]byte("500 Internal Error"))
			log.Println("WARNING:", dataErr.Error())
			return
		case dataRes := <-resChan:
			res.WriteHeader(200)
			res.Write(dataRes)
			return
		}
	}
}

func projectStateKey(body map[string]interface{}) string {
	return fmt.Sprintf("/states/%s/%s/overall/%s", body["Project"], body["Branch"], body["TimeStamp"])
}

func serviceStateKey(body map[string]interface{}) string {
	state := body["State"].(map[string]interface{})
	return fmt.Sprintf("states/%s/%s/%s/%s", state["Project"], state["Branch"], body["Name"], state["TimeStamp"])
}

// Create creats state of projects
func (h Handler) Create(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var body []byte
	_, readErr := req.Body.Read(body)
	if readErr != nil {
		res.WriteHeader(400)
		res.Write([]byte("400 Bad Request"))
		log.Println(readErr.Error())
		return
	}
	var bodyMap map[string]interface{}
	marshErr := json.Unmarshal(body, bodyMap)
	if marshErr != nil {
		res.WriteHeader(400)
		res.Write([]byte("400 Bad Request"))
		log.Println(marshErr.Error())
	}
	var key string
	if bodyMap["Name"] != nil {
		key = serviceStateKey(bodyMap)
	} else {
		key = projectStateKey(bodyMap)
	}
	errChan := make(chan error)
	h.Store.Save(key, bodyMap, func(err error) {
		errChan <- err
	})
	saveErr := <-errChan
	if saveErr != nil {
		res.WriteHeader(500)
		res.Write([]byte("500 Internal Error"))
		return
	}
	res.WriteHeader(201)
	res.Write([]byte("201 Created"))
	h.Running.CheckIn(bodyMap["Project"].(string), bodyMap["Branch"].(string))
}

// HandleUnsupported handles any unsupported methods
func (h Handler) HandleUnsupported(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	res.WriteHeader(405)
	res.Write([]byte("405 Method Not Allowed"))
}

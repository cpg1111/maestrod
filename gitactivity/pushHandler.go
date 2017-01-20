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
	"strings"

	"github.com/cpg1111/maestrod/lifecycle"
)

// PushSubHandler is the subhandler for handling push hooks
type PushSubHandler struct {
	SubHandler
	Error ErrorHandler
	Queue *lifecycle.Queue
}

// NewPushHandler returns a pointer to a RouteHandler for pushes
func NewPushHandler(queue *lifecycle.Queue) *RouteHandler {
	sub := PushSubHandler{
		Error: ErrorHandler{},
		Queue: queue,
	}
	return &RouteHandler{
		Route: "/push",
		sub:   sub,
	}
}

type authorPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"username"`
}

type commitPayload struct {
	ID        string        `json:"id"`
	TreeID    string        `json:"tree_id"`
	Distinct  bool          `json:"distinct"`
	Message   string        `json:"message"`
	Timestamp string        `json:"timestamp"`
	URL       string        `json:"url"`
	Author    authorPayload `json:"author"`
	Commiter  authorPayload `json:"committer"`
	Added     []string      `json:"added"`
	Removed   []string      `json:"removed"`
	Modified  []string      `json:"modified"`
}

type pushPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Created    bool   `json:"created"`
	Deleted    bool   `json:"deleted"`
	Forced     bool   `json:"forced"`
	BaseRef    string `json:"base_ref"`
	Compare    string
	Commits    []commitPayload `json:"commits"`
	Repository RepoPayload     `json:"repository"`
	Pusher     authorPayload   `json:"pusher"`
}

// Get is PushSubHandler's HTTP GET action
func (p PushSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

// Post is PushSubHandler's HTTP POST action
func (p PushSubHandler) Post(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	payload := &pushPayload{}
	decoder.Decode(payload)
	branchName := strings.Replace(payload.Ref, "refs/heads/", "", -1)
	log.Println("Adding Job to Queue: ", payload.Repository.FullName, branchName, payload.Before, payload.After)
	p.Queue.Add(payload.Repository.FullName, branchName, payload.Before, payload.After)
	resp := PostResp{Status: 201, Message: "Created"}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(resp)
}

// Put is PushSubHandler's HTTP PUT action
func (p PushSubHandler) Put(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

// Patch is PushSubHandler's HTTP PATCH action
func (p PushSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

// Delete is PushSubHandler's HTTP DELETE action
func (p PushSubHandler) Delete(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

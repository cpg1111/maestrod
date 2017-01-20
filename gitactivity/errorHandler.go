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

// ErrorHandler handles HTTP errors
type ErrorHandler struct{}

// ErrorResp is the default response for errors
type ErrorResp struct {
	Error  string `json:"error"`
	Status uint   `json:"status"`
}

func (e *ErrorHandler) handle404(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	resp := ErrorResp{
		Error:  "Page Not Found",
		Status: 404,
	}
	encErr := json.NewEncoder(res).Encode(resp)
	if encErr != nil {
		log.Fatal(encErr)
	}
}

func (e *ErrorHandler) handle500(res http.ResponseWriter, req *http.Request, err error) {
	res.WriteHeader(http.StatusInternalServerError)
	resp := ErrorResp{
		Error:  err.Error(),
		Status: 500,
	}
	encErr := json.NewEncoder(res).Encode(resp)
	if encErr != nil {
		log.Fatal(encErr)
	}
}

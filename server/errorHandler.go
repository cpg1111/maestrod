package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorHandler struct{}

type ErrorResp struct {
	Error  string `json:"error"`
	Status uint   `json:"status"`
}

func (e *ErrorHandler) handle404(res http.ResponseWriter, req *http.Request) {
	resp := ErrorResp{
		Error:  "Page Not Found",
		Status: 404,
	}
	encErr := json.NewEncoder(res).Encode(resp)
	if encErr != nil {
		log.Fatal(encErr)
	}
}

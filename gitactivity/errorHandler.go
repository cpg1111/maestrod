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

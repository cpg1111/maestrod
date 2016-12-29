package statecom

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/lifecycle"
)

type SuccessHandler struct {
	Handler
	queue lifecycle.Queue
}

func NewSuccessHandler(s *datastore.Datastore, r *lifecycle.Running, q *lifecycle.Queue) *SuccessHandler {
	return &SuccessHandler{
		Handler: Handler{
			Store:   *s,
			Running: r,
		},
		queue: *q,
	}
}

type successResp struct {
	Proj   string `json:"project"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}

func (s SuccessHandler) Get(res http.ResponseWriter, req *http.Request) {
	s.HandleUnsupported(res, req)
}

func (s SuccessHandler) Post(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body := &successResp{}
	decErr := json.NewDecoder(req.Body).Decode(body)
	if decErr != nil {
		res.WriteHeader(500)
		res.Write([]byte("500 Internal Error"))
		log.Println("WARNING:", decErr.Error())
		return
	}
	saveErr := s.queue.SaveLastSuccess(body.Proj, body.Branch, body.Commit)
	if saveErr != nil {
		res.WriteHeader(500)
		res.Write([]byte("500 Internal Error"))
		log.Println("WARNING:", saveErr.Error())
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("200 OK\n"))
}

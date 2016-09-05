package statecom

import (
	"net/http"

	"github.com/cpg1111/maestrod/datastore"
)

type Handler struct {
	Store *datastore.Datastore
}

func NewHandler(queue *lifecycle.Queue) *Handler {
	return &Handler{
		Queue: queue,
	}
}

func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.Get(res, req)
		break
	case "POST":
		h.Post(res, req)
		break
	case "PUT":
		h.Put(res, req)
		break
	case "DELETE":
		h.Delete(res, req)
		break
	}
}

func (h *Handler) Get(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	project := query.Get("project")
	if len(project) == 0 {
		h.getAll(res, req)
	} else {
		h.getOne(res, req, query, project)
	}
}

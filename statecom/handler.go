package statecom

import (
	"net/http"
	"net/url"

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

func (h *Handler) getAll(res http.ResponseWriter, req *http.Request) {
	resChan := make(chan []byte)
	errChan := make(chan error)
	h.Store.Find("/configs", func(val []byte, err error) {
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
			return
		case dataRes := <-resChan:
			res.WriteHeader(200)
			res.Write(val)
			return
		}
	}
}

func (h *Handler) getOne(res http.ResponseWriter, req *http.Request, query url.Values, project string) {
	resChan := make(chan []byte)
	errChan := make(chan error)
	key := fmt.Sprintf("/configs/%s", project)
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
			return
		case dataRes := <-resChan:
			res.WriteHeader(200)
			res.Write(val)
			return
		}
	}
}

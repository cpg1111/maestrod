package statecom

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/cpg1111/maestrod/datastore"
)

type Handler struct {
	http.Handler
	Store *datastore.Datastore
}

func NewHandler(store *datastore.Datastore) *Handler {
	return &Handler{
		Store: store,
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
		h.HandleUnsupported(res, req)
		break
	case "DELETE":
		h.HandleUnsupported(res, req)
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

func projectStateKey(body map[string]interface{}) string {
	return fmt.Sprintf("/states/%s/%s/overall/%s", body["Project"], body["Branch"], body["TimeStamp"])
}

func serviceStateKey(body map[string]interface{}) string {
	return fmt.Sprintf("states/%s/%s/%s/%s", body["State"]["Project"], body["State"]["Branch"], body["Name"], body["State"]["TimeStamp"])
}

func (h *Handler) Create(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var body []byte
	_, readErr := req.Body.Read(body)
	if readErr != nil {
		res.WriteHeader(400)
		res.Write("400 Bad Request")
		log.Println(readErr.Error())
		return
	}
	var bodyMap map[string]interface{}
	marshErr := json.Unmarhsal(body, bodyMap)
	if marshErr != nil {
		res.WriteHeader(400)
		res.Write("400 Bad Request")
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
		res.Write("500 Internal Error")
		return
	}
	res.WriteHeader(201)
	res.Write("201 Created")
}

func (h *Handler) HandlerUnsupported(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	res.WriteHeader(405)
	res.Write("405 Method Not Allowed")
}

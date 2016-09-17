package gitactivity

import (
	"encoding/json"
	"net/http"

	"github.com/cpg1111/maestrod/config"
)

type ProjectSubHandler struct {
	SubHandler
	Error ErrorHandler
}

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

func (p ProjectSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

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

func NewProjectHandler() *RouteHandler {
	return &RouteHandler{
		Route: "/project",
		sub:   ProjectSubHandler{},
	}
}

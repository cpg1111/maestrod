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
	"fmt"
	"log"
	"net/http"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/lifecycle"
)

// store is a global datastore
var store datastore.Datastore

// SubHandler is an interface to handle http within a RouteHandler
type SubHandler interface {
	Get(res http.ResponseWriter, req *http.Request)
	Post(res http.ResponseWriter, req *http.Request)
	Put(res http.ResponseWriter, req *http.Request)
	Patch(res http.ResponseWriter, req *http.Request)
	Delete(res http.ResponseWriter, req *http.Request)
}

// RouteHandler is a HTTP Handler for the server
type RouteHandler struct {
	http.Handler
	Route string
	sub   SubHandler
}

// ServeHTTP serves http responses
func (rh RouteHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println("Request: ", req.RemoteAddr, req.Method, req.RequestURI)
	switch req.Method {
	case "GET":
		rh.sub.Get(res, req)
		break
	case "POST":
		rh.sub.Post(res, req)
		break
	case "PUT":
		rh.sub.Put(res, req)
		break
	case "PATCH":
		rh.sub.Patch(res, req)
		break
	case "DELETE":
		rh.sub.Delete(res, req)
		break
	}
}

// Run starts a server
func Run(conf *config.Server, dstore *datastore.Datastore, queue *lifecycle.Queue) (*http.ServeMux, error) {
	store = *dstore
	server := http.NewServeMux()
	indexHandler := NewIndexHandler()
	projectHandler := NewProjectHandler()
	pushHandler := NewPushHandler(queue)
	prHandler := NewPRHandler(queue)
	server.Handle(indexHandler.Route, indexHandler)
	server.Handle(projectHandler.Route, projectHandler)
	server.Handle(pushHandler.Route, pushHandler)
	server.Handle(prHandler.Route, prHandler)
	if conf.RuntimeTLSServer {
		sAddr := fmt.Sprintf("%s:%d", conf.Host, conf.SecurePort)
		iAddr := fmt.Sprintf("%s:%d", conf.Host, conf.InsecurePort)
		log.Println("serving securely at ", sAddr)
		log.Println("redirecting insecure traffic at ", iAddr)
		go func() {
			srvErr := http.ListenAndServeTLS(sAddr, conf.ServerCertPath, conf.ServerKeyPath, server)
			if srvErr != nil {
				panic(srvErr)
			}
		}()
		go redirectInsecure(iAddr, conf.InsecurePort, conf.SecurePort)
	} else {
		addr := fmt.Sprintf("%s:%d", conf.Host, conf.InsecurePort)
		log.Println("serving insecurely at ", addr)
		go func() {
			srvErr := http.ListenAndServe(addr, server)
			if srvErr != nil {
				panic(srvErr)
			}
		}()
	}
	return server, nil
}

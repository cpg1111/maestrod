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

package statecom

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/cpg1111/maestrod/cert"
	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/lifecycle"
)

// Run runs the statecom server
func Run(host, certPath, keyPath string, port int, store *datastore.Datastore, running *lifecycle.Running, queue *lifecycle.Queue) {
	go func() {
		handler := NewHandler(store, running)
		successHandler := NewSuccessHandler(store, running, queue)
		mux := http.NewServeMux()
		mux.Handle("/state", handler)
		mux.Handle("/success", successHandler)
		server := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: mux,
		}
		if len(certPath) > 0 && len(keyPath) > 0 {
			certificate, certErr := cert.GetKeyPair(certPath, keyPath)
			if certErr != nil {
				panic(certErr)
			}
			rootCA, rootCAErr := cert.GetRootCA()
			if rootCAErr != nil {
				panic(rootCAErr)
			}
			tlsConf := &tls.Config{
				Certificates: []tls.Certificate{*certificate},
			}
			if rootCA != nil {
				tlsConf.RootCAs = rootCA
			}
			server.TLSConfig = tlsConf
			server.ListenAndServeTLS(certPath, keyPath)
		} else {
			server.ListenAndServe()
		}
	}()
}

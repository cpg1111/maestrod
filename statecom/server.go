package statecom

import (
	"crypto/tls"
	"net/http"

	"github.com/cpg1111/maestrod/cert"
	"github.com/cpg1111/maestrod/datastore"
)

func Run(host, cert, key string, port int, store *datastore.Datastore) {
	handler := NewHandler(store)
	mux := http.NewServeMux()
	mux.Handler("/state", handler)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}
	if len(cert) > 0 && len(key) > 0 {
		certificate := cert.GetKeyPair(cert, key)
		rootCA := cert.GetRootCA()
		tlsConf := &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}
		if rootCA != nil {
			tlsConf.RootCAs = rootCA
		}
		server.TLSConfig = tlsConf
		server.ListenAndServeTLS()
	} else {
		server.ListenAndServe()
	}
}

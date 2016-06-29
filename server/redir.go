package server

import (
	"fmt"
	"net/http"
	"strings"
)

type redirHandler struct {
	http.Handler
	insecureAddr string
	insecurePort uint
	securePort   uint
}

func (r redirHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL)
	fmt.Println(fmt.Sprintf(":%v", r.insecurePort))
	host := strings.Replace(req.URL.Host, (string)(r.insecurePort), "", 1)
	if host == "" {
		host = strings.Replace(r.insecureAddr, fmt.Sprintf(":%v", r.insecurePort), "", 1)
	}
	fmt.Println(fmt.Sprintf("https://%s:%v%s", host, r.securePort, req.URL.RequestURI()))
	http.Redirect(res, req, fmt.Sprintf("https://%s:%v%s", host, r.securePort, req.URL.RequestURI()), 301)
}

func redirectInsecure(iAddr string, insecurePort, securePort uint) {
	http.ListenAndServe(iAddr, &redirHandler{
		insecureAddr: iAddr,
		insecurePort: insecurePort,
		securePort:   securePort,
	})
}

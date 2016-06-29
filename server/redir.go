package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type redirHandler struct {
	http.Handler
	insecureAddr string
	insecurePort string
	securePort   string
}

func (r redirHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	host := strings.Replace(req.Host, r.insecurePort, r.securePort, 1)
	host = strings.Replace(host, "http", "https", 1)
	http.Redirect(res, req, fmt.Sprintf("https://%s%s", host, req.RequestURI), 301)
}

func redirectInsecure(iAddr string, insecurePort, securePort uint) {
	srvErr := http.ListenAndServe(iAddr, &redirHandler{
		insecureAddr: iAddr,
		insecurePort: strconv.FormatUint((uint64)(insecurePort), 10),
		securePort:   strconv.FormatUint((uint64)(securePort), 10),
	})
	if srvErr != nil {
		panic(srvErr)
	}
}

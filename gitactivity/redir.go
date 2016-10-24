package gitactivity

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// redirHandler handles redirecting to a secure server when available
type redirHandler struct {
	http.Handler
	insecureAddr string
	insecurePort string
	securePort   string
}

// ServeHTTP handles serving http redirects
func (r redirHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	host := strings.Replace(req.Host, r.insecurePort, r.securePort, 1)
	host = strings.Replace(host, "http", "https", 1)
	http.Redirect(res, req, fmt.Sprintf("https://%s%s", host, req.RequestURI), 301)
}

// redirectInsecure does the redirecting
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

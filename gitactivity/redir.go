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

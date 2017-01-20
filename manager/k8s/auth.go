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

package k8s

import (
	"crypto/tls"
	"net/http"

	"github.com/cpg1111/maestrod/cert"
	"github.com/cpg1111/maestrod/config"
)

func loadCerts(conf *config.Server) (*tls.Config, error) {
	certificate, certErr := cert.GetKeyPair(conf.ClientCertPath, conf.ClientKeyPath)
	if certErr != nil {
		return nil, certErr
	}
	certPool, poolErr := cert.GetRootCA()
	if poolErr != nil {
		return nil, poolErr
	}
	if certPool != nil {
		return &tls.Config{
			Certificates: []tls.Certificate{*certificate},
			RootCAs:      certPool,
		}, nil
	}
	return &tls.Config{Certificates: []tls.Certificate{*certificate}}, nil
}

// NewAuthTransport returns a pointer to an http.Transport for authentication
// with k8s or returns an error
func NewAuthTransport(conf *config.Server) (*http.Transport, error) {
	if conf.RuntimeTLSClient {
		tlsConf, tlsErr := loadCerts(conf)
		if tlsErr != nil {
			return nil, tlsErr
		}
		return &http.Transport{TLSClientConfig: tlsConf}, nil
	}
	return &http.Transport{}, nil
}

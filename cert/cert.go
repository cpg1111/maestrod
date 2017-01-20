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

package cert

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
)

// GetKeyPair returns a tls certificate keypair or an error in creating the keypair
func GetKeyPair(cert, key string) (*tls.Certificate, error) {
	certData, cDataErr := ioutil.ReadFile(cert)
	if cDataErr != nil {
		return nil, cDataErr
	}
	keyData, kDataErr := ioutil.ReadFile(key)
	if kDataErr != nil {
		return nil, kDataErr
	}
	certificate, certErr := tls.X509KeyPair(certData, keyData)
	return &certificate, certErr
}

// GetRootCA returns a x509 cert pool or an error creating the cert pool
func GetRootCA() (*x509.CertPool, error) {
	rootCAPath := os.Getenv("ROOT_CA_PATH")
	if len(rootCAPath) > 0 {
		rootCAData, rootCAErr := ioutil.ReadFile(rootCAPath)
		if rootCAErr != nil {
			return nil, rootCAErr
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(rootCAData)
		return certPool, nil
	}
	return nil, nil
}

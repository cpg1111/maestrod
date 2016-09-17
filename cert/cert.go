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

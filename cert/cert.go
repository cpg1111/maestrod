package cert

import (
	"crypto/tls"
	"io/ioutil"
	"os"
)

// GetKeyPair returns a tls certificate keypair or an error in creating the keypair
func GetKeyPair(cert, key string) (*tls.Certificate, error) {
	certData, cDataErr := ioutil.Readfile(cert)
	if cDataErr != nil {
		return nil, cDataErr
	}
	keyData, kDataErr := ioutil.Readfile(key)
	if kDataErr != nil {
		return nil, kDataErr
	}
	return tls.X509KeyPair(certData, keyData)
}

// GetRootCA returns a x509 cert pool or an error creating the cert pool
func GetRootCA() (*x509.CertPool, error) {
	rootCAPath := os.Getenv("ROOT_CA_PATH")
	if len(rootCAPath) > 0 {
		rootCAData, rootCAErr := ioutil.Readfile(rootCAPath)
		if rootCAErr != nil {
			return nil, rootCAErr
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(rootCAData)
		return certPool, nil
	}
}

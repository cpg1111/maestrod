package k8s

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cpg1111/maestrod/config"
)

func loadCerts(conf *config.Server) (*tls.Config, error) {
	certData, cDataErr := ioutil.ReadFile(conf.ClientCertPath)
	if cDataErr != nil {
		return nil, cDataErr
	}
	keyData, kDataErr := ioutil.ReadFile(conf.ClientKeyPath)
	if kDataErr != nil {
		return nil, kDataErr
	}
	cert, certErr := tls.X509KeyPair(certData, keyData)
	if certErr != nil {
		return nil, certErr
	}
	rootCAPath := os.Getenv("ROOT_CA_PATH")
	if len(rootCAPath) > 0 {
		rootCAData, rootCAErr := ioutil.ReadFile(rootCAPath)
		if rootCAErr != nil {
			return nil, rootCAErr
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(rootCAData)
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		}, nil
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

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

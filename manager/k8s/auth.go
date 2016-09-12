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

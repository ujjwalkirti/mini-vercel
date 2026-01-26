package kafka_tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func NewTLSConfig() (*tls.Config, error) {
	path := "ca.pem"
	if os.Getenv("ENV") == "production" {
		path = "/secrets/kafka-consumer-ca"
	}

	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		RootCAs: pool,
	}, nil
}

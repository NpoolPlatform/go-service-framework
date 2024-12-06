package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"google.golang.org/grpc/credentials"
)

func LoadTLSConfig() (credentials.TransportCredentials, error) {
	tlsFilesDir := config.GetIntValueWithNameSpace("", config.KeyGRPCSDir)
	return loadTLSConfig(
		fmt.Sprintf("%v/%v", tlsFilesDir, "server.crt"),
		fmt.Sprintf("%v/%v", tlsFilesDir, "server.key"),
		fmt.Sprintf("%v/%v", tlsFilesDir, "ca.crt"),
	)
}

func loadTLSConfig(certFile, keyFile, caFile string) (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certification: %w", err)
	}

	data, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("faild to read CA certificate: %w", err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("unable to append the CA certificate to CA pool")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig), nil
}

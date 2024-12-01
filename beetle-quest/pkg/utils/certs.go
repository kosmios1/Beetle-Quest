package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	caCertPath = "/certs/caCert.pem"
	caKeyPath  = "/certs/caKey.pem"

	serverCertPath = "/serverCert.pem"
	serverKeyPath  = "/serverKey.pem"
)

func readPEMFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}
	return block.Bytes, nil
}

func savePEMFile(filename string, typeStr string, bytes []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(file, &pem.Block{Type: typeStr, Bytes: bytes})
}

func readPEMFileWithoutDecoding(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func loadCACertAndPrvKey() (*x509.Certificate, *rsa.PrivateKey) {
	caCertBytes, err := readPEMFile(caCertPath)
	if err != nil {
		log.Panicf("Failed to read caCert.pem: %s", err)
		return nil, nil
	}
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		log.Panicf("Failed to parse caCert.pem: %s", err)
		return nil, nil
	}

	caPrvKeyBytes, err := readPEMFile(caKeyPath)
	if err != nil {
		log.Panicf("Failed to read caKey.pem: %s", err)
		return nil, nil
	}
	caPrvKey, err := x509.ParsePKCS8PrivateKey(caPrvKeyBytes)
	if err != nil {
		log.Panicf("Failed to parse caKey.pem: %s", err)
		return nil, nil
	}
	return caCert, caPrvKey.(*rsa.PrivateKey)
}

func GenOwnCertAndKey(serviceName string) {
	caCert, caPrvKey := loadCACertAndPrvKey()
	servicePrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicf("Failed to generate private key: %s", err)
		return
	}

	serviceCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization:  []string{"Beetle Quest"},
			Country:       []string{"IT"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    serviceName,
		},

		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// For simplicity, we use the same DNS names for all services
		DNSNames:              []string{"reverse-proxy", "localhost", "admin-service", "auth-service", "user-service", "gacha-service", "market-service", "static-service"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	serviceCertBytes, err := x509.CreateCertificate(rand.Reader, serviceCertTemplate, caCert, &servicePrivKey.PublicKey, caPrvKey)
	if err != nil {
		log.Panicf("Failed to create certificate: %s", err)
		return
	}

	savePEMFile(serverCertPath, "CERTIFICATE", serviceCertBytes)
	savePEMFile(serverKeyPath, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(servicePrivKey))

	log.Println("[INFO] Generated client certificate and key")
}

func getTlsConfig(isServerConfig bool) *tls.Config {
	caCert, err := readPEMFileWithoutDecoding(caCertPath)
	if err != nil {
		log.Panicf("Could not read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("Failed to append CA certificate to pool")
	}

	ownCert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		log.Panicf("Could not load server certificate: %v", err)
	}

	if isServerConfig {
		return &tls.Config{
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{ownCert},
		}
	} else {
		return &tls.Config{
			// TODO: Without insecureSkipVerify = true this error occurs:
			// failed to verify certificate: x509: certificate relies on legacy Common Name field, use SANs instead
			InsecureSkipVerify: true,
			RootCAs:            caCertPool,
			Certificates:       []tls.Certificate{ownCert},
		}
	}
}

func SetupHTPPSServer(h http.Handler) *http.Server {
	tlsConfig := getTlsConfig(true)
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   h,
	}
	return server
}

func SetupHTTPSClient() *http.Client {
	tlsConfig := getTlsConfig(false)

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	return client
}

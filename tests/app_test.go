package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	pb "movies/pkg/api"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	caCertFilePath = "../cert/ca-cert.crt"     // нужен удостовериться в сертификате сервера
	certFilePath   = "../cert/client-cert.crt" // этот сертификат передается серверу
	keyFilePath    = "../cert/client-key.key"  // секретный ключ
)

func CreateClientTLSConfig(caCertFilePath, certFilePath, keyFilePath string) (*tls.Config, error) {
	// Create a pool with the server certificate since it is not signed by a known CA
	caCert, err := os.ReadFile(caCertFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading server certificate: %s", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil, err
	}

	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert}, // <-- provide client cert
		RootCAs:            caCertPool,                    // aka: curl -v --cacert ./cert/server-cert.crt https://127.0.0.1:8443/hello
		InsecureSkipVerify: false,                         // aka: curl -sL https://127.0.0.1:8443/hello --insecure
	}

	return tlsConfig, nil
}

// /usr/local/go/bin/go test -v -timeout 30s -run ^TestApp$ movies/tests
func TestGetMovies(t *testing.T) {
	tlsConfig, err := CreateClientTLSConfig(caCertFilePath, certFilePath, keyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("127.0.0.1:8888",
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	defer conn.Close()

	cli := pb.NewMoviesClient(conn)

	resp, err := cli.GetMovies(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("GetMovies error: %v", err)
	} else {
		t.Log("movies list is ", resp.Movies)
	}
}

// /usr/local/go/bin/go test -v -timeout 30s -run ^TestMovieStream$ movies/tests
func TestMovieStream(t *testing.T) {
	// grpc.WithTransportCredentials(insecure.NewCredentials()
	tlsConfig, err := CreateClientTLSConfig(caCertFilePath, certFilePath, keyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("127.0.0.1:8888",
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := pb.NewMoviesClient(conn)

	stream, err := cli.GetMovie(context.Background(), &pb.MovieRequest{
		MovieName: "Test",
	})
	if err != nil {
		assert.Equal(t, err, nil)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				assert.Equal(t, nil, err, err.Error())
			}
			return
		}

		log.Println(string(resp.Bytes))
	}
}

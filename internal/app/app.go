package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	API "movies/internal/app/api"
	"movies/internal/middleware"
	"movies/internal/middleware/metrics"

	pb "movies/pkg/api"
	"movies/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	port                 = ":8443"
	certFilePath         = "./cert/server-cert.crt"
	keyFilePath          = "./cert/server-key.key"
	caCertFilePath       = "./cert/ca-cert.crt"
	caClientCertFilePath = "./cert/client-cert.crt"
)

func StartMovieApp() {
	// log.SetFlags(log.Lshortfile)

	if os.Getenv("DEBUG") == "true" {
		err := godotenv.Load("./configs/.env")
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	// Logger
	logger.InitLogger()

	// API
	api, err := API.InitAPI()
	if err != nil {
		logger.Global.Fatal(err.Error())
	}

	log.Println("API OK")

	// Debug Sever
	mux := http.NewServeMux()
	// k8s readiness probe
	// mux.HandleFunc("/healthz", srv.healthcheckHandler)
	// prometheus metrics
	mux.Handle("/metrics", promhttp.Handler())
	// pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	httpServer := &http.Server{
		Addr:    os.Getenv("DEBUG_PORT"),
		Handler: mux,
	}

	// Middlewares
	ChainUnaryInterceptors := []grpc.UnaryServerInterceptor{
		middleware.RecoverUnaryInterceptor(),
		middleware.LimiterUnaryInterceptor(),
		middleware.LogErrorUnaryInterceptor(),
		metrics.MetricsUnaryInterceptor(),
	}

	StreamInterceptors := []grpc.StreamServerInterceptor{
		middleware.RecoverStreamInterceptor(),
		middleware.LimiterStreamInterceptor(),
		middleware.RecoverStreamInterceptor(),
	}

	// Options
	tlsConfig, err := createServerTLSConfig(caCertFilePath, certFilePath, keyFilePath)
	if err != nil {
		logger.Global.Fatal(fmt.Sprintf("failed create TLS config: %v", err))
	}

	grpcServerOptions := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}

	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(ChainUnaryInterceptors...),
		grpc.ChainStreamInterceptor(StreamInterceptors...),
	)

	// GRPC server
	grpcServer := grpc.NewServer(
		grpcServerOptions...,
	)
	pb.RegisterMoviesServer(grpcServer, api)

	reflection.Register(grpcServer)

	lisGRPC, err := net.Listen("tcp", os.Getenv("PORT"))
	if err != nil {
		logger.Global.Fatal(fmt.Sprintf("run: %v", err))
	}

	// Serve and gracefull shutdown
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return grpcServer.Serve(lisGRPC)
	})
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		grpcServer.GracefulStop()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil && err != http.ErrServerClosed {
		logger.Global.Error(fmt.Sprintf("exit reason: %s \n", err))
	}
}

func createServerTLSConfig(caCertFilePath, certFile, keyFile string) (*tls.Config, error) {
	// Load server's certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509: %v", err)
	}

	// Load certificate of the CA who signed client's certificate
	rootCA, err := os.ReadFile(caCertFilePath)
	if err != nil {
		return nil, err
	}

	// Load certificate of the CA who signed client's certificate
	clientCA, err := os.ReadFile(caClientCertFilePath)
	if err != nil {
		return nil, err
	}

	clientCAs := x509.NewCertPool()
	if !clientCAs.AppendCertsFromPEM(rootCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	if !clientCAs.AppendCertsFromPEM(clientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	// RequireAndVerifyClientCert -
	// RequireAnyClientCert -
	// RequestClientCert -
	// NoClientCert -

	// Create tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // <--
		ClientCAs:    clientCAs,                      // <--
	}

	return tlsConfig, nil
}

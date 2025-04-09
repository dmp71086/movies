package api

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"

	"movies/internal/app/adapters/pg"

	// pbAuth "movies/pkg/api/auth"
	pb "movies/pkg/api"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
)

// const (
// 	panScoring      = "panScoring"
// 	prometheus_port = 9049
// )

type IApi interface {
	GetMovies(ctx context.Context, in *pb.EmptyRequest) (*pb.MoviesResponse, error)
	GetMovie(req *pb.MovieRequest, stream grpc.ServerStreamingServer[pb.Movie]) error
	BuyMovie(context.Context, *pb.MovieRequest) (*pb.EmptyResponse, error)
}

type API struct {
	pb.UnimplementedMoviesServer
	repo pg.IPG
	// AuthClient pbAuth.AuthClient
	PublicKey *rsa.PublicKey

	validator *protovalidate.Validator
}

func InitAPI() (*API, error) {
	pg, err := pg.InitRepository()
	if err != nil {
		return nil, err
	}
	log.Println("Repo OK")

	// validator
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
			&pb.MovieRequest{},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
	}

	//prometheus.New(panScoring)
	//go prometheus.Run(prometheus_port)
	// log.Println("prometheus OK")

	// nspkService, err := clients.NewNSPKService()
	// if err != nil {
	// 	return nil, err
	// }

	// log.Println("nspkService OK")

	return &API{
		repo:      pg,
		validator: validator,
	}, nil
}

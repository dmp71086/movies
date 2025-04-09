package api

import (
	"crypto/rand"
	"log"

	pb "movies/pkg/api"
	logger "movies/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *API) GetMovie(req *pb.MovieRequest, stream grpc.ServerStreamingServer[pb.Movie]) error {
	if err := s.validator.Validate(req); err != nil {
		return RPCValidationError(err)
	}

	log.Println("ListNotesStream: received")

	ch := make(chan *pb.Movie, 100)
	go func() {
		for range [4]int{} {
			film := make([]byte, 44)
			_, _ = rand.Read(film)
			ch <- &pb.Movie{
				Bytes: film,
			}
		}
		close(ch)
	}()

	for chunk := range ch {
		if err := stream.Send(chunk); err != nil {
			logger.Global.Error(err.Error(),
				zap.String("Method", "GetMovie"))
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

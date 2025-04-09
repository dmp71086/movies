package api

import (
	"context"

	pb "movies/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *API) GetMovies(ctx context.Context, in *pb.EmptyRequest) (*pb.MoviesResponse, error) {
	moviesDTO, err := s.repo.GetMovies(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//fmt.Println(movies.Movies)

	var movies pb.MoviesResponse
	for _, m := range moviesDTO {
		movies.Movies = append(movies.Movies, &pb.MovieDescription{
			Name:        m.Name,
			Description: m.Description,
		})
	}

	return &movies, nil
}

/* func (s *MoviesService) BuyMovie(msg *stan.Msg) {
	//fmt.Println(ctx.Value("email"))

	err := s.MoviesRepository.BuyMovie(context.Background(), "test")
	if err != nil {
		fmt.Print("Error in repository: ")
		fmt.Println(err.Error())
	}

	for i := 0; i < 3; i++ {
		err = msg.Ack()
		if err != nil {
			fmt.Println(err.Error(), whereami.WhereAmI())

			time.Sleep(100 * time.Microsecond)
			continue
		}

		prometheus.CountSuccessActionWithDescription("Количество загруженных объектов", "ChrtMsgCounter", len(list))
		prometheus.CountSuccessActionWithDescription("Количество успешных запросов", "SuccessAck", 1)
		return
	}
} */

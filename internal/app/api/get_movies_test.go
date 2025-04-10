package api

import (
	"errors"
	"testing"


	pb "movies/pkg/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/bufbuild/protovalidate-go"
)

// Мок стрима
type mockMovieStream struct {
	mock.Mock
	grpc.ServerStream
}

func (m *mockMovieStream) Send(movie *pb.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}

// Мок валидатора
type mockValidator struct {
	mock.Mock
}

func (v *mockValidator) Validate(req any) error {
	args := v.Called(req)
	return args.Error(0)
}

func TestGetMovie_Success(t *testing.T) {
	stream := new(mockMovieStream)

	validator, _ := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
			&pb.MovieRequest{},
		),
	)

	req := &pb.MovieRequest{}

	// stream.Send ожидается 4 раза
	stream.On("Send", mock.AnythingOfType("*pb.Movie")).Return(nil).Times(4)

	s := &API{
		validator: validator,
	}

	err := s.GetMovie(req, stream)
	assert.NoError(t, err)

	validator.AssertExpectations(t)
	stream.AssertExpectations(t)
}

func TestGetMovie_ValidationError(t *testing.T) {
	stream := new(mockMovieStream)

	req := &pb.MovieRequest{}

	validator, _ := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
			&pb.MovieRequest{},
		),
	)


	err := s.GetMovie(req, stream)

	assert.Error(t, err)
	assert.Equal(t, RPCValidationError(validationErr), err)
}

func TestGetMovie_StreamError(t *testing.T) {
	stream := new(mockMovieStream)
	validator, _ := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
			&pb.MovieRequest{},
		),
	)


	req := &pb.MovieRequest{}

	// Первые 2 Send успешные, 3-я выдаёт ошибку
	stream.On("Send", mock.Anything).Return(nil).Once()
	stream.On("Send", mock.Anything).Return(nil).Once()
	stream.On("Send", mock.Anything).Return(errors.New("stream error")).Once()

	s := &API{
		validator: validator,
	}

	err := s.GetMovie(req, stream)

	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))

	stream.AssertExpectations(t)
}

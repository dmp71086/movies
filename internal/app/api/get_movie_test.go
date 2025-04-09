package api

import (
	"context"
	"errors"
	"movies/internal/app/models"
	pb "movies/pkg/api"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetMovies(t *testing.T) {
	api := API{
		repo: &RepositoryMock{},
	}

	type testCase struct {
		n        int
		expected error
	}

	testCases := []testCase{
		{
			n:        0,
			expected: nil,
		},
		{
			n:        1,
			expected: status.Error(codes.Internal, ""),
		},
	}

	for _, v := range testCases {
		er = v.n

		_, err := api.GetMovies(context.Background(), &pb.EmptyRequest{})

		assert.Equal(t, err, v.expected)
	}
}

type RepositoryMock struct{}

var er int

func (s *RepositoryMock) GetMovies(ctx context.Context) ([]models.Movie, error) {
	if er == 0 {
		return []models.Movie{{Name: "test_movie"}}, nil
	} else {
		return nil, errors.New("some error")
	}
}

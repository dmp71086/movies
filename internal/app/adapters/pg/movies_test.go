package pg

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

type MockPool struct{}

func (m MockPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if er == 0 {
		return pgxmock.NewRows([]string{"name", "description"}).
			AddRow("test_movie1", "test_description1").AddRow("test_movie2", "test_description2").Kind(), nil
	} else {
		return nil, errors.New("pg select error")
	}

}

func (m MockPool) Ping(ctx context.Context) error {
	return nil
}

var er int

// /usr/local/go/bin/go test -v -timeout 30s -run ^TestGetMovies$ movies/internal/app/adapters/pg
func TestGetMovies(t *testing.T) {
	repo := &Postgres{
		pool: MockPool{},
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
			expected: errors.New("failed to select movies: pg select error"),
		},
	}

	for _, v := range testCases {
		er = v.n

		movies, err := repo.GetMovies(context.Background())
		t.Log(movies)

		assert.Equal(t, err, v.expected)
	}
}

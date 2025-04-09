package pg

import (
	"context"
	"fmt"
	"time"

	"movies/internal/app/models"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) GetMovies(ctx context.Context) ([]models.Movie, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := pg.pool.Query(ctxWithTimeout, `SELECT name, description, path FROM movie_service.movies`)
	if err != nil {
		return nil, fmt.Errorf("failed to select movies: %v", err.Error())
	}
	defer rows.Close()

	moviesDTO, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.MovieDTO])
	if err != nil {
		return nil, fmt.Errorf("failed to scan partner: %v", err.Error())
	}

	movies := make([]models.Movie, len(moviesDTO))

	for i:= 0; i < len(moviesDTO); i++ {
		movies[i] = models.FromDTO(&moviesDTO[i])
	}

	fmt.Println(movies)

	return movies, nil
}

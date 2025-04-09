package pg

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"movies/internal/app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	//prometheus "gitlab.equifax.local/libs/lib-prometheus-monitoring.git/monitoring"
)

const (
	maxConnIdleTimeDefault     = time.Minute
	maxConnLifeTimeDefault     = time.Hour
	minConnectionsCountDefault = 2
	maxConnectionsCountDefault = 10
)

type IPG interface {
	GetMovies(ctx context.Context) ([]models.Movie, error)
}

type IPgPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Ping(ctx context.Context) error
}

type Postgres struct {
	pool IPgPool
	// key string
	sync.RWMutex
}

func InitRepository() (IPG, error) {
	dbPool := initPostgres()

	repo := &Postgres{
		pool: dbPool,
	}

	go repo.checkConnect()

	return repo, nil
}

func (r *Postgres) checkConnect() {
	var count int
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		err := r.pool.Ping(ctx)
		cancel()
		if err == nil {
			count = 0
			time.Sleep(3 * time.Second)
			continue
		}

		count++
		if count > 1 {
			conn := initPostgres()

			r.Lock()
			r.pool = conn
			r.Unlock()
			
			count = 0
		}
	}
}

func initPostgres() *pgxpool.Pool {
	host := os.Getenv("HOST")
	db := os.Getenv("DB")
	password := os.Getenv("PASSWORD")
	user_db := os.Getenv("USER_DB")

	fmt.Printf("postgres://%s:%s@%s:%d/%s\n", user_db, password, host, 5432, db)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx,
		// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user_db, password, host, 5432, db))
	if err != nil {
		//prometheus.AddError(models.Fatal)
		log.Fatal("can't create connection pool: " + err.Error())
	}

	dbpool.Config().MaxConnIdleTime = maxConnIdleTimeDefault
	dbpool.Config().MaxConnLifetime = maxConnLifeTimeDefault
	dbpool.Config().MaxConns = maxConnectionsCountDefault
	dbpool.Config().MinConns = minConnectionsCountDefault

	// dbpool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	// 	pgxUUID.Register(conn.TypeMap())
	// 	return nil
	// }

	err = dbpool.Ping(context.Background())
	if err != nil {
		//prometheus.AddError(models.Fatal)
		log.Fatal("can't create connection pool: " + err.Error())
	}

	return dbpool
}

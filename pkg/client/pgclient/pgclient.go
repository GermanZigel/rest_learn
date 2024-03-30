package pgclient

import (
	"context"
	"fmt"
	"rest/internal/config"
	"rest/internal/logging"
	"rest/pkg/utils"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PoolClient представляет интерфейс для работы с пулом подключений к базе данных.
type PoolClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// PgxPoolClient является оболочкой для *pgxpool.Pool, реализующей интерфейс PoolClient.
type PgxPoolClient struct {
	*pgxpool.Pool
}

// NewClient создает новый клиент для работы с пулом подключений к базе данных.
func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (PoolClient, error) {
	logger := logging.GetLogger()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)

	var pool *pgxpool.Pool
	var err error

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		logger.Fatal("error do with tries postgresql")
		return nil, err
	}

	return &PgxPoolClient{pool}, nil
}

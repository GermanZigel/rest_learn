package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"rest/internal/logging"
	"rest/internal/userProxy"
	"rest/pkg/client/pgclient"

	"github.com/jackc/pgconn"
)

type repository struct {
	client pgclient.PoolClient
	logger *logging.Logger
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) *repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, user userProxy.User) (string, error) {
	var id string
	q := "INSERT INTO public.users\n(id, \"Name\", job, created)\nVALUES($1, $2, $3, $4) returning id"
	row := r.client.QueryRow(ctx, q, user.Id, user.Name, user.Job, user.Created)
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr)
			return "", newErr
		}
	}
	return id, nil
}

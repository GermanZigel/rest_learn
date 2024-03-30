package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"rest/internal/logging"
	"rest/internal/userProxy"
	"rest/pkg/client/pgclient"
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
	var Id string
	q := "INSERT INTO public.users\n(id, \"Name\", job, created)\nVALUES($1, $2, $3, $4) returning id"
	row := r.client.QueryRow(ctx, q, user.Id, user.Name, user.Job, user.Created)
	err := row.Scan(&Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr)
			return "", newErr
		}
	}
	return Id, nil
}
func (r *repository) GetList(ctx context.Context) ([]userProxy.User, error) {
	logger := logging.GetLogger()
	q := "select id, users.\"Name\", job, created from users order by created desc"
	rows, err := r.client.Query(ctx, q)
	logger.Info("result query", &rows)
	if err != nil {
		return nil, err
	}
	users := make([]userProxy.User, 0)
	for rows.Next() {
		var usr userProxy.User
		err = rows.Scan(&usr.Id, &usr.Name, &usr.Job, &usr.Created)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
		logger.Infof("Error after append %s", err)

	}
	logger.Infof("User recieving from database %s", users)
	return users, nil
}

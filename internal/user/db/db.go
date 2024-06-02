package db

import (
	"context"
	"errors"
	"fmt"
	"rest/internal/logging"
	"rest/internal/userProxy"
	"rest/pkg/client/pgclient"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	ctx    context.Context
	client pgclient.PoolClient
	logger *logging.Logger
}

func NewRepository(ctx context.Context, client *pgxpool.Pool, logger *logging.Logger) *repository {
	return &repository{
		ctx:    ctx, // Добавление контекста
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
func (r *repository) GetOnce(ctx context.Context, id int) (userProxy.User, error) {
	logger := logging.GetLogger()
	q := "select id, users.\"Name\", job, created  from users where id in ($1)"
	logger.Infof("SQL Query: %s, id=%d", formatQuery(q), id)

	var usr userProxy.User
	err := r.client.QueryRow(ctx, q, id).Scan(&usr.Id, &usr.Name, &usr.Job, &usr.Created) // Исправлено
	if err != nil {
		if err == pgx.ErrNoRows {
			return userProxy.User{}, pgx.ErrNoRows
			logger.Infof("Ошибка = %s", pgx.ErrNoRows)
		} else {
			return userProxy.User{}, err
		}

	}

	return usr, nil
}

func (r *repository) DeleteOnce(ctx context.Context, id int) (bool, error) {
	logger := logging.GetLogger()
	q := "delete from users where id = $1"

	tx, err := r.client.Begin(ctx)
	if err != nil {
		logger.Printf("Failed to begin transaction: %v\n", err)
		return false, err
	}
	logger.Info("Transaction  DeleteOnce begin")
	logger.Infof("SQL Query: %s, id=%d", formatQuery(q), id)
	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			logger.Printf("Transaction rollback failed: %v\n", err)
			return false, err
		}
		return false, err
	} else {
		err := tx.Commit(ctx)
		if err != nil {
			logger.Printf("Transaction commit failed: %v\n", err)
		}
		logger.Info("Transaction  DeleteOnce commit")
		return true, nil
	}
}
func (r *repository) Update(ctx context.Context, u userProxy.User) (userProxy.User, error) {
	logger := logging.GetLogger()
	var updatedUser userProxy.User
	q := "update users set  \"Name\" = $2, job= $3 where id = $1  returning id, job,\"Name\""
	row := r.client.QueryRow(ctx, q, u.Id, u.Name, u.Job)
	logger.Infof("SQL Query: %s, id=%d", formatQuery(q), u.Id)
	err := row.Scan(&updatedUser.Id, &updatedUser.Job, &updatedUser.Name)
	logger.Infof("query results", err)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr)
			return updatedUser, newErr
		}
	}
	return updatedUser, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

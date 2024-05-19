package storage

import (
	"context"
	"rest/internal/userProxy"
)

type Repository interface {
	Create(ctx context.Context, user userProxy.User) (string, error)
	GetList(ctx context.Context) (u []userProxy.User, err error)
	GetOnce(ctx context.Context, id int) (userProxy.User, error)
	Update(context.Context, userProxy.User) (userProxy.User, error)
	DeleteOnce(ctx context.Context, id int) (bool, error)
}

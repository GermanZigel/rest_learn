package storage

import (
	"context"
	"rest/internal/userProxy"
)

type Repository interface {
	Create(ctx context.Context, author *userProxy.User) error
	FindAll(ctx context.Context) (u []userProxy.User, err error)
	FindOne(ctx context.Context, id string) (userProxy.User, error)
	Update(ctx context.Context, user userProxy.User) error
	Delete(ctx context.Context, id string) error
}

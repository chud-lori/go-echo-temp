package ports

import (
	"context"

	"github.com/chud-lori/go-echo-temp/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (*entities.User, error)
	FindAll(ctx context.Context) ([]*entities.User, error)
}

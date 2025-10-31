package repository

import (
	"context"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

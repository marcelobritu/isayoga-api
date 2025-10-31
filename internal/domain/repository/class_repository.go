package repository

import (
	"context"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassRepository interface {
	Create(ctx context.Context, class *entity.Class) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Class, error)
	FindAll(ctx context.Context) ([]*entity.Class, error)
	Update(ctx context.Context, class *entity.Class) error
	IncrementEnrollmentWithVersion(ctx context.Context, classID primitive.ObjectID, currentVersion int) error
	DecrementEnrollment(ctx context.Context, classID primitive.ObjectID) error
	WithTransaction(ctx context.Context, fn func(context.Context, mongo.SessionContext) error) error
}


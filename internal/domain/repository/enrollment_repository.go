package repository

import (
	"context"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *entity.Enrollment) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Enrollment, error)
	FindByUserAndClass(ctx context.Context, userID, classID primitive.ObjectID) (*entity.Enrollment, error)
	FindByUser(ctx context.Context, userID primitive.ObjectID) ([]*entity.Enrollment, error)
	Update(ctx context.Context, enrollment *entity.Enrollment) error
}


package repository

import (
	"context"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Payment, error)
	FindByEnrollmentID(ctx context.Context, enrollmentID primitive.ObjectID) (*entity.Payment, error)
	FindByMercadoPagoID(ctx context.Context, mpID string) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
}

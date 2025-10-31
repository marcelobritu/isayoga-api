package mongodb

import (
	"context"
	"fmt"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(db *mongo.Database) *PaymentRepository {
	return &PaymentRepository{
		collection: db.Collection("payments"),
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	_, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return fmt.Errorf("erro ao criar pagamento: %w", err)
	}
	return nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("pagamento não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar pagamento: %w", err)
	}
	return &payment, nil
}

func (r *PaymentRepository) FindByEnrollmentID(ctx context.Context, enrollmentID primitive.ObjectID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.collection.FindOne(ctx, bson.M{"enrollment_id": enrollmentID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar pagamento: %w", err)
	}
	return &payment, nil
}

func (r *PaymentRepository) FindByMercadoPagoID(ctx context.Context, mpID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.collection.FindOne(ctx, bson.M{"mercado_pago_id": mpID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("pagamento não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar pagamento: %w", err)
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *entity.Payment) error {
	update := bson.M{
		"$set": payment,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": payment.ID}, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar pagamento: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("pagamento não encontrado")
	}

	return nil
}


package mongodb

import (
	"context"
	"fmt"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnrollmentRepository struct {
	collection *mongo.Collection
}

func NewEnrollmentRepository(db *mongo.Database) *EnrollmentRepository {
	return &EnrollmentRepository{
		collection: db.Collection("enrollments"),
	}
}

func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *entity.Enrollment) error {
	_, err := r.collection.InsertOne(ctx, enrollment)
	if err != nil {
		return fmt.Errorf("erro ao criar inscrição: %w", err)
	}
	return nil
}

func (r *EnrollmentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Enrollment, error) {
	var enrollment entity.Enrollment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&enrollment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("inscrição não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar inscrição: %w", err)
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) FindByUserAndClass(ctx context.Context, userID, classID primitive.ObjectID) (*entity.Enrollment, error) {
	var enrollment entity.Enrollment
	err := r.collection.FindOne(ctx, bson.M{
		"user_id": userID,
		"class_id": classID,
		"status": bson.M{"$in": []string{"pending", "confirmed"}},
	}).Decode(&enrollment)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar inscrição: %w", err)
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) FindByUser(ctx context.Context, userID primitive.ObjectID) ([]*entity.Enrollment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar inscrições: %w", err)
	}
	defer cursor.Close(ctx)

	var enrollments []*entity.Enrollment
	if err = cursor.All(ctx, &enrollments); err != nil {
		return nil, fmt.Errorf("erro ao processar inscrições: %w", err)
	}

	if enrollments == nil {
		enrollments = []*entity.Enrollment{}
	}

	return enrollments, nil
}

func (r *EnrollmentRepository) Update(ctx context.Context, enrollment *entity.Enrollment) error {
	update := bson.M{
		"$set": enrollment,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": enrollment.ID}, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar inscrição: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("inscrição não encontrada")
	}

	return nil
}


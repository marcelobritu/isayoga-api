package mongodb

import (
	"context"
	"fmt"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuário: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*entity.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("erro ao processar usuários: %w", err)
	}

	if users == nil {
		users = []*entity.User{}
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	return nil
}

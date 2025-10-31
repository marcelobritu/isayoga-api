package mongodb

import (
	"context"
	"fmt"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func NewClassRepository(db *mongo.Database, client *mongo.Client) *ClassRepository {
	return &ClassRepository{
		collection: db.Collection("classes"),
		client:     client,
	}
}

func (r *ClassRepository) Create(ctx context.Context, class *entity.Class) error {
	_, err := r.collection.InsertOne(ctx, class)
	if err != nil {
		return fmt.Errorf("erro ao inserir aula: %w", err)
	}
	return nil
}

func (r *ClassRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Class, error) {
	var class entity.Class
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&class)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("aula não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar aula: %w", err)
	}
	return &class, nil
}

func (r *ClassRepository) FindAll(ctx context.Context) ([]*entity.Class, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar aulas: %w", err)
	}
	defer cursor.Close(ctx)

	var classes []*entity.Class
	if err = cursor.All(ctx, &classes); err != nil {
		return nil, fmt.Errorf("erro ao processar aulas: %w", err)
	}

	if classes == nil {
		classes = []*entity.Class{}
	}

	return classes, nil
}

func (r *ClassRepository) Update(ctx context.Context, class *entity.Class) error {
	update := bson.M{
		"$set": class,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": class.ID}, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar aula: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("aula não encontrada")
	}

	return nil
}

func (r *ClassRepository) IncrementEnrollmentWithVersion(ctx context.Context, classID primitive.ObjectID, currentVersion int) error {
	update := bson.M{
		"$inc": bson.M{
			"current_enrolled": 1,
			"version": 1,
		},
		"$set": bson.M{
			"updated_at": primitive.NewDateTimeFromTime(entity.Class{}.UpdatedAt),
		},
	}

	filter := bson.M{
		"_id": classID,
		"version": currentVersion,
		"$expr": bson.M{
			"$lt": []interface{}{"$current_enrolled", "$max_capacity"},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("erro ao incrementar inscrição: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("aula sem vagas ou versão desatualizada")
	}

	return nil
}

func (r *ClassRepository) DecrementEnrollment(ctx context.Context, classID primitive.ObjectID) error {
	update := bson.M{
		"$inc": bson.M{
			"current_enrolled": -1,
		},
	}

	filter := bson.M{
		"_id": classID,
		"current_enrolled": bson.M{"$gt": 0},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ClassRepository) WithTransaction(ctx context.Context, fn func(context.Context, mongo.SessionContext) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("erro ao iniciar sessão: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(ctx, sc)
	})

	return err
}


package database

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Conectando ao MongoDB...",
		zap.String("database", dbName),
	)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão com MongoDB: %w", err)
	}

	logger.Info("Conectado ao MongoDB com sucesso",
		zap.String("database", dbName),
	)

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := m.Client.Disconnect(ctx); err != nil {
		logger.Error("Erro ao desconectar do MongoDB", zap.Error(err))
		return fmt.Errorf("erro ao desconectar do MongoDB: %w", err)
	}
	
	logger.Info("Desconectado do MongoDB com sucesso")
	return nil
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

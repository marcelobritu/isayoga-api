package user

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type GetUserUseCase struct {
	userRepo repository.UserRepository
}

func NewGetUserUseCase(userRepo repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Warn("ID inválido fornecido", zap.String("id", id))
		return nil, fmt.Errorf("ID inválido")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, objectID)
	if err != nil {
		logger.Error("Erro ao buscar usuário",
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

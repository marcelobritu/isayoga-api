package user

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type DeleteUserUseCase struct {
	userRepo repository.UserRepository
}

func NewDeleteUserUseCase(userRepo repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Warn("ID inválido fornecido", zap.String("id", id))
		return fmt.Errorf("ID inválido")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := uc.userRepo.Delete(ctx, objectID); err != nil {
		logger.Error("Erro ao deletar usuário",
			zap.Error(err),
			zap.String("id", id),
		)
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	logger.Info("Usuário deletado com sucesso", zap.String("id", id))

	return nil
}

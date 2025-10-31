package user

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	users, err := uc.userRepo.FindAll(ctx)
	if err != nil {
		logger.Error("Erro ao listar usuários", zap.Error(err))
		return nil, fmt.Errorf("erro ao listar usuários: %w", err)
	}

	logger.Debug("Usuários listados", zap.Int("count", len(users)))

	return users, nil
}

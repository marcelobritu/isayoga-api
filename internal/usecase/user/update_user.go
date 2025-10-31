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

type UpdateUserInput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		logger.Warn("ID inválido fornecido", zap.String("id", input.ID))
		return nil, fmt.Errorf("ID inválido")
	}

	if input.Name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	if input.Email == "" {
		return nil, fmt.Errorf("email é obrigatório")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, objectID)
	if err != nil {
		logger.Error("Erro ao buscar usuário para atualização",
			zap.Error(err),
			zap.String("id", input.ID),
		)
		return nil, fmt.Errorf("usuário não encontrado")
	}

	user.Update(input.Name, input.Email)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.Error("Erro ao atualizar usuário",
			zap.Error(err),
			zap.String("id", input.ID),
		)
		return nil, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	logger.Info("Usuário atualizado com sucesso",
		zap.String("id", user.ID.Hex()),
		zap.String("email", user.Email),
	)

	return user, nil
}

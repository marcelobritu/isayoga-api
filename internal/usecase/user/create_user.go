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

type CreateUserInput struct {
	Name  string          `json:"name"`
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
}

type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*entity.User, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	if input.Email == "" {
		return nil, fmt.Errorf("email é obrigatório")
	}

	if input.Role == "" {
		input.Role = entity.RoleStudent
	}

	if input.Role != entity.RoleStudent && input.Role != entity.RoleInstructor && input.Role != entity.RoleAdmin {
		return nil, fmt.Errorf("role inválido: deve ser student, instructor ou admin")
	}

	user := entity.NewUser(input.Name, input.Email)
	user.Role = input.Role

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.Error("Erro ao criar usuário no repositório",
			zap.Error(err),
			zap.String("email", input.Email),
		)
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	logger.Info("Usuário criado com sucesso",
		zap.String("id", user.ID.Hex()),
		zap.String("email", user.Email),
		zap.String("name", user.Name),
	)

	return user, nil
}

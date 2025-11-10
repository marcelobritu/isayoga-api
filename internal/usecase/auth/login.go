package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"github.com/marcelobritu/isayoga-api/pkg/auth"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

type LoginUseCase struct {
	userRepo repository.UserRepository
}

func NewLoginUseCase(userRepo repository.UserRepository) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	if input.Email == "" || input.Password == "" {
		return nil, fmt.Errorf("email e senha são obrigatórios")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		logger.Warn("Tentativa de login com email não encontrado", zap.String("email", input.Email))
		return nil, fmt.Errorf("credenciais inválidas")
	}

	if !user.CheckPassword(input.Password) {
		logger.Warn("Tentativa de login com senha incorreta", zap.String("email", input.Email))
		return nil, fmt.Errorf("credenciais inválidas")
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.Error("Erro ao gerar token JWT", zap.Error(err), zap.String("user_id", user.ID.Hex()))
		return nil, fmt.Errorf("erro ao gerar token de autenticação")
	}

	logger.Info("Usuário logado com sucesso", zap.String("user_id", user.ID.Hex()), zap.String("email", user.Email))

	return &LoginOutput{
		Token: token,
		User: struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Role  string `json:"role"`
		}{
			ID:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

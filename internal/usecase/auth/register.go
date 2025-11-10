package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	pkgAuth "github.com/marcelobritu/isayoga-api/pkg/auth"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type RegisterInput struct {
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     entity.UserRole `json:"role"`
}

type RegisterOutput struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

type RegisterUseCase struct {
	userRepo repository.UserRepository
}

func NewRegisterUseCase(userRepo repository.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, fmt.Errorf("nome, email e senha são obrigatórios")
	}

	if len(input.Password) < 6 {
		return nil, fmt.Errorf("a senha deve ter no mínimo 6 caracteres")
	}

	if input.Role == "" {
		input.Role = entity.RoleStudent
	}

	if input.Role != entity.RoleStudent && input.Role != entity.RoleInstructor && input.Role != entity.RoleAdmin {
		return nil, fmt.Errorf("role inválido: deve ser student, instructor ou admin")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existingUser, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email já cadastrado")
	}

	user, err := entity.NewUser(input.Name, input.Email, input.Password, input.Role)
	if err != nil {
		logger.Error("Erro ao criar hash da senha", zap.Error(err))
		return nil, fmt.Errorf("erro ao criar usuário")
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.Error("Erro ao criar usuário no repositório", zap.Error(err))
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	token, err := pkgAuth.GenerateToken(user)
	if err != nil {
		logger.Error("Erro ao gerar token", zap.Error(err))
		return nil, fmt.Errorf("erro ao gerar token de autenticação")
	}

	logger.Info("Usuário registrado com sucesso",
		zap.String("id", user.ID.Hex()),
		zap.String("email", user.Email),
		zap.String("role", string(user.Role)),
	)

	output := &RegisterOutput{
		Token: token,
	}
	output.User.ID = user.ID.Hex()
	output.User.Name = user.Name
	output.User.Email = user.Email
	output.User.Role = string(user.Role)

	return output, nil
}


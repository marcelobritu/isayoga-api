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

type ChangePasswordInput struct {
	UserID          string `json:"user_id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangePasswordUseCase struct {
	userRepo repository.UserRepository
}

func NewChangePasswordUseCase(userRepo repository.UserRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) error {
	if input.CurrentPassword == "" || input.NewPassword == "" {
		return fmt.Errorf("senha atual e nova senha são obrigatórias")
	}

	if len(input.NewPassword) < 6 {
		return fmt.Errorf("a nova senha deve ter no mínimo 6 caracteres")
	}

	if input.CurrentPassword == input.NewPassword {
		return fmt.Errorf("a nova senha deve ser diferente da senha atual")
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return fmt.Errorf("ID de usuário inválido")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("usuário não encontrado")
	}

	if !user.CheckPassword(input.CurrentPassword) {
		logger.Warn("Tentativa de alteração de senha com senha incorreta",
			zap.String("user_id", user.ID.Hex()),
			zap.String("email", user.Email),
		)
		return fmt.Errorf("senha atual incorreta")
	}

	if err := user.SetPassword(input.NewPassword); err != nil {
		logger.Error("Erro ao criar hash da nova senha", zap.Error(err))
		return fmt.Errorf("erro ao alterar senha")
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.Error("Erro ao atualizar senha no repositório",
			zap.Error(err),
			zap.String("user_id", user.ID.Hex()),
		)
		return fmt.Errorf("erro ao atualizar senha")
	}

	logger.Info("Senha alterada com sucesso",
		zap.String("user_id", user.ID.Hex()),
		zap.String("email", user.Email),
	)

	return nil
}


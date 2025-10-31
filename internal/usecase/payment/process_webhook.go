package payment

import (
	"context"
	"fmt"

	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type ProcessWebhookUseCase struct {
	paymentRepo    repository.PaymentRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewProcessWebhookUseCase(
	paymentRepo repository.PaymentRepository,
	enrollmentRepo repository.EnrollmentRepository,
) *ProcessWebhookUseCase {
	return &ProcessWebhookUseCase{
		paymentRepo:    paymentRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

type WebhookInput struct {
	Action   string `json:"action"`
	Type     string `json:"type"`
	DataID   string `json:"data.id"`
	LiveMode bool   `json:"live_mode"`
}

func (uc *ProcessWebhookUseCase) Execute(ctx context.Context, input WebhookInput) error {
	if input.Type != "payment" {
		logger.Info("Webhook ignorado: tipo não suportado", zap.String("type", input.Type))
		return nil
	}

	logger.Info("Processando webhook de pagamento",
		zap.String("action", input.Action),
		zap.String("payment_id", input.DataID),
	)

	payment, err := uc.paymentRepo.FindByMercadoPagoID(ctx, input.DataID)
	if err != nil {
		return fmt.Errorf("pagamento não encontrado: %w", err)
	}

	enrollment, err := uc.enrollmentRepo.FindByID(ctx, payment.EnrollmentID)
	if err != nil {
		return err
	}

	if input.Action == "payment.created" || input.Action == "payment.updated" {
		status := "approved"

		payment.UpdateFromMercadoPago(input.DataID, status, "mercadopago")
		if err := uc.paymentRepo.Update(ctx, payment); err != nil {
			return err
		}

		if status == "approved" && enrollment.Status == "pending" {
			enrollment.Confirm(input.DataID)
			if err := uc.enrollmentRepo.Update(ctx, enrollment); err != nil {
				return err
			}

			logger.Info("Inscrição confirmada com sucesso",
				zap.String("enrollment_id", enrollment.ID.Hex()),
				zap.String("payment_id", input.DataID),
			)
		}
	}

	return nil
}

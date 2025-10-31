package enrollment

import (
	"context"
	"fmt"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"github.com/marcelobritu/isayoga-api/internal/infrastructure/payment"
	"github.com/marcelobritu/isayoga-api/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnrollStudentUseCase struct {
	classRepo      repository.ClassRepository
	enrollmentRepo repository.EnrollmentRepository
	paymentRepo    repository.PaymentRepository
	userRepo       repository.UserRepository
	mercadoPago    *payment.MercadoPagoClient
	config         *config.Config
}

func NewEnrollStudentUseCase(
	classRepo repository.ClassRepository,
	enrollmentRepo repository.EnrollmentRepository,
	paymentRepo repository.PaymentRepository,
	userRepo repository.UserRepository,
	mercadoPago *payment.MercadoPagoClient,
	config *config.Config,
) *EnrollStudentUseCase {
	return &EnrollStudentUseCase{
		classRepo:      classRepo,
		enrollmentRepo: enrollmentRepo,
		paymentRepo:    paymentRepo,
		userRepo:       userRepo,
		mercadoPago:    mercadoPago,
		config:         config,
	}
}

type EnrollStudentInput struct {
	UserID  string `json:"user_id"`
	ClassID string `json:"class_id"`
}

type EnrollStudentOutput struct {
	Enrollment *entity.Enrollment `json:"enrollment"`
	Payment    *entity.Payment    `json:"payment"`
	PaymentURL string             `json:"payment_url"`
}

func (uc *EnrollStudentUseCase) Execute(ctx context.Context, input EnrollStudentInput) (*EnrollStudentOutput, error) {
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user_id inválido")
	}

	classID, err := primitive.ObjectIDFromHex(input.ClassID)
	if err != nil {
		return nil, fmt.Errorf("class_id inválido")
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if !user.IsStudent() {
		return nil, fmt.Errorf("apenas estudantes podem se inscrever em aulas")
	}

	existing, err := uc.enrollmentRepo.FindByUserAndClass(ctx, userID, classID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("usuário já está inscrito nesta aula")
	}

	var result *EnrollStudentOutput

	for {
		class, err := uc.classRepo.FindByID(ctx, classID)
		if err != nil {
			return nil, err
		}

		if !class.HasAvailableSpots() {
			return nil, fmt.Errorf("aula sem vagas disponíveis")
		}

		enrollment := entity.NewEnrollment(userID, classID)
		paymentEntity := entity.NewPayment(enrollment.ID, class.PriceInCents)

		err = uc.classRepo.WithTransaction(ctx, func(ctx context.Context, sc mongo.SessionContext) error {
			if err := uc.classRepo.IncrementEnrollmentWithVersion(sc, classID, class.Version); err != nil {
				return err
			}

			if err := uc.enrollmentRepo.Create(sc, enrollment); err != nil {
				return err
			}

			preferenceResp, err := uc.mercadoPago.CreatePreference(sc, &payment.PreferenceRequest{
				Title:       class.Title,
				Description: class.Description,
				Quantity:    1,
				UnitPrice:   float64(class.PriceInCents) / 100.0,
				ExternalRef: enrollment.ID.Hex(),
				NotifyURL:   uc.config.MercadoPago.NotifyURL,
				BackURL:     uc.config.MercadoPago.BackURL,
			})
			if err != nil {
				return fmt.Errorf("erro ao criar preferência de pagamento: %w", err)
			}

			paymentEntity.SetPreference(preferenceResp.ID, preferenceResp.InitPointURL)

			if err := uc.paymentRepo.Create(sc, paymentEntity); err != nil {
				return err
			}

			result = &EnrollStudentOutput{
				Enrollment: enrollment,
				Payment:    paymentEntity,
				PaymentURL: preferenceResp.InitPointURL,
			}

			return nil
		})

		if err == nil {
			return result, nil
		}
	}
}

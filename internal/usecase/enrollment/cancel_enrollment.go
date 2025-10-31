package enrollment

import (
	"context"
	"fmt"

	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CancelEnrollmentUseCase struct {
	enrollmentRepo repository.EnrollmentRepository
	classRepo      repository.ClassRepository
}

func NewCancelEnrollmentUseCase(
	enrollmentRepo repository.EnrollmentRepository,
	classRepo repository.ClassRepository,
) *CancelEnrollmentUseCase {
	return &CancelEnrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		classRepo:      classRepo,
	}
}

func (uc *CancelEnrollmentUseCase) Execute(ctx context.Context, enrollmentID string) error {
	id, err := primitive.ObjectIDFromHex(enrollmentID)
	if err != nil {
		return fmt.Errorf("enrollment_id inválido")
	}

	enrollment, err := uc.enrollmentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if enrollment.Status != "confirmed" {
		return fmt.Errorf("apenas inscrições confirmadas podem ser canceladas")
	}

	enrollment.Cancel()
	if err := uc.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return err
	}

	if err := uc.classRepo.DecrementEnrollment(ctx, enrollment.ClassID); err != nil {
		return err
	}

	return nil
}


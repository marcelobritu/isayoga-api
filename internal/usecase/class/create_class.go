package class

import (
	"context"
	"time"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateClassUseCase struct {
	classRepo repository.ClassRepository
}

func NewCreateClassUseCase(classRepo repository.ClassRepository) *CreateClassUseCase {
	return &CreateClassUseCase{
		classRepo: classRepo,
	}
}

type CreateClassInput struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	InstructorID   string    `json:"instructor_id"`
	InstructorName string    `json:"instructor_name"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	MaxCapacity    int       `json:"max_capacity"`
	PriceInCents   int64     `json:"price_in_cents"`
}

func (uc *CreateClassUseCase) Execute(ctx context.Context, input CreateClassInput) (*entity.Class, error) {
	instructorID, err := primitive.ObjectIDFromHex(input.InstructorID)
	if err != nil {
		return nil, err
	}

	class := entity.NewClass(
		input.Title,
		input.Description,
		instructorID,
		input.InstructorName,
		input.StartTime,
		input.EndTime,
		input.MaxCapacity,
		input.PriceInCents,
	)

	if err := uc.classRepo.Create(ctx, class); err != nil {
		return nil, err
	}

	return class, nil
}


package class

import (
	"context"

	"github.com/marcelobritu/isayoga-api/internal/domain/entity"
	"github.com/marcelobritu/isayoga-api/internal/domain/repository"
)

type ListClassesUseCase struct {
	classRepo repository.ClassRepository
}

func NewListClassesUseCase(classRepo repository.ClassRepository) *ListClassesUseCase {
	return &ListClassesUseCase{
		classRepo: classRepo,
	}
}

func (uc *ListClassesUseCase) Execute(ctx context.Context) ([]*entity.Class, error) {
	return uc.classRepo.FindAll(ctx)
}


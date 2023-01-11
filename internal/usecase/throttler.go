package usecase

import (
	"context"
	"fmt"

	"github.com/FRahimov84/throttler/internal/entity"
)

type ThrottlerUseCase struct {
	repo        ThrottlerRepo
	externalSvc ExternalSvc
}

func New(r ThrottlerRepo, e ExternalSvc) *ThrottlerUseCase {
	return &ThrottlerUseCase{repo: r, externalSvc: e}
}

func (uc *ThrottlerUseCase) NewRequest(ctx context.Context) (entity.UUID, error) {
	uuid, err := uc.repo.StoreRequest(ctx)
	if err != nil {
		return entity.UUID{}, fmt.Errorf("ThrottlerUseCase - NewRequest - s.repo.StoreRequest: %w", err)
	}

	return uuid, nil
}

func (uc *ThrottlerUseCase) RequestByUUID(ctx context.Context, uuid entity.UUID) (entity.Request, error) {
	req, err := uc.repo.RequestByUUID(ctx, uuid)
	if err != nil {
		return entity.Request{}, fmt.Errorf("ThrottlerUseCase - RequestByUUID - s.repo.RequestByUUID: %w", err)
	}

	return req, nil
}

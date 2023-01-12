package usecase

import (
	"context"
	"fmt"
	"github.com/FRahimov84/throttler/internal/entity"
	"github.com/FRahimov84/throttler/pkg/logger"
	"go.uber.org/zap"
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

func (uc *ThrottlerUseCase) Call(ctx context.Context) {
	l := logger.FromCtx(ctx)
	request, err := uc.repo.GetReqByFilter(ctx, entity.Filter{Status: "new"})
	if err != nil {
		if err == entity.NoRows {
			return
		}
		l.Error("ThrottlerUseCase - Call - s.repo.GetReqByFilter", zap.Error(err))
		return
	}

	resp, err := uc.externalSvc.Call(ctx, request)
	if err != nil {
		l.Error("ThrottlerUseCase - Call - s.externalSvc.Call", zap.Error(err))
		return
	}

	err = uc.repo.UpdateRequest(ctx, entity.Request{
		ID:       request.ID,
		Status:   resp.Status,
		Response: resp.Response,
	})
	if err != nil {
		l.Error("ThrottlerUseCase - Call - s.repo.UpdateRequest", zap.Error(err))
		return
	}

	l.Info("ThrottlerUseCase - Call - request processed", zap.Any("request", request), zap.Any("resp", resp))
}

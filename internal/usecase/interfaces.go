package usecase

import (
	"context"

	"github.com/FRahimov84/throttler/internal/entity"
)

type (
	Throttler interface {
		NewRequest(context.Context) (entity.UUID, error)
		RequestByUUID(context.Context, entity.UUID) (entity.Request, error)
	}

	ThrottlerRepo interface {
		StoreRequest(context.Context) (entity.UUID, error)
		RequestByUUID(context.Context, entity.UUID) (entity.Request, error)
		GetReqByFilter(context.Context, entity.Filter) (entity.Request, error)
		UpdateRequest(context.Context, entity.Request) error
	}

	ExternalSvc interface {
		Call(context.Context, entity.Request) (entity.ExternalSvcResp, error)
	}

	Caller interface {
		Call(ctx context.Context)
	}
)

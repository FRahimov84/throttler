package usecase

import (
	"context"

	"github.com/FRahimov84/throttler/internal/entity"
)

type (
	Throttler interface {
		NewRequest(context.Context, entity.Request) (entity.UUID, error)
		GetRequestByID(context.Context, entity.UUID) (entity.Request, error)
	}

	ThrottlerRepo interface {
		StoreRequest(context.Context, entity.Request) (entity.UUID, error)
		GetRequestByID(context.Context, entity.UUID) (entity.Request, error)
		GetRequestByFilter(context.Context, entity.Filter) (entity.Request, error)
		UpdateRequest(context.Context, entity.Request) error
	}

	ExternalSvc interface {
		CallRemoteSvc(context.Context, entity.Request) (entity.ExternalSvcResp, error)
	}
)

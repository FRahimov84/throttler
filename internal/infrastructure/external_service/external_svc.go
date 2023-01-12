package external_service

import (
	"context"
	"fmt"
	"time"

	"github.com/FRahimov84/throttler/internal/entity"
)

type ExternalService struct {
	URL string
}

func New(url string) *ExternalService {
	return &ExternalService{URL: url}
}

func (e *ExternalService) Call(ctx context.Context, r entity.Request) (entity.ExternalSvcResp, error) {
	//	TODO: implement caller
	time.Sleep(100 * time.Millisecond)
	return entity.ExternalSvcResp{
		Status:   "Success",
		Response: fmt.Sprintf("Request %s processed", r.ID.String()),
	}, nil
}

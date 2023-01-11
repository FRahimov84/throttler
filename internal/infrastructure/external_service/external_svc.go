package external_service

import (
	"context"

	"github.com/FRahimov84/throttler/internal/entity"
)

type ExternalService struct {
	URL string
	N   int
	K   int
	X   int
}

func New(url string, n int, k int, x int) *ExternalService {
	return &ExternalService{URL: url, N: n, K: k, X: x}
}


func (e *ExternalService) Call(ctx context.Context) (entity.ExternalSvcResp, error) {
	panic("implement me")
}

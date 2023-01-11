package postgres

import (
	"context"
	"fmt"

	"github.com/FRahimov84/throttler/internal/entity"
	"github.com/FRahimov84/throttler/pkg/postgres"
)

type ThrottlerRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *ThrottlerRepo {
	return &ThrottlerRepo{pg}
}

func (r *ThrottlerRepo) StoreRequest(ctx context.Context) (uuid entity.UUID, err error) {
	sql, args, err := r.Builder.Insert("requests").
		Columns("status").
		Values("new").Suffix("RETURNING id").ToSql()
	if err != nil {
		return entity.UUID{}, fmt.Errorf("ThrottlerRepo - StoreRequest - r.Builder: %w", err)
	}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&uuid)
	if err != nil {
		return entity.UUID{}, fmt.Errorf("ThrottlerRepo - StoreRequest - r.Pool.QueryRow: %w", err)
	}

	return uuid, nil
}

func (r *ThrottlerRepo) RequestByUUID(ctx context.Context, uuid entity.UUID) (req entity.Request,err error) {
	sql, args, err := r.Builder.Select("id", "status", "response").
		From("requests").
		Where("id = ?", uuid).
		ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("ThrottlerRepo - RequestByUUID - r.Builder: %w", err)
	}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&req.ID, &req.Status, &req.Response)
	if err != nil {
		return entity.Request{}, fmt.Errorf("ThrottlerRepo - RequestByUUID - r.Pool.QueryRow: %w", err)
	}

	return req, nil
}

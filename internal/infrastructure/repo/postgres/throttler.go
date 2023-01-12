package postgres

import (
	"context"
	"fmt"
	"github.com/FRahimov84/throttler/internal/entity"
	"github.com/FRahimov84/throttler/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

type ThrottlerRepo struct {
	*postgres.Postgres
}

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

func (r *ThrottlerRepo) RequestByUUID(ctx context.Context, uuid entity.UUID) (req entity.Request, err error) {
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

func (r *ThrottlerRepo) GetReqByFilter(ctx context.Context, filter entity.Filter) (req entity.Request, err error) {
	builder := r.Builder.Select("id", "status", "response").From("requests")
	if filter.ID != entity.EmptyUUID {
		builder = builder.Where("id = ?", filter.ID)
	}
	if filter.Status != "" {
		builder = builder.Where("status = ?", filter.Status)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("ThrottlerRepo - GetReqByFilter - r.Builder: %w", err)
	}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&req.ID, &req.Status, &req.Response)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Request{}, entity.NoRows
		}
		return entity.Request{}, fmt.Errorf("ThrottlerRepo - RequestByUUID - r.Pool.QueryRow: %w", err)
	}

	return
}

func (r *ThrottlerRepo) UpdateRequest(ctx context.Context, request entity.Request) error {
	sql, args, err := r.Builder.Update("requests").
		Set("status", request.Status).
		Set("response", request.Response).
		Where("id = ?", request.ID).ToSql()
	if err != nil {
		return fmt.Errorf("ThrottlerRepo - UpdateRequest - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ThrottlerRepo - UpdateRequest - r.Pool.Exec: %w", err)
	}

	return nil
}

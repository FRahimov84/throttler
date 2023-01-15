package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/FRahimov84/throttler/internal/entity"
	"github.com/FRahimov84/throttler/pkg/postgres"
)

type ThrottlerRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *ThrottlerRepo {
	repo := &ThrottlerRepo{pg}
	repo.Init()

	return repo
}

func (r *ThrottlerRepo) Init() {
	r.Pool.Exec(context.Background(), `UPDATE requests SET status = 'new' WHERE status = 'processing'`)
}

func (r *ThrottlerRepo) StoreRequest(ctx context.Context, request entity.Request) (uuid entity.UUID, err error) {
	sql, args, err := r.Builder.Insert("requests").
		Columns("status", "response").
		Values(request.Status, request.Response).Suffix("RETURNING id").ToSql()
	if err != nil {
		return entity.UUID{}, fmt.Errorf("PostgresRepo - StoreRequest - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&uuid)
	if err != nil {
		return entity.UUID{}, fmt.Errorf("PostgresRepo - StoreRequest - r.Pool.QueryRow: %w", err)
	}

	return uuid, nil
}

func (r *ThrottlerRepo) GetRequestByID(ctx context.Context, uuid entity.UUID) (req entity.Request, err error) {
	sql, args, err := r.Builder.Select("id", "status", "response").
		From("requests").
		Where("id = ?", uuid).
		ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("PostgresRepo - RequestByUUID - r.Builder: %w", err)
	}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&req.ID, &req.Status, &req.Response)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Request{}, entity.RepoNotFoundErr
		}
		return entity.Request{}, fmt.Errorf("PostgresRepo - RequestByUUID - r.Pool.QueryRow: %w", err)
	}

	return req, nil
}

func (r *ThrottlerRepo) GetRequestByFilter(ctx context.Context, filter entity.Filter) (req entity.Request, err error) {
	//builder := r.Builder.Select("id", "status", "response").From("requests")
	//if filter.ID != entity.EmptyUUID {
	//	builder = builder.Where("id = ?", filter.ID)
	//}
	//if filter.Status != "" {
	//	builder = builder.Where("status = ?", filter.Status)
	//}
	//sql, args, err := builder.ToSql()
	//if err != nil {
	//	return entity.Request{}, fmt.Errorf("PostgresRepo - GetReqByFilter - r.Builder: %w", err)
	//}
	//err = r.Pool.QueryRow(ctx, sql, args...).Scan(&req.ID, &req.Status, &req.Response)
	//if err != nil {
	//	if err == pgx.ErrNoRows {
	//		return entity.Request{}, entity.RepoNotFoundErr
	//	}
	//	return entity.Request{}, fmt.Errorf("PostgresRepo - RequestByUUID - r.Pool.QueryRow: %w", err)
	//}

	return r.GetRequestForProcessing(ctx)
}

func (r *ThrottlerRepo) GetRequestForProcessing(ctx context.Context) (req entity.Request, err error) {
	// case with select for update
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entity.Request{}, fmt.Errorf("PostgresRepo - GetRequestForUpdate - r.Pool.Begin: %w", err)
	}

	err = tx.QueryRow(ctx, `SELECT id, status, response FROM requests WHERE status = 'new' LIMIT 1 FOR UPDATE`).
		Scan(&req.ID, &req.Status, &req.Response)
	if err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows{
			return entity.Request{}, entity.RepoNotFoundErr
		}
		return entity.Request{}, fmt.Errorf("PostgresRepo - GetRequestForUpdate - tx.QueryRow: %w", err)
	}

	exec, err := tx.Exec(ctx, `UPDATE requests SET status = 'processing' WHERE id = $1`, req.ID)
	if err != nil {
		tx.Rollback(ctx)
		return entity.Request{}, fmt.Errorf("PostgresRepo - GetRequestForUpdate - tx.Exec: %w", err)
	}
	if exec.RowsAffected() != 1 {
		tx.Rollback(ctx)
		return entity.Request{},
			fmt.Errorf("PostgresRepo - GetRequestForUpdate - tx.Exec finished with <no rows affected> %d",
				exec.RowsAffected())
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return entity.Request{}, fmt.Errorf("PostgresRepo - GetRequestForUpdate - tx.Commit: %w", err)
	}

	return req, nil
}

func (r *ThrottlerRepo) UpdateRequest(ctx context.Context, request entity.Request) error {
	sql, args, err := r.Builder.Update("requests").
		Set("status", request.Status).
		Set("response", request.Response).
		Where("id = ?", request.ID).ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - UpdateRequest - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PostgresRepo - UpdateRequest - r.Pool.Exec: %w", err)
	}

	return nil
}

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FRahimov84/throttler/internal/entity"
	"github.com/go-redis/redis/v8"
)

type ThrottlerRedis struct {
	*redis.Client
}

func New(rd *redis.Client) *ThrottlerRedis {
	return &ThrottlerRedis{rd}
}

func (r *ThrottlerRedis) StoreRequest(ctx context.Context, request entity.Request) (entity.UUID, error) {
	if request.ID == entity.EmptyUUID {
		request.ID = entity.New()
	}
	marshal, err := json.Marshal(request)
	if err != nil {
		return entity.UUID{}, fmt.Errorf("RedisRepo - StoreRequest - json.Marshal: %w", err)
	}

	err = r.Set(ctx, request.ID.String()+":"+request.Status, string(marshal), 0).Err()
	if err != nil {
		return entity.UUID{}, fmt.Errorf("RedisRepo - StoreRequest - repo.Set: %w", err)
	}

	return request.ID, nil
}

func (r *ThrottlerRedis) GetRequestByID(ctx context.Context, uuid entity.UUID) (req entity.Request, err error) {
	keyPattern := uuid.String() + ":*"

	return r.getFirstInByKeyPattern(ctx, keyPattern)
}

func (r *ThrottlerRedis) GetRequestByFilter(ctx context.Context, filter entity.Filter) (req entity.Request, err error) {
	keyPattern := ""
	if filter.ID != entity.EmptyUUID {
		if filter.Status != "" {
			keyPattern = filter.ID.String() + ":" + filter.Status
		} else {
			keyPattern = filter.ID.String() + ":*"
		}
	} else {
		if filter.Status != "" {
			keyPattern = "*:" + filter.Status
		} else {
			keyPattern = "*:*"
		}
	}

	return r.getFirstInByKeyPattern(ctx, keyPattern)
}

func (r *ThrottlerRedis) getFirstInByKeyPattern(ctx context.Context, keyPattern string) (req entity.Request, err error) {
	keys, err := r.Keys(ctx, keyPattern).Result()
	if err != nil {
		return entity.Request{}, fmt.Errorf("RedisRepo - GetRequestByFilter - repo.Keys: %w", err)
	}

	if len(keys) > 0 {
		val, err := r.Get(ctx, keys[0]).Result()
		if err == redis.Nil {
			return entity.Request{}, entity.RepoNotFoundErr
		} else if err != nil {
			return entity.Request{}, fmt.Errorf("RedisRepo - getFirstInByKeyPattern - repo.Get: %w", err)
		}

		err = json.Unmarshal([]byte(val), &req)
		if err != nil {
			return entity.Request{}, fmt.Errorf("RedisRepo - getFirstInByKeyPattern - json.Unmarshal: %w", err)
		}

		return req, nil
	}

	return entity.Request{}, entity.RepoNotFoundErr
}

func (r *ThrottlerRedis) UpdateRequest(ctx context.Context, request entity.Request) error {
	uuid, err := r.StoreRequest(ctx, request)
	if err != nil {
		return fmt.Errorf("RedisRepo - UpdateRequest - repo.StoreRequest: %w", err)
	}

	iter := r.Scan(ctx, 0, uuid.String()+":*", 0).Iterator()
	for iter.Next(ctx) {
		if iter.Val() != uuid.String()+":"+request.Status {
			r.Del(ctx, iter.Val())
		}

	}
	if err = iter.Err(); err != nil {
		return fmt.Errorf("RedisRepo - UpdateRequest - got error on deleting: %w", err)
	}

	return nil
}

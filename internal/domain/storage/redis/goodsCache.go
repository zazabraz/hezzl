package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"hezzl/internal/domain/storage"
	"hezzl/internal/domain/storage/dto"
	"strconv"
	"time"
)

type goodsCacheStorage struct {
	rdb               *redis.Client
	expirationSeconds time.Duration
}

func NewGoodsCache(rdb *redis.Client, expirationSeconds int) storage.GoodsCacheStorage {
	return &goodsCacheStorage{
		rdb:               rdb,
		expirationSeconds: time.Duration(expirationSeconds) * time.Second,
	}
}

const GoodsCacheKey = "goods"

func (g goodsCacheStorage) keyById(id int) string {
	return fmt.Sprintf("%s:%s", GoodsCacheKey, strconv.Itoa(id))
}

func (g goodsCacheStorage) Create(ctx context.Context, good *dto.Good) error {
	return g.rdb.Set(ctx, g.keyById(good.ID), good, g.expirationSeconds).Err()
}

func (g goodsCacheStorage) GetByIdAndProjectID(ctx context.Context, id int) (*dto.Good, error) {
	resStr, err := g.rdb.Get(ctx, g.keyById(id)).Result()
	if err != nil {
		return nil, err
	}

	var good *dto.Good
	err = json.Unmarshal([]byte(resStr), good)
	if err != nil {
		return nil, err
	}

	return good, err
}

func (g goodsCacheStorage) Invalidate(ctx context.Context, good *dto.Good) error {
	err := g.rdb.Del(ctx, g.keyById(good.ID)).Err()
	if err != nil {
		return err
	}
	err = g.rdb.Set(ctx, g.keyById(good.ID), good, g.expirationSeconds).Err()
	if err != nil {
		return err
	}
	return nil
}

func (g goodsCacheStorage) PopByIdAndProjectID(ctx context.Context, id int) error {
	return g.rdb.Del(ctx, g.keyById(id)).Err()
}

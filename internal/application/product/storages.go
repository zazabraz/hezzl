package product

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/redis/go-redis/v9"
	"hezzl/_pkg/postgresql"
	"hezzl/internal/domain/storage"
	"hezzl/internal/domain/storage/clickhouse"
	"hezzl/internal/domain/storage/postgres"
	redisStorage "hezzl/internal/domain/storage/redis"
)

type Storages struct {
	Goods      storage.GoodsStorage
	GoodsCache storage.GoodsCacheStorage
	GoodsEvent storage.GoodsEventStorage
}

func NewStorages(
	pgPool postgresql.PGXPool,

	rdb *redis.Client,
	redisExpirationSeconds int,

	conn *driver.Conn,
) *Storages {
	return &Storages{
		GoodsEvent: clickhouse.NewGoodsEventStorage(*conn),
		Goods:      postgres.NewGoodsStorage(pgPool),
		GoodsCache: redisStorage.NewGoodsCache(rdb, redisExpirationSeconds),
	}
}

package storage

import (
	"context"
	"hezzl/internal/domain/event"
)

type GoodsEventStorage interface {
	SaveEvent(ctx context.Context, change event.GoodsChange) error
}

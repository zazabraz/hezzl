package storage

import (
	"context"
	"hezzl/internal/domain/storage/dto"
)

type GoodsCacheStorage interface {
	Create(context.Context, *dto.Good) error
	GetByIdAndProjectID(context.Context, int) (*dto.Good, error)
	Invalidate(ctx context.Context, good *dto.Good) error
	PopByIdAndProjectID(ctx context.Context, id int) error
}

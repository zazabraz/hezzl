package storage

import (
	"context"
	"hezzl/internal/domain/storage/dto"
)

type GoodsStorage interface {
	Create(context.Context, *dto.Good) (*dto.Good, error)
	Get(ctx context.Context, id int) (*dto.Good, error)
	GetList(ctx context.Context, limit int, offset int) ([]*dto.Good, error)
	Update(context.Context, *dto.Good) (*dto.Good, error)
	DeleteByIDAndProjectID(ctx context.Context, id int, projectId int) (*dto.Good, error)
}

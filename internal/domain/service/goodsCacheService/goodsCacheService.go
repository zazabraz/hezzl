package goodsCacheService

import (
	"context"
	"github.com/sirupsen/logrus"
	"hezzl/internal/domain/storage"
	"hezzl/internal/domain/storage/dto"
)

type GoodsCacheService interface {
	Set(context.Context, *dto.Good) error
	GetById(context.Context, int) (*dto.Good, error)
	Invalidate(context.Context, *dto.Good) error
	PopByIdAndProjectID(ctx context.Context, id int) error
}

type goodsCacheService struct {
	log               logrus.Entry
	goodsCacheStorage storage.GoodsCacheStorage
}

func NewGoodsCacheService(
	log logrus.Entry,
	goodsCacheStorage storage.GoodsCacheStorage,
) GoodsCacheService {
	return &goodsCacheService{
		log:               log,
		goodsCacheStorage: goodsCacheStorage,
	}
}

func (g goodsCacheService) Set(ctx context.Context, good *dto.Good) error {
	return g.goodsCacheStorage.Create(ctx, good)
}

func (g goodsCacheService) GetById(ctx context.Context, id int) (*dto.Good, error) {
	return g.goodsCacheStorage.GetByIdAndProjectID(ctx, id)
}

func (g goodsCacheService) Invalidate(ctx context.Context, good *dto.Good) error {
	return g.goodsCacheStorage.Invalidate(ctx, good)
}

func (g goodsCacheService) PopByIdAndProjectID(ctx context.Context, id int) error {
	return g.goodsCacheStorage.PopByIdAndProjectID(ctx, id)
}

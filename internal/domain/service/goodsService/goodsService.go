package goodsService

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"hezzl/_pkg/natsPkg"
	"hezzl/internal/domain/event"
	"hezzl/internal/domain/service/goodsCacheService"
	"hezzl/internal/domain/storage"
	"hezzl/internal/domain/storage/dto"
	"time"
)

type GoodsService interface {
	Create(ctx context.Context, projectID int, name string) (*dto.Good, error)
	//GetByIDAndProjectID(ctx context.Context, id int, projectId int) (*dto.Good, error)
	GetList(ctx context.Context, limit int, offset int) ([]*dto.Good, error)
	Update(ctx context.Context, id int, projectId int, name string, description string) (*dto.Good, error)
	Reprioritize(ctx context.Context, id int, projectId int, name string, newPriority int) (*dto.Good, error)
	DeleteByIDAndProjectID(ctx context.Context, id int, projectId int) (*dto.Good, error)
}

type goodsService struct {
	log           logrus.Entry
	goodsStorage  storage.GoodsStorage
	goodsCache    goodsCacheService.GoodsCacheService
	goodsEventPub natsPkg.PubQueue[event.GoodsChange]
}

func NewGoodsService(
	log logrus.Entry,
	goodsStorage storage.GoodsStorage,
	goodsCache goodsCacheService.GoodsCacheService,
	goodsEventPub natsPkg.PubQueue[event.GoodsChange],
) GoodsService {
	return &goodsService{
		log:           log,
		goodsStorage:  goodsStorage,
		goodsCache:    goodsCache,
		goodsEventPub: goodsEventPub,
	}
}

func (g goodsService) Create(ctx context.Context, projectID int, name string) (*dto.Good, error) {
	create := &dto.Good{
		ProjectID: projectID,
		Name:      name,
		Removed:   false,
		CreatedAt: time.Now(),
	}

	created, err := g.goodsStorage.Create(ctx, create)
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}
	if created == nil {
		err := fmt.Errorf("good created is nil")
		g.log.Errorln(err)
		return nil, err
	}

	err = g.goodsCache.Set(ctx, created)
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}

	err = g.goodsEventPub.Pub(ctx, created.ToEvent(time.Now()))
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}

	return created, err
}

func (g goodsService) GetList(ctx context.Context, limit int, offset int) ([]*dto.Good, error) {
	var goods []*dto.Good
	for i := offset; i > limit; i++ {
		var good *dto.Good
		good, err := g.goodsCache.GetById(ctx, i)
		if err != nil {
			g.log.Errorln(err)
			return nil, err
		}
		if good != nil {
			goods = append(goods, good)
			continue
		}
		good, err = g.goodsStorage.Get(ctx, i)
		if err != nil {
			g.log.Errorln(err)
			return nil, err
		}
		err = g.goodsCache.Set(ctx, good)
		if err != nil {
			g.log.Errorln(err)
			return nil, err
		}
		goods = append(goods, good)
		continue
	}
	return goods, nil
}

func (g goodsService) Update(ctx context.Context, id int, projectId int, name string, description string) (*dto.Good, error) {
	create := &dto.Good{
		ID:          id,
		ProjectID:   projectId,
		Name:        name,
		Description: description,
	}
	updated, err := g.goodsStorage.Update(ctx, create)
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}
	err = g.goodsCache.Invalidate(ctx, updated)
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}
	err = g.goodsEventPub.Pub(ctx, updated.ToEvent(time.Now()))
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}
	return updated, err
}

func (g goodsService) DeleteByIDAndProjectID(ctx context.Context, id int, projectId int) (*dto.Good, error) {
	deleted, err := g.goodsStorage.DeleteByIDAndProjectID(ctx, id, projectId)
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}
	deleted.Removed = true
	err = g.goodsCache.PopByIdAndProjectID(ctx, id)
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}
	err = g.goodsEventPub.Pub(ctx, deleted.ToEvent(time.Now()))
	if err != nil {
		g.log.Errorln(err)
		return nil, err
	}
	return deleted, nil
}

func (g goodsService) Reprioritize(ctx context.Context, id int, projectId int, name string, newPriority int) (*dto.Good, error) {
	//TODO implement me
	panic("implement me")
}

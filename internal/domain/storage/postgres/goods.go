package postgres

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"hezzl/_pkg/postgresql"
	"hezzl/internal/domain/storage"
	"hezzl/internal/domain/storage/dto"
)

type goodsStorage struct {
	pool postgresql.PGXPool
}

func NewGoodsStorage(pool postgresql.PGXPool) storage.GoodsStorage {
	return &goodsStorage{pool: pool}
}

func (g goodsStorage) Create(ctx context.Context, good *dto.Good) (*dto.Good, error) {
	var created *dto.Good

	tx, err := g.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO goods
			VALUES (default,$1,$2,default,default,$3,$4) 
			RETURNING id,project_id,name,description,priority,removed,created_at
		`,
		good.ProjectID,
		good.Name,
		good.Removed,
		good.CreatedAt,
	).Scan(
		&created.ID,
		&created.ProjectID,
		&created.Name,
		&created.Description,
		&created.Priority,
		&created.Removed,
		&created.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.PgErrNoEffect
	}
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (g goodsStorage) Get(ctx context.Context, id int, projectId int) (*dto.Good, error) {
	var good *dto.Good
	err := g.pool.QueryRow(
		ctx,
		`
			SELECT id,project_id,name,description,priority,removed,created_at
			FROM goods
			WHERE id = $1 AND project_id = $2
		`,
		id,
		projectId,
	).Scan(
		&good.ID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return good, nil
}

func (g goodsStorage) GetList(ctx context.Context, limit int, offset int) ([]*dto.Good, error) {
	rows, err := g.pool.Query(
		ctx,
		`
			SELECT id,project_id,name,description,priority,removed,created_at
			FROM goods
			WHERE id >= $1
			LIMIT $2
		`,
		offset,
		limit,
	)

	goods := make([]*dto.Good, 0)
	err = pgxscan.ScanAll(&goods, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		return []*dto.Good{}, nil
	}
	if err != nil {
		return nil, err
	}

	return goods, nil
}

func (g goodsStorage) Update(ctx context.Context, good *dto.Good) (*dto.Good, error) {
	var updated *dto.Good

	tx, err := g.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = g.pool.QueryRow(
		ctx,
		`
			UPDATE goods
			SET id = $1, project_id = $2, name = $3, description = $4, priority = $5, removed = $6, created_at = $7
			WHERE id = $3 AND project_id = $4
			RETURNING id,project_id,name,description,priority,removed,created_at
		`,
		good.Name, good.Description,
		good.ID, good.ProjectID,
	).Scan(
		&updated.ID,
		&updated.ProjectID,
		&updated.Name,
		&updated.Description,
		&updated.Priority,
		&updated.Removed,
		&updated.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.PgErrNoEffect
	}
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (g goodsStorage) DeleteByIDAndProjectID(ctx context.Context, id int, projectId int) (*dto.Good, error) {
	var good *dto.Good

	tx, err := g.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = g.pool.QueryRow(
		ctx,
		`
			delete from goods
			WHERE id = $1 AND project_id = $2
			RETURNING id,project_id,name,description,priority,removed,created_at
		`,
		id,
		projectId,
	).Scan(
		&good.ID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.PgErrNoEffect
	}
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return good, nil
}

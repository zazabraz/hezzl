package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"hezzl/internal/domain/event"
	"hezzl/internal/domain/storage"
)

type goodsEventStorage struct {
	conn driver.Conn
}

func NewGoodsEventStorage(conn driver.Conn) storage.GoodsEventStorage {
	return &goodsEventStorage{conn: conn}
}

func (g *goodsEventStorage) SaveEvent(ctx context.Context, change event.GoodsChange) error {
	query :=
		`
        INSERT INTO goods 
        (Id, ProjectId, Name, Description, Priority, Removed, EventTime) 
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	err := g.conn.Exec(ctx, query, change.ID, change.ProjectID, change.Name,
		change.Description, change.Priority,
		boolToInt(change.Removed), change.EventTime)
	if err != nil {
		return fmt.Errorf("failed to execute save event: %w", err)
	}

	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

package dto

import (
	"hezzl/internal/domain/event"
	"time"
)

type Good struct {
	ID          int       `db:"id"`
	ProjectID   int       `db:"project_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

func (g *Good) ToEvent(eventTime time.Time) event.GoodsChange {
	return event.GoodsChange{
		ID:          g.ID,
		ProjectID:   g.ID,
		Name:        g.Name,
		Description: g.Description,
		Priority:    g.Priority,
		Removed:     g.Removed,
		EventTime:   eventTime,
	}
}

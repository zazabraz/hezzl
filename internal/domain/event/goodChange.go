package event

import (
	"encoding/json"
	"fmt"
	"hezzl/_pkg/natsPkg"
	"time"
)

type GoodsChange struct {
	ID          int       `json:"id" binding:"required"`
	ProjectID   int       `json:"projectId" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Priority    int       `json:"priority" binding:"required"`
	Removed     bool      `json:"removed" binding:"required"`
	EventTime   time.Time `json:"eventTime" binding:"required"`
}

func (gc GoodsChange) Bytes() ([]byte, error) {
	bs, err := json.Marshal(gc)
	if err != nil {
		return nil, fmt.Errorf("marshal GoodsChange to bytes: %w", err)
	}
	return bs, nil
}

// Unmarshal scans []byte and unmarshals and returns it
func (gc GoodsChange) Unmarshal(bytes []byte) (natsPkg.Message, error) {
	err := json.Unmarshal(bytes, &gc)
	if err != nil {
		return GoodsChange{}, fmt.Errorf("unmarshal bytes: %w", err)
	}
	return gc, nil
}

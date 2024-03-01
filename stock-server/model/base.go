package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GormSlice []string

func (g *GormSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, &g)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (g GormSlice) Value() (driver.Value, error) {
	if len(g) == 0 {
		return nil, nil
	}
	return json.Marshal(g)
}

type BaseModel struct {
	ID        int32          `grom:"primaryKey;type:int"`
	CratedAt  time.Time      `gorm:"column:create_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeleteAt  gorm.DeletedAt `gorm:"column:delete_time"`
	IsDelete  bool
}

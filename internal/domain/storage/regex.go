package storage

import (
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

type RegexRecord struct {
	ID        uint        `gorm:"column:id;primaryKey"`
	Name      string      `gorm:"column:name;default:none"`
	Username  string      `gorm:"column:username;not null"`
	Pattern   model.Regex `gorm:"column:pattern;type:jsonb;not null"`
	CreatedAt time.Time   `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

func (RegexRecord) TableName() string {
	return "regexes"
}

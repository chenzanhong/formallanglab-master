package storage

import (
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

type RegexRecord struct {
	ID        uint        `gorm:"primaryKey"`
	Name      string      `gorm:"default:none"`
	Username  string      `gorm:"not null"`
	Pattern   model.Regex `gorm:"type:jsonb;not null"`
	CreatedAt time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (RegexRecord) TableName() string {
	return "regexes"
}

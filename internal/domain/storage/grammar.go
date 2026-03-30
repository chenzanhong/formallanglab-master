package storage

import (
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

type GrammarRecord struct {
	ID          uint          `gorm:"primaryKey"`
	Name        string        `gorm:"default:none"`
	Username    string        `gorm:"not null"`
	Grammar     model.Grammar `gorm:"type:jsonb;not null;serializer:json"`
	GrammarHash string        `gorm:"type:char(64);not null"`
	CreatedAt   time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (GrammarRecord) TableName() string {
	return "grammars"
}

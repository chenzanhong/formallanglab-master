package storage

import (
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

type GrammarRecord struct {
	ID          uint          `gorm:"column:id;primaryKey"`
	Name        string        `gorm:"column:name;default:none"`
	Username    string        `gorm:"column:username;not null"`
	Grammar     model.Grammar `gorm:"column:grammar;type:jsonb;not null;serializer:json"`
	GrammarHash string        `gorm:"column:grammar_hash;type:char(64);not null"`
	CreatedAt   time.Time     `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

func (GrammarRecord) TableName() string {
	return "grammars"
}

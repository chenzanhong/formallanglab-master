package storage

import (
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

type AutomatonRecord struct {
	ID            uint            `gorm:"primaryKey"`
	Name          string          `gorm:"default:none"`
	Username      string          `gorm:"not null"`
	Automaton     model.Automaton `gorm:"type:jsonb;not null;serializer:json"`
	AutomatonHash string          `gorm:"type:char(64);not null"`
	CreatedAt     time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (AutomatonRecord) TableName() string {
	return "automatons"
}

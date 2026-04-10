package storage

import (
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

type AutomatonRecord struct {
	ID            uint            `gorm:"column:id;primaryKey"`
	Name          string          `gorm:"column:name;default:none"`
	Username      string          `gorm:"column:username;not null"`
	Automaton     model.Automaton `gorm:"column:automaton;type:jsonb;not null;serializer:json"`
	AutomatonHash string          `gorm:"column:automaton_hash;type:char(64);not null"`
	AutomatonType string          `gorm:"column:automaton_type;type:varchar(50);not null;default:'DFA'"`
	CreatedAt     time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

func (AutomatonRecord) TableName() string {
	return "automatons"
}

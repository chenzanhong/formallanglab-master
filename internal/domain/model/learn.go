package model

import (
	"time"
)

// LearnMaterial 学习资源模型
type LearnMaterial struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category" gorm:"type:varchar(50);not null"` // grammar/automaton/regex/general
	FileKey     string    `json:"file_key" gorm:"type:varchar(512);not null;unique"`
	FileName    string    `json:"file_name" gorm:"type:varchar(255);not null"`
	MimeType    string    `json:"mime_type" gorm:"type:varchar(100)"`
	SizeBytes   int64     `json:"size_bytes" gorm:"type:bigint"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (LearnMaterial) TableName() string {
	return "learn_materials"
}

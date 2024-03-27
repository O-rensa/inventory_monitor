package c_sharedTypes

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelID struct {
	ID uuid.UUID `gorm:"primaryKey"`
}

func NewModelID() *ModelID {
	modelID := &ModelID{
		ID: uuid.New(),
	}
	return modelID
}

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

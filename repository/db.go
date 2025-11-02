package repository

import (
	"time"

	"gorm.io/gorm"
)

type db struct {
	db *gorm.DB
}

type GormModel struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GormModelWithoutID struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

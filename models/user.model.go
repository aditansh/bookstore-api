package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	Name       string    `gorm:"not null" json:"name"`
	Username   string    `gorrm:"type:varchar(100);unique;not null" json:"username"`
	Email      string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string    `gorm:"not null" json:"-"`
	IsVerified bool      `gorm:"default:false" json:"-"`
	IsDeleted  bool      `gorm:"default:false" json:"-"`
	Role       string    `gorm:"default:'user'" json:"-"`
	Reviews    []Review  `gorm:"foreignKey:UserID" json:"reviews"`
	Cart       Cart      `gorm:"foreignKey:UserID" json:"cart"`
	Orders     []Order   `gorm:"foreignKey:UserID" json:"orders"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"-"`
}
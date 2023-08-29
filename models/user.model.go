package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	Name       string    `gorm:"not null" json:"name"`
	Username   string    `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email      string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string    `gorm:"not null" json:"-"`
	IsVerified bool      `gorm:"default:false" json:"isVerified"`
	IsActive   bool      `gorm:"default:true" json:"isActive"`
	IsDeleted  bool      `gorm:"default:false" json:"isDeleted"`
	Role       string    `gorm:"default:'user'" json:"role"`
	Reviews    []Review  `gorm:"foreignKey:UserID" json:"reviews"`
	Cart       Cart      `gorm:"foreignKey:UserID" json:"cart"`
	Orders     []Order   `gorm:"foreignKey:UserID" json:"orders"`
	IsFlagged  bool      `gorm:"default:false" json:"isFlagged"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"-"`
}

package models

import (
	"time"

	"github.com/google/uuid"
	pq "github.com/lib/pq"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	Name       string         `gorm:"not null" json:"name"`
	Username   string         `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email      string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string         `gorm:"not null" json:"-"`
	IsVerified bool           `gorm:"default:false" json:"isVerified"`
	IsActive   bool           `gorm:"default:true" json:"isActive"`
	IsDeleted  bool           `gorm:"default:false" json:"isDeleted"`
	Role       string         `gorm:"default:'user'" json:"role"`
	Reviews    []Review       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"reviews"`
	Cart       Cart           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cart"`
	Orders     []Order        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"orders"`
	Books      pq.StringArray `gorm:"type:varchar(64)[]" json:"books"`
	IsFlagged  bool           `gorm:"default:false" json:"isFlagged"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"-"`
}

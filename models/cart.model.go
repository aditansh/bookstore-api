package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID       uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	UserID   uuid.UUID  `gorm:"type:uuid;not null" json:"-"`
	Username string     `gorm:"not null" json:"username"`
	Items    []CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
	Value    float64    `gorm:"not null; default:0" json:"value"`
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"-"`
	CartID    uuid.UUID `gorm:"type:uuid;primary_key;not null" json:"-"`
	BookID    uuid.UUID `gorm:"type:uuid;primary_key;not null" json:"-"`
	BookName  string    `gorm:"type:varchar;not null" json:"bookName"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	VendorID  uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

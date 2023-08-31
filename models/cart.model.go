package models

import "github.com/google/uuid"

type Cart struct {
	ID     string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	UserID string     `gorm:"type:uuid;not null" json:"-"`
	Items  []CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
	Value  float64    `gorm:"not null; default:0" json:"value"`
}

type CartItem struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	CartID    uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	BookID    uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	VendorID  uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	CreatedAt string    `gorm:"autoCreateTime" json:"-"`
	UpdatedAt string    `gorm:"autoUpdateTime" json:"-"`
}

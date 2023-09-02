package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID   `gorm:"type:uuid;not null" json:"-"`
	Items     []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
	Address   string      `gorm:"not null" json:"address"`
	Value     float64     `gorm:"not null; default:0" json:"value"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"-"`
}

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	BookID     uuid.UUID `gorm:"type:uuid;not null" json:"book_id"`
	BookName   string    `gorm:"not null" json:"book_name"`
	Author     string    `gorm:"not null" json:"author"`
	VendorName string    `gorm:"not null" json:"vendor_name"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"-"`
}

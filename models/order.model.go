package models

import "github.com/google/uuid"

type Order struct {
	ID        string      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	UserID    string      `gorm:"type:uuid;not null" json:"-"`
	Items     []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
	Address   string      `gorm:"not null" json:"address"`
	Value     float64     `gorm:"not null; default:0" json:"value"`
	CreatedAt string      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt string      `gorm:"autoUpdateTime" json:"-"`
}

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	BookName   string    `gorm:"not null" json:"book_name"`
	Author     string    `gorm:"not null" json:"author"`
	VendorName string    `gorm:"not null" json:"vendor_name"`
	CreatedAt  string    `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  string    `gorm:"autoUpdateTime" json:"-"`
}

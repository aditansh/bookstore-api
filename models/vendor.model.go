package models

import (
	"time"

	"github.com/google/uuid"
)

type Vendor struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	Name       string     `gorm:"not null" json:"name"`
	Email      string     `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string     `gorm:"not null" json:"-"`
	IsVerified bool       `gorm:"default:false" json:"-"`
	IsActive   bool       `gorm:"default:true" json:"isActive"`
	IsDeleted  bool       `gorm:"default:false" json:"-"`
	Books      []Book     `gorm:"many2many:book_vendor" json:"books"`
	CartItems  []CartItem `gorm:"foreignKey:VendorID" json:"cart_items"`
	IsApproved bool       `gorm:"default:false" json:"-"`
	IsFlagged  bool       `gorm:"default:false" json:"isFlagged"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"-"`
}

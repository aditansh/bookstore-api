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
	IsDeleted  bool       `gorm:"default:false" json:"-"`
	Role       string     `gorm:"default:'user'" json:"-"`
	Books      []Book     `gorm:"many2many:book_vendor" json:"books"`
	CartItems  []CartItem `gorm:"foreignKey:VendorID" json:"cart_items"`
	IsApproved bool       `gorm:"default:false" json:"-"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"-"`
}

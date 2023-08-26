package models

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	Name        string     `gorm:"not null" json:"name"`
	Author      string     `gorm:"not null" json:"author"`
	Description string     `gorm:"not null" json:"description"`
	Categories  []string   `gorm:"not null" json:"categories"`
	Price       float64    `gorm:"not null" json:"price"`
	Stock       int        `gorm:"not null" json:"stock"`
	Rating      float64    `gorm:"not null; default:0" json:"rating"`
	Reviews     []Review   `gorm:"foreignKey:BookID" json:"reviews"`
	InCart      []CartItem `gorm:"foreignKey:BookID" json:"inCart"`
	Vendors     []Vendor   `gorm:"many2many:book_vendor" json:"vendors"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"-"`
}

type Review struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	UserID    string    `gorm:"type:uuid;not null" json:"-"`
	BookID    string    `gorm:"type:uuid;not null" json:"-"`
	Comment   string    `gorm:"not null" json:"comment"`
	Rating    int       `gorm:"not null" json:"rating"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

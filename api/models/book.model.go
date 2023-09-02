package models

import (
	"time"

	"github.com/google/uuid"
	pq "github.com/lib/pq"
)

type Book struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Author      string         `gorm:"not null" json:"author"`
	Description string         `gorm:"not null" json:"description"`
	Categories  pq.StringArray `gorm:"type:varchar(64)[];not null" json:"categories"`
	Stock       int            `gorm:"not null" json:"stock"`
	Reviews     []Review       `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews"`
	InCart      []CartItem     `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"inCart"`
	VendorID    uuid.UUID      `gorm:"type:uuid;not null" json:"-"`
	Cost        float64        `gorm:"not null" json:"cost"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"-"`
}

type Review struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"-"`
	BookID    uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Username  string    `gorm:"type:varchar;default:User Deleted" json:"username"`
	Comment   string    `gorm:"not null" json:"comment"`
	Rating    int       `gorm:"not null" json:"rating"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

package models

type Order struct {
	ID        string      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	UserID    string      `gorm:"type:uuid;not null" json:"-"`
	Items     []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Address   string      `gorm:"not null" json:"address"`
	CreatedAt string      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt string      `gorm:"autoUpdateTime" json:"-"`
}

type OrderItem struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	OrderID   string `gorm:"type:uuid;not null" json:"-"`
	BookID    string `gorm:"type:uuid;not null" json:"-"`
	Quantity  int    `gorm:"not null" json:"quantity"`
	CreatedAt string `gorm:"autoCreateTime" json:"-"`
	UpdatedAt string `gorm:"autoUpdateTime" json:"-"`
}

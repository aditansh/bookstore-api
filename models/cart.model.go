package models

type Cart struct {
	ID     string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	UserID string     `gorm:"type:uuid;not null" json:"-"`
	Items  []CartItem `gorm:"foreignKey:CartID" json:"items"`
	Value  float64    `gorm:"not null; default:0" json:"value"`
}

type CartItem struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	CartID    string `gorm:"type:uuid;not null" json:"-"`
	BookID    string `gorm:"type:uuid;not null" json:"-"`
	Quantity  int    `gorm:"not null" json:"quantity"`
	VendorID  string `gorm:"type:uuid;not null" json:"-"`
	CreatedAt string `gorm:"autoCreateTime" json:"-"`
	UpdatedAt string `gorm:"autoUpdateTime" json:"-"`
}

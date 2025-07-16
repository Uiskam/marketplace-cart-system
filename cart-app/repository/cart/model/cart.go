package model

import (
	"time"
)

const (
	AddProductEventType    = "added"
	RemoveProductEventType = "removed"
	CheckoutEventType      = "checkout"
)

const (
	CartStateActive     = "active"
	CartStateCheckedOut = "checked_out"
)

type CartEvent struct {
	SeqNumber   uint64    `gorm:"primaryKey;autoIncrement:false" json:"seq_number"`
	UUID        string    `gorm:"primaryKey" json:"uuid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProductUUID string    `json:"product_uuid"`
	ProductName string    `json:"product_name"`
	Price       int       `json:"price"`
	EventType   string    `json:"event_type"` // e.g., "added", "removed"
}

type Cart struct {
	CartUUID  string    `gorm:"primaryKey;index" json:"cart_uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartDTO struct {
	CartUUID string   `json:"cart_uuid"`
	Products []string `json:"products"`
	Total    int      `json:"total"`
	State    string   `json:"state"` // e.g., "active", "checked_out", "expired"
}

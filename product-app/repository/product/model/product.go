package model

import (
	"time"
)

type Product struct {
	UUID  string `gorm:"primarykey" json:"uuid"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductLock struct {
	SeqNumber     uint64    `gorm:"primaryKey;autoIncrement:false" json:"seq_number"`
	ProductUUID   string    `gorm:"primaryKey;foreignKey:ProductUUID;references:UUID" json:"product_uuid"`
	LockingEntity string    `json:"locking_entity"`
	LockEnd       time.Time `json:"lock_end"`
}

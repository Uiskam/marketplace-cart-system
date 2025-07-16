package product

import (
	"gorm.io/gorm"
	"product-app/repository/product/model"
)

type ReadRepository struct {
	db *gorm.DB
}

func NewReadRepository(db *gorm.DB) *ReadRepository {
	return &ReadRepository{db: db}
}

func (r *ReadRepository) GetProducts() ([]model.Product, error) {
	var products []model.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

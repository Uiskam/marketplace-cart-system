package getProduct

import (
	"product-app/repository/product"
	"product-app/repository/product/model"
)

type Handler struct {
	repo *product.ReadRepository
}

func NewHandler(readRepo *product.ReadRepository) *Handler {
	return &Handler{
		repo: readRepo,
	}
}

func (h *Handler) Handle(query interface{}) (any, error) {
	result, err := h.repo.GetProducts()
	if err != nil {
		return []model.Product{}, err
	}
	return result, nil
}

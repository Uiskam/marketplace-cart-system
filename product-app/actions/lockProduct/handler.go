package lockProduct

import (
	"product-app/repository/product"
)

type Handler struct {
	repo *product.WriteRepository
}

func NewHandler(writeRepo *product.WriteRepository) *Handler {
	return &Handler{
		repo: writeRepo,
	}
}

func (h *Handler) Handle(command interface{}) (any, error) {
	cmd := command.(Command)

	price, err := h.repo.Lock(cmd.ProductUUID, cmd.LockingEntity)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}, err
	}

	return map[string]interface{}{
		"message": "Ok",
		"price":   price,
	}, nil
}

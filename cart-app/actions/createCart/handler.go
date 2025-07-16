package createCart

import (
	"cart-app/repository/cart"
)

type Handler struct {
	repo *cart.WriteRepository
}

func NewHandler(writeRepo *cart.WriteRepository) *Handler {
	return &Handler{
		repo: writeRepo,
	}
}

func (h *Handler) Handle(command interface{}) (any, error) {
	cartUUID, err := h.repo.CreateCart()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}, err
	}

	return map[string]interface{}{
		"message":   "Ok",
		"cart_uuid": cartUUID,
	}, nil
}

package removeFromCart

import (
	"cart-app/external"
	"cart-app/repository/cart"
	"fmt"
)

type Handler struct {
	repo          *cart.WriteRepository
	productClient *external.ProductClient
}

func NewHandler(writeRepo *cart.WriteRepository, productClient *external.ProductClient) *Handler {
	return &Handler{
		repo:          writeRepo,
		productClient: productClient,
	}
}

func (h *Handler) Handle(command interface{}) (any, error) {
	cmd := command.(Command)

	// Unlock the product first
	err := h.productClient.UnlockProduct(cmd.ProductUUID, cmd.CartUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock product: %w", err)
	}

	// Remove product from cart
	err = h.repo.RemoveProduct(cmd.CartUUID, cmd.ProductUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove product from cart: %w", err)
	}

	return map[string]interface{}{
		"message": "Product removed from cart",
	}, nil
}

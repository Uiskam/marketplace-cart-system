package addToCart

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

	// Lock the product first
	price, err := h.productClient.LockProduct(cmd.ProductUUID, cmd.CartUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to lock product: %w", err)
	}

	// Add product to cart
	err = h.repo.AddProduct(cmd.CartUUID, cmd.ProductUUID, price)
	if err != nil {
		// Unlock the product if adding to cart fails
		err := h.productClient.UnlockProduct(cmd.ProductUUID, cmd.CartUUID)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to add product to cart: %w", err)
	}

	return map[string]interface{}{
		"message": "Product added to cart",
		"price":   price,
	}, nil
}

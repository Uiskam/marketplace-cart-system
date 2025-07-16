package checkoutCart

import (
	"cart-app/actions/getCart"
	"cart-app/app/cqrs"
	"cart-app/external"
	"cart-app/repository/cart"
	"cart-app/repository/cart/model"
	"fmt"
	"log"
)

type Handler struct {
	repo               *cart.WriteRepository
	queryBus           *cqrs.Bus
	productClient      *external.ProductClient
	notificationClient *external.NotificationClient
}

func NewHandler(writeRepo *cart.WriteRepository, queryBus *cqrs.Bus, productClient *external.ProductClient, notificationClient *external.NotificationClient) *Handler {
	return &Handler{
		repo:               writeRepo,
		queryBus:           queryBus,
		productClient:      productClient,
		notificationClient: notificationClient,
	}
}

func (h *Handler) Handle(command interface{}) (any, error) {
	cmd := command.(Command)

	// Get cart to find all products that need to be sold
	cartQuery := getCart.Query{CartUUID: cmd.CartUUID}
	result, err := h.queryBus.Send(cartQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	cartDTO := result.(*model.CartDTO)

	// Sell all products in the cart
	for _, productUUID := range cartDTO.Products {
		err := h.productClient.SellProduct(productUUID, cmd.CartUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to sell product %s: %w", productUUID, err)
		}
	}

	// Mark cart as checked out
	err = h.repo.CheckoutCart(cmd.CartUUID)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}, err
	}

	// Notify the user about the successful checkout
	err = h.notificationClient.SendTask("email", external.NewEmailNotification("Checkout Successful", "Your cart has been successfully checked out."))
	if err != nil {
		log.Println(err)
	}
	return map[string]interface{}{
		"message": "Cart checked out successfully",
		"total":   cartDTO.Total,
	}, nil
}

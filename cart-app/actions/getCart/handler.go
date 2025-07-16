package getCart

import (
	cart_repo "cart-app/repository/cart"
	"cart-app/repository/cart/model"
)

type Handler struct {
	repo *cart_repo.ReadRepository
}

func NewHandler(readRepo *cart_repo.ReadRepository) *Handler {
	return &Handler{
		repo: readRepo,
	}
}

func (h *Handler) Handle(query interface{}) (any, error) {
	q := query.(Query)

	// First, get the cart to get the username
	cart, err := h.repo.GetCart(q.CartUUID)
	if err != nil {
		return nil, err
	}

	// Get all events for this cart
	events, err := h.repo.GetDefiningCartEvents(q.CartUUID)
	if err != nil {
		return nil, err
	}

	// Build CartDTO from events
	cartDTO := h.buildCartDTO(cart, events)
	return cartDTO, nil
}

func (h *Handler) buildCartDTO(cart *model.Cart, events []model.CartEvent) *model.CartDTO {
	productMap := make(map[string]bool)
	total := 0
	isCheckedOut := false

	// Process events to build current cart state
outer:
	for _, event := range events {
		switch event.EventType {
		case model.AddProductEventType:
			productMap[event.ProductUUID] = true
			total += event.Price
		case model.RemoveProductEventType:
			continue
		case model.CheckoutEventType:
			// Cart was checked out
			isCheckedOut = true
			productMap = make(map[string]bool)
			break outer // Exit the loop immediately

		}
	}

	// Determine cart state
	var state string
	if isCheckedOut {
		state = model.CartStateCheckedOut
	} else {
		state = model.CartStateActive
	}

	// Convert product map to slice
	products := make([]string, 0, len(productMap))
	for productUUID := range productMap {
		products = append(products, productUUID)
	}

	return &model.CartDTO{
		CartUUID: cart.CartUUID,
		Products: products,
		Total:    total,
		State:    state,
	}
}

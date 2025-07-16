package cart

import (
	"cart-app/repository/cart/model"
	"fmt"
	"gorm.io/gorm"
)

type ReadRepository struct {
	db *gorm.DB
}

func NewReadRepository(db *gorm.DB) *ReadRepository {
	return &ReadRepository{db: db}
}

func (r *ReadRepository) GetCart(cartUUID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("cart_uuid = ?", cartUUID).First(&cart).Error
	if err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}
	return &cart, nil
}

func (r *ReadRepository) GetCartEvents(cartUUID string) ([]model.CartEvent, error) {
	var events []model.CartEvent

	query := r.db.Where("uuid = ?", cartUUID)

	err := query.Order("seq_number ASC").Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get cart events: %w", err)
	}

	return events, nil
}

func (r *ReadRepository) GetDefiningCartEvents(cartUUID string) ([]model.CartEvent, error) {
	var events []model.CartEvent

	query := r.db.Raw(`
    SELECT DISTINCT ON (product_uuid)
    *
	FROM cart_events
	WHERE uuid = ?
	ORDER BY product_uuid, seq_number DESC;
	`, cartUUID)
	err := query.Scan(&events).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get cart events: %w", err)
	}

	return events, nil
}

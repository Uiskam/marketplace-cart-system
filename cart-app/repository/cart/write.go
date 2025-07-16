package cart

import (
	"cart-app/external"
	"cart-app/repository/cart/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WriteRepository struct {
	db *gorm.DB
}

func NewWriteRepository(db *gorm.DB) *WriteRepository {
	return &WriteRepository{db: db}
}

// isCheckedOut checks if a checkout event exists for the given cart
func (r *WriteRepository) isCheckedOut(uow *external.UnitOfWork, cartUUID string) bool {
	var count int64
	err := uow.Tx().Model(&model.CartEvent{}).
		Where("uuid = ? AND event_type = ?", cartUUID, model.CheckoutEventType).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

func (r *WriteRepository) AddProduct(cartUUID, productUUID string, price int) error {
	return external.WithRetryTx(r.db, 3, func(uow *external.UnitOfWork) error {
		// Validate cart exists
		if r.isCheckedOut(uow, cartUUID) {
			return errors.New("cart has already been checked out")
		}
		var cart model.Cart
		err := uow.Tx().Where("cart_uuid = ?", cartUUID).First(&cart).Error
		if err != nil {
			return err
		}

		// Get next sequence number
		var maxSeqNumber uint64
		err = uow.Tx().Model(&model.CartEvent{}).Where("uuid = ?", cartUUID).Select("COALESCE(MAX(seq_number), 0)").Scan(&maxSeqNumber).Error
		if err != nil {
			return err
		}

		// Create add product event
		now := time.Now()
		addEvent := model.CartEvent{
			SeqNumber:   maxSeqNumber + 1,
			UUID:        cartUUID,
			CreatedAt:   now,
			UpdatedAt:   now,
			ProductUUID: productUUID,
			ProductName: "",
			Price:       price,
			EventType:   model.AddProductEventType,
		}

		return uow.Tx().Create(&addEvent).Error
	})
}

func (r *WriteRepository) RemoveProduct(cartUUID, productUUID string) error {
	return external.WithRetryTx(r.db, 3, func(uow *external.UnitOfWork) error {
		// Validate cart exists

		if r.isCheckedOut(uow, cartUUID) {
			return errors.New("cart has already been checked out")
		}
		var cart model.Cart
		err := uow.Tx().Where("cart_uuid = ?", cartUUID).First(&cart).Error
		if err != nil {
			return err
		}

		// Check if product is in the cart by getting the latest event for this product
		var latestEvent model.CartEvent
		err = uow.Tx().Where("uuid = ? AND product_uuid = ?", cartUUID, productUUID).
			Order("seq_number DESC").
			First(&latestEvent).Error

		// If no event exists for this product, abort
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound // Product not in cart
			}
			return err
		}

		// If the latest event is a remove event, the product is not in cart, abort
		if latestEvent.EventType == model.RemoveProductEventType {
			return gorm.ErrRecordNotFound // Product already removed
		}

		// Get next sequence number
		var maxSeqNumber uint64
		err = uow.Tx().Model(&model.CartEvent{}).Where("uuid = ?", cartUUID).Select("COALESCE(MAX(seq_number), 0)").Scan(&maxSeqNumber).Error
		if err != nil {
			return err
		}

		// Create remove product event
		now := time.Now()
		removeEvent := model.CartEvent{
			SeqNumber:   maxSeqNumber + 1,
			UUID:        cartUUID,
			CreatedAt:   now,
			UpdatedAt:   now,
			ProductUUID: productUUID,
			ProductName: "",
			Price:       0,
			EventType:   model.RemoveProductEventType,
		}

		return uow.Tx().Create(&removeEvent).Error
	})
}

func (r *WriteRepository) CreateCart() (string, error) {
	var cartUUID string
	err := external.WithRetryTx(r.db, 3, func(uow *external.UnitOfWork) error {
		newUUID := uuid.New().String()
		now := time.Now()

		cart := model.Cart{
			CartUUID:  newUUID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := uow.Tx().Create(&cart).Error
		if err != nil {
			return err
		}

		cartUUID = newUUID
		return nil
	})
	return cartUUID, err
}

func (r *WriteRepository) CheckoutCart(cartUUID string) error {
	return external.WithRetryTx(r.db, 3, func(uow *external.UnitOfWork) error {
		// Validate cart exists
		if r.isCheckedOut(uow, cartUUID) {
			return errors.New("cart has already been checked out")
		}
		var cart model.Cart
		err := uow.Tx().Where("cart_uuid = ?", cartUUID).First(&cart).Error
		if err != nil {
			return err
		}

		// Add checkout event
		var maxSeqNumber uint64
		err = uow.Tx().Model(&model.CartEvent{}).Where("uuid = ?", cartUUID).Select("COALESCE(MAX(seq_number), 0)").Scan(&maxSeqNumber).Error
		if err != nil {
			return err
		}

		now := time.Now()
		checkoutEvent := model.CartEvent{
			SeqNumber:   maxSeqNumber + 1,
			UUID:        cartUUID,
			CreatedAt:   now,
			UpdatedAt:   now,
			ProductUUID: "",
			ProductName: "",
			Price:       0,
			EventType:   model.CheckoutEventType,
		}

		return uow.Tx().Create(&checkoutEvent).Error
	})
}

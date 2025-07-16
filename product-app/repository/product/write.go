package product

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"product-app/external"
	"product-app/repository/product/model"
)

const (
	SellLockEndDate = "2200-01-01"
)

type SellError struct {
	FailedProducts []string
	Message        string
}

func (e *SellError) Error() string {
	return e.Message
}

type WriteRepository struct {
	db *gorm.DB
}

func NewWriteRepository(db *gorm.DB) *WriteRepository {
	return &WriteRepository{db: db}
}

func (repo *WriteRepository) IsLocked(productUUID string) (bool, error) {
	uow := external.BeginTransaction(repo.db)
	defer func() { _ = uow.Rollback() }()

	var count int64
	now := time.Now()
	err := uow.Tx().Model(&model.ProductLock{}).
		Where("product_uuid = ? AND lock_end > ?", productUUID, now).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check lock status: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *WriteRepository) Lock(productUUID, lockingEntity string) (int, error) {
	var price int
	err := external.WithRetryTx(repo.db, 3, func(uow *external.UnitOfWork) error {
		var product model.Product
		err := uow.Tx().Where("uuid = ?", productUUID).First(&product).Error
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}
		price = product.Price

		return repo.createLock(uow, productUUID, lockingEntity, 15*time.Minute)
	})
	return price, err
}

func (repo *WriteRepository) Unlock(productUUID, lockingEntity string) error {
	return external.WithRetryTx(repo.db, 3, func(uow *external.UnitOfWork) error {
		return repo.updateLockEnd(uow, productUUID, lockingEntity, time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC))
	})
}

func (repo *WriteRepository) SellProduct(productUUIDs []string, lockingEntity string) error {
	sellDate, _ := time.Parse("2006-01-02", SellLockEndDate)
	var failedProducts []string
	
	return external.WithRetryTx(repo.db, 3, func(uow *external.UnitOfWork) error {
		failedProducts = nil // Reset for each retry
		
		for _, productUUID := range productUUIDs {
			err := repo.createLock(uow, productUUID, lockingEntity, time.Until(sellDate))
			if err != nil {
				failedProducts = append(failedProducts, productUUID)
			}
		}
		
		if len(failedProducts) > 0 {
			return &SellError{
				FailedProducts: failedProducts,
				Message:        fmt.Sprintf("failed to sell %d products: %v", len(failedProducts), failedProducts),
			}
		}
		
		return nil
	})
}

func (repo *WriteRepository) createLock(uow *external.UnitOfWork, productUUID, lockingEntity string, duration time.Duration) error {
	now := time.Now()
	sellDate, _ := time.Parse("2006-01-02", SellLockEndDate)

	// Check if product is already sold (has lock until sell date)
	var soldCount int64
	err := uow.Tx().Model(&model.ProductLock{}).
		Where("product_uuid = ? AND lock_end >= ?", productUUID, sellDate.Add(-24*time.Hour)).
		Count(&soldCount).Error
	if err != nil {
		return fmt.Errorf("failed to check if product is sold: %w", err)
	}
	if soldCount > 0 {
		return fmt.Errorf("product is already sold")
	}

	// Check for existing valid lock by another entity
	var count int64
	err = uow.Tx().Model(&model.ProductLock{}).
		Where("product_uuid = ? AND lock_end > ? AND locking_entity != ?", productUUID, now, lockingEntity).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check for conflicting locks: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("product is locked by another entity")
	}

	// Get next sequence number with locking
	var maxSeq uint64
	err = uow.Tx().Model(&model.ProductLock{}).
		Where("product_uuid = ?", productUUID).
		Select("COALESCE(MAX(seq_number), 0)").
		Scan(&maxSeq).Error
	if err != nil {
		return fmt.Errorf("failed to get max sequence number: %w", err)
	}

	lock := model.ProductLock{
		SeqNumber:     maxSeq + 1,
		ProductUUID:   productUUID,
		LockingEntity: lockingEntity,
		LockEnd:       now.Add(duration),
	}

	if err := uow.Tx().Create(&lock).Error; err != nil {
		return fmt.Errorf("failed to create lock: %w", err)
	}
	return nil
}

func (repo *WriteRepository) updateLockEnd(uow *external.UnitOfWork, productUUID, lockingEntity string, newEnd time.Time) error {
	now := time.Now()
	var lock model.ProductLock

	err := uow.Tx().Where("product_uuid = ? AND locking_entity = ?", productUUID, lockingEntity).
		Order("seq_number DESC").First(&lock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("no lock found for product %s by entity %s", productUUID, lockingEntity)
	} else if err != nil {
		return fmt.Errorf("failed to fetch lock: %w", err)
	}

	if lock.LockEnd.Before(now) {
		return fmt.Errorf("lock has already expired")
	}

	res := uow.Tx().Model(&model.ProductLock{}).
		Where("product_uuid = ? AND seq_number = ? AND locking_entity = ?", productUUID, lock.SeqNumber, lockingEntity).
		Update("lock_end", newEnd)

	if res.Error != nil {
		return fmt.Errorf("failed to update lock_end: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("lock was modified by another process")
	}
	return nil
}

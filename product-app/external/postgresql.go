package external

import (
	"fmt"
	"log"
	"os"
	"product-app/repository/product/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}

func NewPostgresConnection() (*gorm.DB, error) {
	host, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	port, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	user, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	password, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbname, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&model.Product{}, &model.ProductLock{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// Check if products table is empty and execute base.sql if so
	var count int64
	err = db.Model(&model.Product{}).Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	if count == 0 {
		log.Println("Products table is empty, executing base.sql...")
		err = executeSQLFile(db, "external/base.sql")
		if err != nil {
			return nil, fmt.Errorf("failed to execute base.sql: %w", err)
		}
		log.Println("Successfully executed base.sql")
	}

	return db, nil
}

// executeSQLFile reads and executes SQL statements from a file
func executeSQLFile(db *gorm.DB, filePath string) error {
	// Read the SQL file
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", filePath, err)
	}

	// Execute the SQL statements
	sqlContent := string(sqlBytes)
	err = db.Exec(sqlContent).Error
	if err != nil {
		return fmt.Errorf("failed to execute SQL from %s: %w", filePath, err)
	}

	return nil
}

// UnitOfWork represents a database transaction wrapper
type UnitOfWork struct {
	tx        *gorm.DB
	committed bool
}

// BeginTransaction creates a new Unit of Work with a database transaction
func BeginTransaction(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{tx: db.Begin()}
}

// Commit commits the transaction
func (uow *UnitOfWork) Commit() error {
	err := uow.tx.Commit().Error
	if err == nil {
		uow.committed = true
	}
	return err
}

// Rollback rolls back the transaction if not already committed
func (uow *UnitOfWork) Rollback() error {
	if !uow.committed {
		return uow.tx.Rollback().Error
	}
	return nil
}

// Tx returns the transaction database instance
func (uow *UnitOfWork) Tx() *gorm.DB {
	return uow.tx
}

// WithRetryTx executes a function with retry logic using Unit of Work pattern
func WithRetryTx(db *gorm.DB, attempts int, fn func(uow *UnitOfWork) error) error {
	var err error
	for i := 0; i < attempts; i++ {
		uow := BeginTransaction(db)
		err = fn(uow)
		if err == nil {
			if commitErr := uow.Commit(); commitErr == nil {
				return nil
			} else {
				err = commitErr
			}
		}
		_ = uow.Rollback()
		time.Sleep(50 * time.Millisecond)
	}
	return err
}

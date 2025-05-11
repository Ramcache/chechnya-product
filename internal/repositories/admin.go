package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// AdminRepoInterface интерфейс для админских операций с БД
type AdminRepoInterface interface {
	TruncateTable(tableName string) error
	TruncateAllTables() error
}

// AdminRepo реализация AdminRepoInterface
type AdminRepo struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) TruncateTable(tableName string) error {
	allowedTables := map[string]bool{
		"announcements": true,
		"cart_items":    true,
		"categories":    true,
		"products":      true,
		"orders":        true,
		"order_items":   true,
		"users":         true,
		"reviews":       true,
	}

	if !allowedTables[tableName] {
		return fmt.Errorf("недопустимая таблица: %s", tableName)
	}

	_, err := r.db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE`, tableName))
	return err
}

func (r *AdminRepo) TruncateAllTables() error {
	tables := []string{"announcements", "cart_items", "categories", "order_items", "orders", "products", "reviews"}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE`, table)); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

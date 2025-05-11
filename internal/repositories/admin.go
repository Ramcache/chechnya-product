package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// AdminRepoInterface интерфейс для админских операций с БД
type AdminRepoInterface interface {
	TruncateTable(tableName string) error
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
		"carts":         true,
		"reviews":       true,
	}

	if !allowedTables[tableName] {
		return fmt.Errorf("недопустимая таблица: %s", tableName)
	}

	_, err := r.db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE`, tableName))
	return err
}

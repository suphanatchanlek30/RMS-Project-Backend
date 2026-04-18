package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type PaymentMethodRepository struct {
	DB *pgxpool.Pool
}

func NewPaymentMethodRepository(db *pgxpool.Pool) *PaymentMethodRepository {
	return &PaymentMethodRepository{DB: db}
}

func (r *PaymentMethodRepository) GetAll(ctx context.Context) ([]models.PaymentMethodItem, error) {
	query := `
		SELECT payment_method_id, method_name
		FROM payment_methods
		ORDER BY payment_method_id ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	items := make([]models.PaymentMethodItem, 0)
	for rows.Next() {
		var item models.PaymentMethodItem
		if err := rows.Scan(&item.PaymentMethodID, &item.MethodName); err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return items, nil
}

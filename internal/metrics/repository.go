package metrics

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Overview(ctx context.Context) (Overview, error) {
	var out Overview
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE((SELECT SUM(quantity * COALESCE(unit_cost_cents, 0)) FROM stock_movements WHERE movement_type = 'entry'), 0),
			COALESCE((SELECT SUM(subtotal_cents) FROM orders WHERE status <> 'canceled'), 0),
			COALESCE((SELECT SUM(quantity * COALESCE(unit_cost_cents, 0)) FROM stock_movements WHERE movement_type = 'sale'), 0),
			COALESCE((SELECT COUNT(*) FROM orders WHERE status <> 'canceled'), 0)
	`).Scan(&out.EntriesCents, &out.SalesCents, &out.CostOfGoodsCents, &out.OrdersCount)
	if err != nil {
		return Overview{}, err
	}
	out.GrossProfitCents = out.SalesCents - out.CostOfGoodsCents
	out.NetProfitCents = out.GrossProfitCents
	return out, nil
}

package products

import (
	"context"
	"database/sql"
)

type Repository interface {
	List(ctx context.Context, category string, onlyActive bool) ([]Product, error)
	Get(ctx context.Context, id string) (Product, error)
	Create(ctx context.Context, req CreateProductRequest) (Product, error)
	Update(ctx context.Context, id string, req UpdateProductRequest) (Product, error)
	Delete(ctx context.Context, id string) error
	SetAvailability(ctx context.Context, id string, isActive bool) error
	AddStockMovement(ctx context.Context, productID string, req StockMovementRequest) error
	ReserveStock(ctx context.Context, tx *sql.Tx, productID string, quantity int) (Product, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, category string, onlyActive bool) ([]Product, error) {
	query := `
		SELECT id::text, COALESCE(sku, ''), name, COALESCE(description, ''), category, COALESCE(image_url, ''),
		       price_cents, cost_cents, stock_quantity, min_stock_quantity, is_active, created_at, updated_at
		FROM products
		WHERE ($1 = '' OR category = $1) AND ($2 = false OR is_active = true)
		ORDER BY category, name`
	rows, err := r.db.QueryContext(ctx, query, category, onlyActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Product, 0)
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.ImageURL, &p.PriceCents, &p.CostCents, &p.StockQuantity, &p.MinStockQuantity, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *repository) Get(ctx context.Context, id string) (Product, error) {
	var p Product
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, COALESCE(sku, ''), name, COALESCE(description, ''), category, COALESCE(image_url, ''),
		       price_cents, cost_cents, stock_quantity, min_stock_quantity, is_active, created_at, updated_at
		FROM products
		WHERE id = $1`, id).
		Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.ImageURL, &p.PriceCents, &p.CostCents, &p.StockQuantity, &p.MinStockQuantity, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *repository) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	var p Product
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO products (sku, name, description, category, image_url, price_cents, cost_cents, stock_quantity, min_stock_quantity)
		VALUES (NULLIF($1, ''), $2, NULLIF($3, ''), $4, NULLIF($5, ''), $6, $7, $8, $9)
		RETURNING id::text, COALESCE(sku, ''), name, COALESCE(description, ''), category, COALESCE(image_url, ''),
		          price_cents, cost_cents, stock_quantity, min_stock_quantity, is_active, created_at, updated_at`,
		req.SKU, req.Name, req.Description, req.Category, req.ImageURL, req.PriceCents, req.CostCents, req.StockQuantity, req.MinStockQuantity).
		Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.ImageURL, &p.PriceCents, &p.CostCents, &p.StockQuantity, &p.MinStockQuantity, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *repository) Update(ctx context.Context, id string, req UpdateProductRequest) (Product, error) {
	var p Product
	err := r.db.QueryRowContext(ctx, `
		UPDATE products
		SET sku = NULLIF($2, ''), name = $3, description = NULLIF($4, ''), category = $5, image_url = NULLIF($6, ''),
		    price_cents = $7, cost_cents = $8, stock_quantity = $9, min_stock_quantity = $10, is_active = $11, updated_at = NOW()
		WHERE id = $1
		RETURNING id::text, COALESCE(sku, ''), name, COALESCE(description, ''), category, COALESCE(image_url, ''),
		          price_cents, cost_cents, stock_quantity, min_stock_quantity, is_active, created_at, updated_at`,
		id, req.SKU, req.Name, req.Description, req.Category, req.ImageURL, req.PriceCents, req.CostCents, req.StockQuantity, req.MinStockQuantity, req.IsActive).
		Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.ImageURL, &p.PriceCents, &p.CostCents, &p.StockQuantity, &p.MinStockQuantity, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (r *repository) SetAvailability(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE products SET is_active = $2, updated_at = NOW() WHERE id = $1`, id, isActive)
	return err
}

func (r *repository) AddStockMovement(ctx context.Context, productID string, req StockMovementRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	unitCost := req.UnitCostCents
	if req.TotalCostCents > 0 && req.Quantity > 0 {
		unitCost = req.TotalCostCents / req.Quantity
	}

	delta := req.Quantity
	if req.Type == "loss" || req.Type == "sale" {
		delta = -req.Quantity
	}

	if req.Type == "entry" && unitCost > 0 {
		if _, err := tx.ExecContext(ctx, `
			UPDATE products
			SET cost_cents = CASE
					WHEN stock_quantity + $2 <= 0 THEN $3
					ELSE ((cost_cents * stock_quantity) + ($3 * $2)) / (stock_quantity + $2)
				END,
				stock_quantity = stock_quantity + $2,
				updated_at = NOW()
			WHERE id = $1`, productID, req.Quantity, unitCost); err != nil {
			return err
		}
	} else if _, err := tx.ExecContext(ctx, `UPDATE products SET stock_quantity = stock_quantity + $2, updated_at = NOW() WHERE id = $1`, productID, delta); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO stock_movements (product_id, movement_type, quantity, unit_cost_cents, notes)
		VALUES ($1, $2, $3, NULLIF($4, 0), NULLIF($5, ''))`, productID, req.Type, req.Quantity, unitCost, req.Notes); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *repository) ReserveStock(ctx context.Context, tx *sql.Tx, productID string, quantity int) (Product, error) {
	var p Product
	err := tx.QueryRowContext(ctx, `
		UPDATE products
		SET stock_quantity = stock_quantity - $2, updated_at = NOW()
		WHERE id = $1 AND is_active = true AND stock_quantity >= $2
		RETURNING id::text, COALESCE(sku, ''), name, COALESCE(description, ''), category, COALESCE(image_url, ''),
		          price_cents, cost_cents, stock_quantity, min_stock_quantity, is_active, created_at, updated_at`,
		productID, quantity).
		Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.ImageURL, &p.PriceCents, &p.CostCents, &p.StockQuantity, &p.MinStockQuantity, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return Product{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO stock_movements (product_id, movement_type, quantity, unit_cost_cents, notes)
		VALUES ($1, 'sale', $2, NULLIF($3, 0), 'Reserva criada por pedido')`, productID, quantity, p.CostCents)
	return p, err
}

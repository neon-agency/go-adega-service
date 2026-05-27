package products

import "time"

type Product struct {
	ID               string    `json:"id"`
	SKU              string    `json:"sku,omitempty"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	Category         string    `json:"category"`
	ImageURL         string    `json:"image_url,omitempty"`
	PriceCents       int       `json:"price_cents"`
	CostCents        int       `json:"cost_cents"`
	StockQuantity    int       `json:"stock_quantity"`
	MinStockQuantity int       `json:"min_stock_quantity"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	SKU              string `json:"sku"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Category         string `json:"category"`
	ImageURL         string `json:"image_url"`
	PriceCents       int    `json:"price_cents"`
	CostCents        int    `json:"cost_cents"`
	StockQuantity    int    `json:"stock_quantity"`
	MinStockQuantity int    `json:"min_stock_quantity"`
}

type UpdateProductRequest struct {
	SKU              string `json:"sku"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Category         string `json:"category"`
	ImageURL         string `json:"image_url"`
	PriceCents       int    `json:"price_cents"`
	CostCents        int    `json:"cost_cents"`
	StockQuantity    int    `json:"stock_quantity"`
	MinStockQuantity int    `json:"min_stock_quantity"`
	IsActive         bool   `json:"is_active"`
}

type StockMovementRequest struct {
	Type           string `json:"type"`
	Quantity       int    `json:"quantity"`
	UnitCostCents  int    `json:"unit_cost_cents"`
	TotalCostCents int    `json:"total_cost_cents"`
	Notes          string `json:"notes"`
}

type AvailabilityRequest struct {
	On bool `json:"on"`
}

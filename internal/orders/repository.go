package orders

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	CreateCustomer(ctx context.Context, tx *sql.Tx, req CustomerRequest) (string, error)
	CreateAddress(ctx context.Context, tx *sql.Tx, customerID string, req AddressRequest) (string, error)
	ReserveProduct(ctx context.Context, tx *sql.Tx, productID string, quantity int) (OrderItem, error)
	CreateOrder(ctx context.Context, tx *sql.Tx, customerID, addressID string, req CreateOrderRequest, subtotal int, payment PaymentInfo) (string, error)
	AddOrderItem(ctx context.Context, tx *sql.Tx, orderID string, item OrderItem) error
	List(ctx context.Context, status string) ([]Order, error)
	Get(ctx context.Context, id string) (Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Begin(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *repository) CreateCustomer(ctx context.Context, tx *sql.Tx, req CustomerRequest) (string, error) {
	id := uuid.NewString()
	err := tx.QueryRowContext(ctx, `
		INSERT INTO customers (id, name, phone, email, document)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''))
		RETURNING id::text`, id, req.Name, req.Phone, req.Email, req.Document).Scan(&id)
	return id, err
}

func (r *repository) CreateAddress(ctx context.Context, tx *sql.Tx, customerID string, req AddressRequest) (string, error) {
	id := uuid.NewString()
	err := tx.QueryRowContext(ctx, `
		INSERT INTO customer_addresses
			(id, customer_id, street, number, complement, neighborhood, city, state, zip_code, latitude, longitude)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7, $8, $9, NULLIF($10, 0), NULLIF($11, 0))
		RETURNING id::text`,
		id, customerID, req.Street, req.Number, req.Complement, req.Neighborhood, req.City, req.State, req.ZipCode, req.Latitude, req.Longitude).Scan(&id)
	return id, err
}

func (r *repository) ReserveProduct(ctx context.Context, tx *sql.Tx, productID string, quantity int) (OrderItem, error) {
	var item OrderItem
	err := tx.QueryRowContext(ctx, `
		UPDATE products
		SET stock_quantity = stock_quantity - $2, updated_at = NOW()
		WHERE id = $1 AND is_active = true AND stock_quantity >= $2
		RETURNING id::text, name, price_cents`, productID, quantity).
		Scan(&item.ProductID, &item.ProductName, &item.UnitPriceCents)
	if errors.Is(err, sql.ErrNoRows) {
		return OrderItem{}, errors.New("produto indisponível ou estoque insuficiente")
	}
	if err != nil {
		return OrderItem{}, err
	}
	item.ID = uuid.NewString()
	item.Quantity = quantity
	item.TotalCents = item.UnitPriceCents * quantity
	return item, nil
}

func (r *repository) CreateOrder(ctx context.Context, tx *sql.Tx, customerID, addressID string, req CreateOrderRequest, subtotal int, payment PaymentInfo) (string, error) {
	id := uuid.NewString()
	status := "awaiting_payment"
	paymentStatus := "pending"
	if req.PaymentMode == "delivery" {
		status = "separating"
	}
	if payment.Status == "approved" || payment.Status == "paid" {
		status = "paid"
		paymentStatus = "approved"
	}

	err := tx.QueryRowContext(ctx, `
		INSERT INTO orders
			(id, customer_id, delivery_address_id, status, payment_method, payment_mode, payment_status, payment_provider, payment_reference,
			 subtotal_cents, delivery_fee_cents, total_cents, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, ''), NULLIF($9, ''), $10, $11, $12, NULLIF($13, ''))
		RETURNING id::text`,
		id, customerID, addressID, status, req.PaymentMethod, req.PaymentMode, paymentStatus, payment.Provider, payment.Reference,
		subtotal, req.DeliveryFee, subtotal+req.DeliveryFee, req.Notes).Scan(&id)
	return id, err
}

func (r *repository) AddOrderItem(ctx context.Context, tx *sql.Tx, orderID string, item OrderItem) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO order_items (id, order_id, product_id, product_name, quantity, unit_price_cents, total_cents)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		item.ID, orderID, item.ProductID, item.ProductName, item.Quantity, item.UnitPriceCents, item.TotalCents)
	return err
}

func (r *repository) List(ctx context.Context, status string) ([]Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT o.id::text, o.customer_id::text, c.name, c.phone,
		       CONCAT_WS(', ', a.street, NULLIF(a.number, ''), NULLIF(a.neighborhood, ''), NULLIF(a.city, '')) AS delivery_address,
		       o.status::text, o.subtotal_cents, o.delivery_fee_cents,
		       o.total_cents, COALESCE(o.notes, ''), o.payment_method::text, o.payment_mode, o.payment_status::text,
		       COALESCE(o.payment_provider, ''), COALESCE(o.payment_reference, ''), o.created_at, o.updated_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN customer_addresses a ON a.id = o.delivery_address_id
		WHERE ($1 = '' OR o.status::text = $1)
		ORDER BY o.created_at DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Order, 0)
	for rows.Next() {
		var o Order
		if err := scanOrder(rows, &o); err != nil {
			return nil, err
		}
		if err := r.loadItems(ctx, &o); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func (r *repository) Get(ctx context.Context, id string) (Order, error) {
	var o Order
	err := r.db.QueryRowContext(ctx, `
		SELECT o.id::text, o.customer_id::text, c.name, c.phone,
		       CONCAT_WS(', ', a.street, NULLIF(a.number, ''), NULLIF(a.neighborhood, ''), NULLIF(a.city, '')) AS delivery_address,
		       o.status::text, o.subtotal_cents, o.delivery_fee_cents,
		       o.total_cents, COALESCE(o.notes, ''), o.payment_method::text, o.payment_mode, o.payment_status::text,
		       COALESCE(o.payment_provider, ''), COALESCE(o.payment_reference, ''), o.created_at, o.updated_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN customer_addresses a ON a.id = o.delivery_address_id
	WHERE o.id = $1`, id).Scan(
		&o.ID, &o.CustomerID, &o.CustomerName, &o.CustomerPhone, &o.DeliveryAddress, &o.Status, &o.SubtotalCents, &o.DeliveryFeeCents,
		&o.TotalCents, &o.Notes, &o.Payment.PaymentMethod, &o.Payment.PaymentMode, &o.Payment.Status, &o.Payment.Provider, &o.Payment.Reference,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return Order{}, err
	}

	if err := r.loadItems(ctx, &o); err != nil {
		return Order{}, err
	}
	return o, nil
}

func (r *repository) loadItems(ctx context.Context, o *Order) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, product_id::text, product_name, quantity, unit_price_cents, total_cents
		FROM order_items WHERE order_id = $1 ORDER BY product_name`, o.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	o.Items = nil
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.ProductName, &item.Quantity, &item.UnitPriceCents, &item.TotalCents); err != nil {
			return err
		}
		o.Items = append(o.Items, item)
	}
	return rows.Err()
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1`, id, status)
	return err
}

type orderScanner interface {
	Scan(dest ...any) error
}

func scanOrder(scanner orderScanner, o *Order) error {
	return scanner.Scan(
		&o.ID, &o.CustomerID, &o.CustomerName, &o.CustomerPhone, &o.DeliveryAddress, &o.Status, &o.SubtotalCents, &o.DeliveryFeeCents,
		&o.TotalCents, &o.Notes, &o.Payment.PaymentMethod, &o.Payment.PaymentMode, &o.Payment.Status, &o.Payment.Provider, &o.Payment.Reference,
		&o.CreatedAt, &o.UpdatedAt,
	)
}

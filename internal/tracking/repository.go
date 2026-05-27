package tracking

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"
)

type Repository interface {
	CreateDelivery(ctx context.Context, orderID string, req CreateDeliveryRequest) (Tracking, error)
	UpdateLocation(ctx context.Context, trackingCode string, req UpdateLocationRequest) error
	UpdateDriverLocation(ctx context.Context, driverID string, trackingCode string, req UpdateLocationRequest) error
	GetByCode(ctx context.Context, trackingCode string) (Tracking, error)
	ListByDriver(ctx context.Context, driverID string) ([]Tracking, error)
	ListAvailable(ctx context.Context) ([]Tracking, error)
	Claim(ctx context.Context, driverID string, trackingCode string) error
	CountActiveByDriver(ctx context.Context, driverID string) (int, int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateDelivery(ctx context.Context, orderID string, req CreateDeliveryRequest) (Tracking, error) {
	code := generateTrackingCode()
	var trackingCode string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO deliveries (order_id, driver_id, tracking_code, estimated_arrival_at)
		VALUES ($1, NULLIF($2, '')::uuid, $3, $4)
		ON CONFLICT (order_id) DO UPDATE
		SET driver_id = EXCLUDED.driver_id, estimated_arrival_at = EXCLUDED.estimated_arrival_at, updated_at = NOW()
		RETURNING tracking_code`,
		orderID, req.DriverID, code, req.EstimatedArrivalAt).Scan(&trackingCode)
	if err != nil {
		return Tracking{}, err
	}
	return r.GetByCode(ctx, trackingCode)
}

func (r *repository) UpdateLocation(ctx context.Context, trackingCode string, req UpdateLocationRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var deliveryID string
	err = tx.QueryRowContext(ctx, `
		UPDATE deliveries
		SET status = $2, current_latitude = NULLIF($3, 0), current_longitude = NULLIF($4, 0),
		    estimated_arrival_at = COALESCE($5, estimated_arrival_at), updated_at = NOW()
		WHERE tracking_code = $1
		RETURNING id::text`, trackingCode, req.Status, req.Latitude, req.Longitude, req.EstimatedArrivalAt).Scan(&deliveryID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO delivery_events (delivery_id, status, latitude, longitude, note)
		VALUES ($1, $2, NULLIF($3, 0), NULLIF($4, 0), NULLIF($5, ''))`,
		deliveryID, req.Status, req.Latitude, req.Longitude, req.Note)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *repository) UpdateDriverLocation(ctx context.Context, driverID string, trackingCode string, req UpdateLocationRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var deliveryID string
	err = tx.QueryRowContext(ctx, `
		UPDATE deliveries
		SET status = $3, current_latitude = NULLIF($4, 0), current_longitude = NULLIF($5, 0),
		    estimated_arrival_at = COALESCE($6, estimated_arrival_at),
		    started_at = COALESCE(started_at, NOW()),
		    delivered_at = CASE WHEN $3 = 'delivered' THEN NOW() ELSE delivered_at END,
		    updated_at = NOW()
		WHERE tracking_code = $1 AND driver_id = $2::uuid
		RETURNING id::text`, strings.ToUpper(trackingCode), driverID, req.Status, req.Latitude, req.Longitude, req.EstimatedArrivalAt).Scan(&deliveryID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO delivery_events (delivery_id, status, latitude, longitude, note)
		VALUES ($1, $2, NULLIF($3, 0), NULLIF($4, 0), NULLIF($5, ''))`,
		deliveryID, req.Status, req.Latitude, req.Longitude, req.Note)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *repository) GetByCode(ctx context.Context, trackingCode string) (Tracking, error) {
	var t Tracking
	err := r.db.QueryRowContext(ctx, `
		SELECT d.order_id::text, d.tracking_code, d.status::text, c.name, c.phone,
		       concat_ws(', ', a.street || ', ' || a.number, NULLIF(a.complement, ''), a.neighborhood, a.city || ' - ' || a.state, a.zip_code),
		       COALESCE(a.latitude, 0), COALESCE(a.longitude, 0), COALESCE(dr.name, ''), COALESCE(dr.phone, ''),
		       COALESCE(d.current_latitude, 0), COALESCE(d.current_longitude, 0), d.estimated_arrival_at
		FROM deliveries d
		JOIN orders o ON o.id = d.order_id
		JOIN customers c ON c.id = o.customer_id
		JOIN customer_addresses a ON a.id = o.delivery_address_id
		LEFT JOIN drivers dr ON dr.id = d.driver_id
		WHERE d.tracking_code = $1`, strings.ToUpper(trackingCode)).
		Scan(&t.OrderID, &t.TrackingCode, &t.Status, &t.CustomerName, &t.CustomerPhone, &t.DeliveryAddress, &t.DeliveryLatitude, &t.DeliveryLongitude, &t.DriverName, &t.DriverPhone, &t.CurrentLatitude, &t.CurrentLongitude, &t.EstimatedArrivalAt)
	if err != nil {
		return Tracking{}, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT status::text, COALESCE(latitude, 0), COALESCE(longitude, 0), COALESCE(note, ''), created_at
		FROM delivery_events
		WHERE delivery_id = (SELECT id FROM deliveries WHERE tracking_code = $1)
		ORDER BY created_at`, strings.ToUpper(trackingCode))
	if err != nil {
		return Tracking{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var event TrackingEvent
		if err := rows.Scan(&event.Status, &event.Latitude, &event.Longitude, &event.Note, &event.CreatedAt); err != nil {
			return Tracking{}, err
		}
		t.Events = append(t.Events, event)
	}
	return t, rows.Err()
}

func (r *repository) ListByDriver(ctx context.Context, driverID string) ([]Tracking, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.order_id::text, d.tracking_code, d.status::text, c.name, c.phone,
		       concat_ws(', ', a.street || ', ' || a.number, NULLIF(a.complement, ''), a.neighborhood, a.city || ' - ' || a.state, a.zip_code),
		       COALESCE(a.latitude, 0), COALESCE(a.longitude, 0), COALESCE(dr.name, ''), COALESCE(dr.phone, ''),
		       COALESCE(d.current_latitude, 0), COALESCE(d.current_longitude, 0), d.estimated_arrival_at
		FROM deliveries d
		JOIN orders o ON o.id = d.order_id
		JOIN customers c ON c.id = o.customer_id
		JOIN customer_addresses a ON a.id = o.delivery_address_id
		LEFT JOIN drivers dr ON dr.id = d.driver_id
		WHERE d.driver_id = $1::uuid AND d.status <> 'delivered'
		ORDER BY d.created_at DESC`, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Tracking, 0)
	for rows.Next() {
		var t Tracking
		if err := rows.Scan(&t.OrderID, &t.TrackingCode, &t.Status, &t.CustomerName, &t.CustomerPhone, &t.DeliveryAddress, &t.DeliveryLatitude, &t.DeliveryLongitude, &t.DriverName, &t.DriverPhone, &t.CurrentLatitude, &t.CurrentLongitude, &t.EstimatedArrivalAt); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, rows.Err()
}

func (r *repository) ListAvailable(ctx context.Context) ([]Tracking, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.order_id::text, d.tracking_code, d.status::text, c.name, c.phone,
		       concat_ws(', ', a.street || ', ' || a.number, NULLIF(a.complement, ''), a.neighborhood, a.city || ' - ' || a.state, a.zip_code),
		       COALESCE(a.latitude, 0), COALESCE(a.longitude, 0), '', '',
		       COALESCE(d.current_latitude, 0), COALESCE(d.current_longitude, 0), d.estimated_arrival_at
		FROM deliveries d
		JOIN orders o ON o.id = d.order_id
		JOIN customers c ON c.id = o.customer_id
		JOIN customer_addresses a ON a.id = o.delivery_address_id
		WHERE d.driver_id IS NULL AND d.status IN ('separating', 'paid')
		ORDER BY d.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTrackingRows(rows)
}

func (r *repository) Claim(ctx context.Context, driverID string, trackingCode string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE deliveries
		SET driver_id = $2::uuid, updated_at = NOW()
		WHERE tracking_code = $1 AND driver_id IS NULL AND status <> 'delivered'`,
		strings.ToUpper(trackingCode), driverID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) CountActiveByDriver(ctx context.Context, driverID string) (int, int, error) {
	var active, max int
	err := r.db.QueryRowContext(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM deliveries WHERE driver_id = d.id AND status NOT IN ('delivered', 'canceled'))::int,
		  d.max_active_deliveries
		FROM drivers d
		WHERE d.id = $1::uuid`, driverID).Scan(&active, &max)
	return active, max, err
}

func scanTrackingRows(rows *sql.Rows) ([]Tracking, error) {
	items := make([]Tracking, 0)
	for rows.Next() {
		var t Tracking
		if err := rows.Scan(&t.OrderID, &t.TrackingCode, &t.Status, &t.CustomerName, &t.CustomerPhone, &t.DeliveryAddress, &t.DeliveryLatitude, &t.DeliveryLongitude, &t.DriverName, &t.DriverPhone, &t.CurrentLatitude, &t.CurrentLongitude, &t.EstimatedArrivalAt); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, rows.Err()
}

func generateTrackingCode() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

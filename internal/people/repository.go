package people

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

func (r *Repository) ListDrivers(ctx context.Context) ([]Person, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, name, phone, COALESCE(email, '') AS email, 'driver' AS role, is_active,
		       must_change_password, max_active_deliveries, created_at, updated_at
		FROM drivers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPeople(rows)
}

func (r *Repository) LoginDriver(ctx context.Context, email string) (Person, string, error) {
	var p Person
	var hash string
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, name, phone, COALESCE(email, '') AS email, 'driver' AS role, is_active,
		       must_change_password, max_active_deliveries, created_at, updated_at, COALESCE(password_hash, '')
		FROM drivers
		WHERE is_active = true AND LOWER(email) = LOWER($1)
		LIMIT 1`, email).
		Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt, &hash)
	return p, hash, err
}

func (r *Repository) GetDriverAuthByID(ctx context.Context, id string) (Person, string, error) {
	var p Person
	var hash string
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, name, phone, COALESCE(email, '') AS email, 'driver' AS role, is_active,
		       must_change_password, max_active_deliveries, created_at, updated_at, COALESCE(password_hash, '')
		FROM drivers
		WHERE id = $1::uuid AND is_active = true`, id).
		Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt, &hash)
	return p, hash, err
}

func (r *Repository) CreateDriver(ctx context.Context, req UpsertPersonRequest, passwordHash string) (Person, error) {
	var p Person
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO drivers (name, phone, email, password_hash, must_change_password, max_active_deliveries, is_active)
		VALUES ($1, $2, LOWER($3), $4, true, $5, $6)
		RETURNING id::text, name, phone, COALESCE(email, ''), 'driver' AS role, is_active,
		          must_change_password, max_active_deliveries, created_at, updated_at`,
		req.Name, req.Phone, req.Email, passwordHash, req.MaxActiveDeliveries, req.IsActive).
		Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) UpdateDriver(ctx context.Context, id string, req UpsertPersonRequest) (Person, error) {
	var p Person
	err := r.db.QueryRowContext(ctx, `
		UPDATE drivers
		SET name = $2, phone = $3, email = LOWER($4), max_active_deliveries = $5, is_active = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id::text, name, phone, COALESCE(email, ''), 'driver' AS role, is_active,
		          must_change_password, max_active_deliveries, created_at, updated_at`,
		id, req.Name, req.Phone, req.Email, req.MaxActiveDeliveries, req.IsActive).
		Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) ChangeDriverPassword(ctx context.Context, id string, passwordHash string) (Person, error) {
	var p Person
	err := r.db.QueryRowContext(ctx, `
		UPDATE drivers
		SET password_hash = $2, must_change_password = false, updated_at = NOW()
		WHERE id = $1
		RETURNING id::text, name, phone, COALESCE(email, ''), 'driver' AS role, is_active,
		          must_change_password, max_active_deliveries, created_at, updated_at`,
		id, passwordHash).
		Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) ListEmployees(ctx context.Context) ([]Person, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, name, phone, COALESCE(email, ''), role::text, is_active, false, 0, created_at, updated_at
		FROM employees ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPeople(rows)
}

func (r *Repository) CreateEmployee(ctx context.Context, req UpsertPersonRequest) (Person, error) {
	var p Person
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO employees (name, phone, email, role, is_active)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5)
		RETURNING id::text, name, phone, COALESCE(email, ''), role::text, is_active, false, 0, created_at, updated_at`,
		req.Name, req.Phone, req.Email, req.Role, req.IsActive).Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) UpdateEmployee(ctx context.Context, id string, req UpsertPersonRequest) (Person, error) {
	var p Person
	err := r.db.QueryRowContext(ctx, `
		UPDATE employees SET name = $2, phone = $3, email = NULLIF($4, ''), role = $5, is_active = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id::text, name, phone, COALESCE(email, ''), role::text, is_active, false, 0, created_at, updated_at`,
		id, req.Name, req.Phone, req.Email, req.Role, req.IsActive).Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func scanPeople(rows *sql.Rows) ([]Person, error) {
	items := make([]Person, 0)
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Email, &p.Role, &p.IsActive, &p.MustChangePassword, &p.MaxActiveDeliveries, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

package auth

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, req RegisterRequest, passwordHash string) (User, error)
	FindByEmail(ctx context.Context, email string) (UserWithPassword, error)
	FindByID(ctx context.Context, id string) (User, error)
}

type UserWithPassword struct {
	User
	PasswordHash string
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, req RegisterRequest, passwordHash string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO admin_users (name, email, phone, store_name, password_hash, role)
		VALUES ($1, LOWER($2), NULLIF($3, ''), NULLIF($4, ''), $5, 'owner')
		RETURNING id::text, name, email, role, COALESCE(store_name, '')`,
		req.Name, req.Email, req.Phone, req.StoreName, passwordHash,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.StoreName)
	if err != nil {
		return User{}, err
	}
	user.Initials = initials(user.Name)
	user.Color = "#6E1F2C"
	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (UserWithPassword, error) {
	var user UserWithPassword
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, name, email, role, COALESCE(store_name, ''), password_hash
		FROM admin_users
		WHERE LOWER(email) = LOWER($1) AND is_active = true`, email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.StoreName, &user.PasswordHash)
	if err != nil {
		return UserWithPassword{}, err
	}
	user.Initials = initials(user.Name)
	user.Color = "#6E1F2C"
	return user, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, name, email, role, COALESCE(store_name, '')
		FROM admin_users
		WHERE id = $1 AND is_active = true`, id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.StoreName)
	if err != nil {
		return User{}, err
	}
	user.Initials = initials(user.Name)
	user.Color = "#6E1F2C"
	return user, nil
}

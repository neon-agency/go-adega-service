package people

import "time"

type Person struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Phone               string    `json:"phone"`
	Email               string    `json:"email,omitempty"`
	Role                string    `json:"role,omitempty"`
	IsActive            bool      `json:"is_active"`
	MustChangePassword  bool      `json:"must_change_password,omitempty"`
	MaxActiveDeliveries int       `json:"max_active_deliveries,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}

type UpsertPersonRequest struct {
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Role                string `json:"role"`
	IsActive            bool   `json:"is_active"`
	MaxActiveDeliveries int    `json:"max_active_deliveries"`
}

type DriverLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DriverChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

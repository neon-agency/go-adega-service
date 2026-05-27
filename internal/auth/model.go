package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	StoreName string `json:"storeName"`
	Password  string `json:"password"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	StoreName string `json:"storeName"`
	Initials  string `json:"initials"`
	Color     string `json:"color"`
}

type Session struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

package orders

import "time"

type CustomerRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Document string `json:"document"`
}

type AddressRequest struct {
	Street       string  `json:"street"`
	Number       string  `json:"number"`
	Complement   string  `json:"complement"`
	Neighborhood string  `json:"neighborhood"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	ZipCode      string  `json:"zip_code"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CreateOrderRequest struct {
	Customer      CustomerRequest    `json:"customer"`
	Address       AddressRequest     `json:"address"`
	PaymentMethod string             `json:"payment_method"`
	PaymentMode   string             `json:"payment_mode"`
	PaymentToken  string             `json:"payment_token"`
	Provider      string             `json:"provider"`
	Installments  int                `json:"installments"`
	Notes         string             `json:"notes"`
	Items         []OrderItemRequest `json:"items"`
	DeliveryFee   int                `json:"-"`
}

type OrderItem struct {
	ID             string `json:"id"`
	ProductID      string `json:"product_id"`
	ProductName    string `json:"product_name"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int    `json:"unit_price_cents"`
	TotalCents     int    `json:"total_cents"`
}

type PaymentInfo struct {
	Provider      string `json:"provider,omitempty"`
	Reference     string `json:"reference,omitempty"`
	Status        string `json:"status"`
	QRCode        string `json:"qr_code,omitempty"`
	CopyPaste     string `json:"copy_paste,omitempty"`
	PaymentURL    string `json:"payment_url,omitempty"`
	PaymentMethod string `json:"payment_method"`
	PaymentMode   string `json:"payment_mode"`
}

type Order struct {
	ID               string      `json:"id"`
	CustomerID       string      `json:"customer_id"`
	CustomerName     string      `json:"customer_name"`
	CustomerPhone    string      `json:"customer_phone"`
	DeliveryAddress  string      `json:"delivery_address,omitempty"`
	Status           string      `json:"status"`
	SubtotalCents    int         `json:"subtotal_cents"`
	DeliveryFeeCents int         `json:"delivery_fee_cents"`
	TotalCents       int         `json:"total_cents"`
	Notes            string      `json:"notes,omitempty"`
	Payment          PaymentInfo `json:"payment"`
	Items            []OrderItem `json:"items"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

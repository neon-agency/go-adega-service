package tracking

import "time"

type Tracking struct {
	OrderID            string          `json:"order_id"`
	TrackingCode       string          `json:"tracking_code"`
	Status             string          `json:"status"`
	CustomerName       string          `json:"customer_name,omitempty"`
	CustomerPhone      string          `json:"customer_phone,omitempty"`
	DeliveryAddress    string          `json:"delivery_address,omitempty"`
	DeliveryLatitude   float64         `json:"delivery_latitude,omitempty"`
	DeliveryLongitude  float64         `json:"delivery_longitude,omitempty"`
	DriverName         string          `json:"driver_name,omitempty"`
	DriverPhone        string          `json:"driver_phone,omitempty"`
	CurrentLatitude    float64         `json:"current_latitude,omitempty"`
	CurrentLongitude   float64         `json:"current_longitude,omitempty"`
	EstimatedArrivalAt *time.Time      `json:"estimated_arrival_at,omitempty"`
	Events             []TrackingEvent `json:"events"`
}

type TrackingEvent struct {
	Status    string    `json:"status"`
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDeliveryRequest struct {
	DriverID           string     `json:"driver_id"`
	EstimatedArrivalAt *time.Time `json:"estimated_arrival_at"`
}

type UpdateLocationRequest struct {
	Status             string     `json:"status"`
	Latitude           float64    `json:"latitude"`
	Longitude          float64    `json:"longitude"`
	EstimatedArrivalAt *time.Time `json:"estimated_arrival_at"`
	Note               string     `json:"note"`
}

type ClaimDeliveryRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

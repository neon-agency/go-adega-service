package settings

import "time"

type StoreSettings struct {
	ID                       string        `json:"id"`
	StoreName                string        `json:"store_name"`
	CNPJ                     string        `json:"cnpj,omitempty"`
	ShortDescription         string        `json:"short_description,omitempty"`
	BrandColor               string        `json:"brand_color,omitempty"`
	Instagram                string        `json:"instagram,omitempty"`
	Phone                    string        `json:"phone"`
	Whatsapp                 string        `json:"whatsapp,omitempty"`
	LogoURL                  string        `json:"logo_url,omitempty"`
	BannerURL                string        `json:"banner_url,omitempty"`
	AddressStreet            string        `json:"address_street"`
	AddressNumber            string        `json:"address_number"`
	AddressComplement        string        `json:"address_complement,omitempty"`
	AddressNeighborhood      string        `json:"address_neighborhood"`
	AddressCity              string        `json:"address_city"`
	AddressState             string        `json:"address_state"`
	AddressZipCode           string        `json:"address_zip_code"`
	DeliveryFeeCents         int           `json:"delivery_fee_cents"`
	FreeDeliveryFromCents    int           `json:"free_delivery_from_cents"`
	MinOrderCents            int           `json:"min_order_cents"`
	AllowDelivery            bool          `json:"allow_delivery"`
	AllowPickup              bool          `json:"allow_pickup"`
	DeliveryRadiusKM         float64       `json:"delivery_radius_km"`
	AverageDeliveryMin       int           `json:"average_delivery_min_minutes"`
	AverageDeliveryMax       int           `json:"average_delivery_max_minutes"`
	Latitude                 float64       `json:"latitude,omitempty"`
	Longitude                float64       `json:"longitude,omitempty"`
	DriverPickupRadiusMeters int           `json:"driver_pickup_radius_meters"`
	OpeningTime              string        `json:"opening_time"`
	ClosingTime              string        `json:"closing_time"`
	OpeningHours             []OpeningHour `json:"opening_hours"`
	Notifications            Notifications `json:"notifications"`
	IsOpen                   bool          `json:"is_open"`
	AcceptOnlinePix          bool          `json:"accept_online_pix"`
	AcceptOnlineCard         bool          `json:"accept_online_card"`
	AcceptDeliveryPix        bool          `json:"accept_delivery_pix"`
	AcceptDeliveryCard       bool          `json:"accept_delivery_card"`
	AcceptDeliveryCash       bool          `json:"accept_delivery_cash"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}

type OpeningHour struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Open  bool   `json:"open"`
	From  string `json:"from"`
	To    string `json:"to"`
}

type Notifications struct {
	SoundNew     bool `json:"sound_new"`
	SoundConfirm bool `json:"sound_confirm"`
	Browser      bool `json:"browser"`
	Daily        bool `json:"daily"`
	Weekly       bool `json:"weekly"`
	LowStock     bool `json:"low_stock"`
	CustomerMsg  bool `json:"customer_msg"`
}

type UpdateStoreSettingsRequest struct {
	StoreName                string        `json:"store_name"`
	CNPJ                     string        `json:"cnpj"`
	ShortDescription         string        `json:"short_description"`
	BrandColor               string        `json:"brand_color"`
	Instagram                string        `json:"instagram"`
	Phone                    string        `json:"phone"`
	Whatsapp                 string        `json:"whatsapp"`
	LogoURL                  string        `json:"logo_url"`
	BannerURL                string        `json:"banner_url"`
	AddressStreet            string        `json:"address_street"`
	AddressNumber            string        `json:"address_number"`
	AddressComplement        string        `json:"address_complement"`
	AddressNeighborhood      string        `json:"address_neighborhood"`
	AddressCity              string        `json:"address_city"`
	AddressState             string        `json:"address_state"`
	AddressZipCode           string        `json:"address_zip_code"`
	DeliveryFeeCents         int           `json:"delivery_fee_cents"`
	FreeDeliveryFromCents    int           `json:"free_delivery_from_cents"`
	MinOrderCents            int           `json:"min_order_cents"`
	AllowDelivery            bool          `json:"allow_delivery"`
	AllowPickup              bool          `json:"allow_pickup"`
	DeliveryRadiusKM         float64       `json:"delivery_radius_km"`
	AverageDeliveryMin       int           `json:"average_delivery_min_minutes"`
	AverageDeliveryMax       int           `json:"average_delivery_max_minutes"`
	Latitude                 float64       `json:"latitude"`
	Longitude                float64       `json:"longitude"`
	DriverPickupRadiusMeters int           `json:"driver_pickup_radius_meters"`
	OpeningTime              string        `json:"opening_time"`
	ClosingTime              string        `json:"closing_time"`
	OpeningHours             []OpeningHour `json:"opening_hours"`
	Notifications            Notifications `json:"notifications"`
	IsOpen                   bool          `json:"is_open"`
	AcceptOnlinePix          bool          `json:"accept_online_pix"`
	AcceptOnlineCard         bool          `json:"accept_online_card"`
	AcceptDeliveryPix        bool          `json:"accept_delivery_pix"`
	AcceptDeliveryCard       bool          `json:"accept_delivery_card"`
	AcceptDeliveryCash       bool          `json:"accept_delivery_cash"`
}

package settings

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Repository interface {
	Get(ctx context.Context) (StoreSettings, error)
	Update(ctx context.Context, req UpdateStoreSettingsRequest) (StoreSettings, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Get(ctx context.Context) (StoreSettings, error) {
	var s StoreSettings
	var openingHours []byte
	var notifications []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, store_name, COALESCE(cnpj, ''), COALESCE(short_description, ''),
		       COALESCE(brand_color, '#6E1F2C'), COALESCE(instagram, ''),
		       phone, COALESCE(whatsapp, ''), COALESCE(logo_url, ''), COALESCE(banner_url, ''), address_street, address_number,
		       COALESCE(address_complement, ''), address_neighborhood, address_city, address_state,
		       address_zip_code, delivery_fee_cents, free_delivery_from_cents, min_order_cents,
		       allow_delivery, allow_pickup, delivery_radius_km, average_delivery_min_minutes, average_delivery_max_minutes,
		       COALESCE(latitude, 0), COALESCE(longitude, 0), driver_pickup_radius_meters,
		       opening_time::text, closing_time::text, opening_hours, notifications,
		       is_open, accept_online_pix, accept_online_card, accept_delivery_pix, accept_delivery_card,
		       accept_delivery_cash, created_at, updated_at
		FROM store_settings
		ORDER BY created_at
		LIMIT 1`).Scan(
		&s.ID, &s.StoreName, &s.CNPJ, &s.ShortDescription, &s.BrandColor, &s.Instagram,
		&s.Phone, &s.Whatsapp, &s.LogoURL, &s.BannerURL, &s.AddressStreet, &s.AddressNumber,
		&s.AddressComplement, &s.AddressNeighborhood, &s.AddressCity, &s.AddressState,
		&s.AddressZipCode, &s.DeliveryFeeCents, &s.FreeDeliveryFromCents, &s.MinOrderCents,
		&s.AllowDelivery, &s.AllowPickup, &s.DeliveryRadiusKM, &s.AverageDeliveryMin, &s.AverageDeliveryMax,
		&s.Latitude, &s.Longitude, &s.DriverPickupRadiusMeters,
		&s.OpeningTime, &s.ClosingTime, &openingHours, &notifications,
		&s.IsOpen, &s.AcceptOnlinePix, &s.AcceptOnlineCard, &s.AcceptDeliveryPix, &s.AcceptDeliveryCard,
		&s.AcceptDeliveryCash, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return s, err
	}
	if len(openingHours) > 0 {
		_ = json.Unmarshal(openingHours, &s.OpeningHours)
	}
	if len(notifications) > 0 {
		_ = json.Unmarshal(notifications, &s.Notifications)
	}
	return s, err
}

func (r *repository) Update(ctx context.Context, req UpdateStoreSettingsRequest) (StoreSettings, error) {
	current, err := r.Get(ctx)
	if err != nil {
		return StoreSettings{}, err
	}

	var s StoreSettings
	openingHours, err := json.Marshal(req.OpeningHours)
	if err != nil {
		return StoreSettings{}, err
	}
	notifications, err := json.Marshal(req.Notifications)
	if err != nil {
		return StoreSettings{}, err
	}
	var savedOpeningHours []byte
	var savedNotifications []byte
	err = r.db.QueryRowContext(ctx, `
		UPDATE store_settings
		SET store_name = $2,
		    cnpj = NULLIF($3, ''),
		    short_description = NULLIF($4, ''),
		    brand_color = COALESCE(NULLIF($5, ''), '#6E1F2C'),
		    instagram = NULLIF($6, ''),
		    phone = $7,
		    whatsapp = NULLIF($8, ''),
		    logo_url = NULLIF($9, ''),
		    banner_url = NULLIF($10, ''),
		    address_street = $11,
		    address_number = $12,
		    address_complement = NULLIF($13, ''),
		    address_neighborhood = $14,
		    address_city = $15,
		    address_state = $16,
		    address_zip_code = $17,
		    delivery_fee_cents = $18,
		    free_delivery_from_cents = $19,
		    min_order_cents = $20,
		    allow_delivery = $21,
		    allow_pickup = $22,
		    delivery_radius_km = $23,
		    average_delivery_min_minutes = $24,
		    average_delivery_max_minutes = $25,
		    latitude = NULLIF($26, 0),
		    longitude = NULLIF($27, 0),
		    driver_pickup_radius_meters = $28,
		    opening_time = $29,
		    closing_time = $30,
		    opening_hours = $31::jsonb,
		    notifications = $32::jsonb,
		    is_open = $33,
		    accept_online_pix = $34,
		    accept_online_card = $35,
		    accept_delivery_pix = $36,
		    accept_delivery_card = $37,
		    accept_delivery_cash = $38,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id::text, store_name, COALESCE(cnpj, ''), COALESCE(short_description, ''),
		          COALESCE(brand_color, '#6E1F2C'), COALESCE(instagram, ''),
		          phone, COALESCE(whatsapp, ''), COALESCE(logo_url, ''), COALESCE(banner_url, ''), address_street, address_number,
		          COALESCE(address_complement, ''), address_neighborhood, address_city, address_state,
		          address_zip_code, delivery_fee_cents, free_delivery_from_cents, min_order_cents,
		          allow_delivery, allow_pickup, delivery_radius_km, average_delivery_min_minutes, average_delivery_max_minutes,
		          COALESCE(latitude, 0), COALESCE(longitude, 0), driver_pickup_radius_meters,
		          opening_time::text, closing_time::text, opening_hours, notifications,
		          is_open, accept_online_pix, accept_online_card, accept_delivery_pix, accept_delivery_card,
		          accept_delivery_cash, created_at, updated_at`,
		current.ID, req.StoreName, req.CNPJ, req.ShortDescription, req.BrandColor, req.Instagram,
		req.Phone, req.Whatsapp, req.LogoURL, req.BannerURL, req.AddressStreet, req.AddressNumber,
		req.AddressComplement, req.AddressNeighborhood, req.AddressCity, req.AddressState,
		req.AddressZipCode, req.DeliveryFeeCents, req.FreeDeliveryFromCents, req.MinOrderCents,
		req.AllowDelivery, req.AllowPickup, req.DeliveryRadiusKM, req.AverageDeliveryMin, req.AverageDeliveryMax,
		req.Latitude, req.Longitude, req.DriverPickupRadiusMeters,
		req.OpeningTime, req.ClosingTime, string(openingHours), string(notifications),
		req.IsOpen, req.AcceptOnlinePix, req.AcceptOnlineCard, req.AcceptDeliveryPix, req.AcceptDeliveryCard, req.AcceptDeliveryCash,
	).Scan(
		&s.ID, &s.StoreName, &s.CNPJ, &s.ShortDescription, &s.BrandColor, &s.Instagram,
		&s.Phone, &s.Whatsapp, &s.LogoURL, &s.BannerURL, &s.AddressStreet, &s.AddressNumber,
		&s.AddressComplement, &s.AddressNeighborhood, &s.AddressCity, &s.AddressState,
		&s.AddressZipCode, &s.DeliveryFeeCents, &s.FreeDeliveryFromCents, &s.MinOrderCents,
		&s.AllowDelivery, &s.AllowPickup, &s.DeliveryRadiusKM, &s.AverageDeliveryMin, &s.AverageDeliveryMax,
		&s.Latitude, &s.Longitude, &s.DriverPickupRadiusMeters,
		&s.OpeningTime, &s.ClosingTime, &savedOpeningHours, &savedNotifications,
		&s.IsOpen, &s.AcceptOnlinePix, &s.AcceptOnlineCard, &s.AcceptDeliveryPix, &s.AcceptDeliveryCard,
		&s.AcceptDeliveryCash, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return s, err
	}
	if len(savedOpeningHours) > 0 {
		_ = json.Unmarshal(savedOpeningHours, &s.OpeningHours)
	}
	if len(savedNotifications) > 0 {
		_ = json.Unmarshal(savedNotifications, &s.Notifications)
	}
	return s, err
}

package orders

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/GabsMeloTI/go_adega/internal/paymentclient"
	"github.com/GabsMeloTI/go_adega/internal/settings"
)

type Service struct {
	repo          Repository
	settingsRepo  settings.Repository
	paymentClient *paymentclient.Client
	provider      string
}

func NewService(repo Repository, settingsRepo settings.Repository, paymentClient *paymentclient.Client, provider string) *Service {
	if strings.TrimSpace(provider) == "" {
		provider = "efi"
	}
	return &Service{repo: repo, settingsRepo: settingsRepo, paymentClient: paymentClient, provider: provider}
}

func (s *Service) Create(ctx context.Context, req CreateOrderRequest) (Order, error) {
	if strings.TrimSpace(req.Customer.Name) == "" || strings.TrimSpace(req.Customer.Phone) == "" {
		return Order{}, errors.New("cliente inválido")
	}
	if len(req.Items) == 0 {
		return Order{}, errors.New("pedido sem itens")
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = "pix"
	}
	if req.PaymentMode == "" {
		req.PaymentMode = "online"
	}
	if req.Provider == "" {
		req.Provider = s.provider
	}

	storeSettings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return Order{}, err
	}
	if !storeSettings.IsOpen {
		return Order{}, errors.New("adega fechada no momento")
	}
	if err := validatePayment(storeSettings, req.PaymentMode, req.PaymentMethod); err != nil {
		return Order{}, err
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()

	customerID, err := s.repo.CreateCustomer(ctx, tx, req.Customer)
	if err != nil {
		return Order{}, err
	}
	addressID, err := s.repo.CreateAddress(ctx, tx, customerID, req.Address)
	if err != nil {
		return Order{}, err
	}

	items := make([]OrderItem, 0, len(req.Items))
	subtotal := 0
	for _, requested := range req.Items {
		if requested.Quantity <= 0 {
			return Order{}, errors.New("quantidade do item deve ser maior que zero")
		}
		item, err := s.repo.ReserveProduct(ctx, tx, requested.ProductID, requested.Quantity)
		if err != nil {
			return Order{}, err
		}
		items = append(items, item)
		subtotal += item.TotalCents
	}

	if storeSettings.MinOrderCents > 0 && subtotal < storeSettings.MinOrderCents {
		return Order{}, errors.New("pedido abaixo do valor mínimo")
	}

	req.DeliveryFee = storeSettings.DeliveryFeeCents
	if storeSettings.FreeDeliveryFromCents > 0 && subtotal >= storeSettings.FreeDeliveryFromCents {
		req.DeliveryFee = 0
	}

	payment := PaymentInfo{PaymentMethod: req.PaymentMethod, PaymentMode: req.PaymentMode, Status: "pending"}
	if req.PaymentMode == "online" && (req.PaymentMethod == "pix" || req.PaymentMethod == "credit_card" || req.PaymentMethod == "debit_card") {
		total := subtotal + req.DeliveryFee
		paymentReq := buildPaymentRequest(req, items, total)
		log.Printf("[orders] processing online payment method=%s provider=%s subtotal=%d delivery_fee=%d total=%d customer=%q items=%d",
			req.PaymentMethod, paymentReq.Provider, subtotal, req.DeliveryFee, total, req.Customer.Name, len(items))
		resp, err := s.paymentClient.Process(ctx, paymentReq)
		if err != nil {
			log.Printf("[orders] payment failed method=%s provider=%s total=%d error=%v", req.PaymentMethod, paymentReq.Provider, total, err)
			return Order{}, err
		}
		log.Printf("[orders] payment processed provider=%s provider_id=%s status=%s amount=%d", resp.Provider, resp.ProviderID, resp.Status, resp.Amount)
		payment.Provider = resp.Provider
		payment.Reference = resp.ProviderID
		payment.Status = resp.Status
		payment.QRCode = resp.QRCode
		payment.CopyPaste = resp.CopyPaste
		payment.PaymentURL = resp.PaymentURL
	}

	orderID, err := s.repo.CreateOrder(ctx, tx, customerID, addressID, req, subtotal, payment)
	if err != nil {
		return Order{}, err
	}
	for _, item := range items {
		if err := s.repo.AddOrderItem(ctx, tx, orderID, item); err != nil {
			return Order{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return Order{}, err
	}

	order, err := s.repo.Get(ctx, orderID)
	if err != nil {
		return Order{}, err
	}
	order.Payment.QRCode = payment.QRCode
	order.Payment.CopyPaste = payment.CopyPaste
	order.Payment.PaymentURL = payment.PaymentURL
	return order, nil
}

func validatePayment(storeSettings settings.StoreSettings, mode string, method string) error {
	switch mode {
	case "online":
		switch method {
		case "pix":
			if !storeSettings.AcceptOnlinePix {
				return errors.New("pix online não está habilitado")
			}
		case "credit_card", "debit_card":
			if !storeSettings.AcceptOnlineCard {
				return errors.New("cartão online não está habilitado")
			}
		case "cash":
			return errors.New("dinheiro não é permitido no pagamento online")
		default:
			return errors.New("forma de pagamento inválida")
		}
	case "delivery":
		switch method {
		case "pix":
			if !storeSettings.AcceptDeliveryPix {
				return errors.New("pix na entrega não está habilitado")
			}
		case "credit_card", "debit_card":
			if !storeSettings.AcceptDeliveryCard {
				return errors.New("cartão na entrega não está habilitado")
			}
		case "cash":
			if !storeSettings.AcceptDeliveryCash {
				return errors.New("dinheiro na entrega não está habilitado")
			}
		default:
			return errors.New("forma de pagamento inválida")
		}
	default:
		return errors.New("modo de pagamento inválido")
	}
	return nil
}

func (s *Service) List(ctx context.Context, status string) ([]Order, error) {
	return s.repo.List(ctx, status)
}

func (s *Service) Get(ctx context.Context, id string) (Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status string) (Order, error) {
	switch status {
	case "created", "awaiting_payment", "paid", "separating", "out_for_delivery", "delivered", "canceled":
		if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
			return Order{}, err
		}
		return s.repo.Get(ctx, id)
	default:
		return Order{}, errors.New("status inválido")
	}
}

func buildPaymentRequest(req CreateOrderRequest, items []OrderItem, total int) paymentclient.PaymentRequest {
	method := "pix"
	if req.PaymentMethod == "credit_card" || req.PaymentMethod == "debit_card" {
		method = "card"
	}

	paymentItems := make([]paymentclient.Item, 0, len(items))
	for _, item := range items {
		paymentItems = append(paymentItems, paymentclient.Item{
			Name:   item.ProductName,
			Value:  item.UnitPriceCents,
			Amount: item.Quantity,
		})
	}

	customer := paymentclient.Customer{
		Name:  req.Customer.Name,
		Email: req.Customer.Email,
		Phone: req.Customer.Phone,
	}
	if cpf := digitsOnly(req.Customer.Document); len(cpf) == 11 {
		customer.CPF = cpf
	}
	if address := buildPaymentAddress(req.Address); address != nil {
		customer.Address = address
	}

	metadata := map[string]string{
		"source": "go_adega",
	}
	if req.PaymentMethod == "credit_card" || req.PaymentMethod == "debit_card" {
		metadata["pagarme_payment_method"] = req.PaymentMethod
	}

	return paymentclient.PaymentRequest{
		Provider:     req.Provider,
		Method:       method,
		Amount:       total,
		Currency:     "BRL",
		Description:  "Pedido Adega",
		PaymentToken: req.PaymentToken,
		Installments: req.Installments,
		Customer:     customer,
		Items:        paymentItems,
		Metadata:     metadata,
	}
}

func buildPaymentAddress(address AddressRequest) *paymentclient.Address {
	if strings.TrimSpace(address.Street) == "" &&
		strings.TrimSpace(address.Number) == "" &&
		strings.TrimSpace(address.Neighborhood) == "" &&
		strings.TrimSpace(address.City) == "" &&
		strings.TrimSpace(address.State) == "" &&
		strings.TrimSpace(address.ZipCode) == "" {
		return nil
	}

	return &paymentclient.Address{
		Street:       address.Street,
		Number:       address.Number,
		Neighborhood: address.Neighborhood,
		City:         address.City,
		State:        address.State,
		ZipCode:      address.ZipCode,
	}
}

func digitsOnly(value string) string {
	var b strings.Builder
	for _, char := range value {
		if char >= '0' && char <= '9' {
			b.WriteRune(char)
		}
	}
	return b.String()
}

package orders

import "testing"

func TestBuildPaymentRequestForPix(t *testing.T) {
	req := CreateOrderRequest{
		Customer: CustomerRequest{
			Name:     "Gabriel",
			Document: "123.456.789-01",
		},
		PaymentMethod: "pix",
		Provider:      "efi",
	}

	paymentReq := buildPaymentRequest(req, nil, 10789)

	if paymentReq.Provider != "efi" {
		t.Fatalf("expected provider efi, got %s", paymentReq.Provider)
	}
	if paymentReq.Method != "pix" {
		t.Fatalf("expected method pix, got %s", paymentReq.Method)
	}
	if paymentReq.Amount != 10789 {
		t.Fatalf("expected amount in cents, got %d", paymentReq.Amount)
	}
	if paymentReq.Customer.Name != "Gabriel" {
		t.Fatalf("expected customer name Gabriel, got %s", paymentReq.Customer.Name)
	}
	if paymentReq.Customer.CPF != "12345678901" {
		t.Fatalf("expected only CPF digits, got %s", paymentReq.Customer.CPF)
	}
}

func TestBuildPaymentRequestOmitsInvalidCPF(t *testing.T) {
	req := CreateOrderRequest{
		Customer: CustomerRequest{
			Name:     "Gabriel",
			Document: "123",
		},
		PaymentMethod: "pix",
		Provider:      "efi",
	}

	paymentReq := buildPaymentRequest(req, nil, 10789)

	if paymentReq.Customer.CPF != "" {
		t.Fatalf("expected invalid CPF to be omitted, got %s", paymentReq.Customer.CPF)
	}
}

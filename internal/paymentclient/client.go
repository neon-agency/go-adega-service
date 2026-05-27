package paymentclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Address struct {
	Street       string `json:"street,omitempty"`
	Number       string `json:"number,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	ZipCode      string `json:"zipcode,omitempty"`
}

type Customer struct {
	Name    string   `json:"name"`
	Email   string   `json:"email,omitempty"`
	CPF     string   `json:"cpf,omitempty"`
	Phone   string   `json:"phone,omitempty"`
	Address *Address `json:"address,omitempty"`
}

type Item struct {
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Amount      int    `json:"amount"`
	Description string `json:"description,omitempty"`
}

type PaymentRequest struct {
	Provider     string            `json:"provider"`
	Method       string            `json:"method"`
	Amount       int               `json:"amount"`
	Currency     string            `json:"currency,omitempty"`
	Description  string            `json:"description,omitempty"`
	Customer     Customer          `json:"customer"`
	Items        []Item            `json:"items"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	PaymentToken string            `json:"payment_token,omitempty"`
	Installments int               `json:"installments,omitempty"`
	PixKey       string            `json:"pix_key,omitempty"`
}

type PaymentResponse struct {
	TransactionID string            `json:"transaction_id"`
	ProviderID    string            `json:"provider_id"`
	Status        string            `json:"status"`
	Amount        int               `json:"amount"`
	Method        string            `json:"method"`
	Provider      string            `json:"provider"`
	QRCode        string            `json:"qr_code,omitempty"`
	CopyPaste     string            `json:"copy_paste,omitempty"`
	PaymentURL    string            `json:"payment_url,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

func New(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: httpClient}
}

func (c *Client) Process(ctx context.Context, req PaymentRequest) (PaymentResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return PaymentResponse{}, err
	}

	url := c.baseURL + "/payments"
	log.Printf("[paymentclient] POST %s payload=%s", url, string(maskPaymentPayload(body)))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return PaymentResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		log.Printf("[paymentclient] request failed error=%v", err)
		return PaymentResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[paymentclient] cannot read response status=%d error=%v", resp.StatusCode, err)
		return PaymentResponse{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[paymentclient] non-2xx response status=%d body=%s", resp.StatusCode, string(respBody))
		return PaymentResponse{}, fmt.Errorf("payment service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var out PaymentResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		log.Printf("[paymentclient] decode response failed status=%d body=%s error=%v", resp.StatusCode, string(respBody), err)
		return PaymentResponse{}, err
	}
	log.Printf("[paymentclient] response status=%d provider=%s method=%s provider_id=%s payment_status=%s amount=%d",
		resp.StatusCode, out.Provider, out.Method, out.ProviderID, out.Status, out.Amount)
	return out, nil
}

func maskPaymentPayload(body []byte) []byte {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return body
	}
	if token, ok := payload["payment_token"].(string); ok && token != "" {
		payload["payment_token"] = "***"
	}
	if customer, ok := payload["customer"].(map[string]any); ok {
		if cpf, ok := customer["cpf"].(string); ok && len(cpf) > 3 {
			customer["cpf"] = strings.Repeat("*", len(cpf)-3) + cpf[len(cpf)-3:]
		}
	}
	masked, err := json.Marshal(payload)
	if err != nil {
		return body
	}
	return masked
}

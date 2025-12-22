package yoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const PaymentCurrency = "RUB"


type Client struct {
	shopID string
	secretKey string
}

func NewClient(shopID, secretKey string) (*Client, error) {
	if shopID == "" || secretKey == "" {
		return &Client{}, fmt.Errorf("missing YOOKASSA credentials")
	}
	return &Client{shopID: shopID, secretKey: secretKey}, nil
}


type YooAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}
type YooItem struct {
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Amount      YooAmount `json:"amount"`
	VatCode        int    `json:"vat_code"`
	PaymentMode    string `json:"payment_mode"`
	PaymentSubject string `json:"payment_subject"`
}
type YooCustomer struct {
	Email string `json:"email"`
}
type Receipt struct {
	Customer YooCustomer `json:"customer"`
	Items         []YooItem `json:"items"`
	TaxSystemCode int    `json:"tax_system_code"`
}
type YooConfirmation struct {
	Type string `json:"type"`
}
type PaymentRequest struct {
	Amount YooAmount `json:"amount"`
	Confirmation YooConfirmation `json:"confirmation"`
	Capture     bool    `json:"capture"`
	Description string  `json:"description"`
	Receipt     Receipt `json:"receipt"`
}
type PaymentConfirmation struct {
	ConfirmationToken string `json:"confirmation_token"`
}
type PaymentResponse struct {
	Confirmation PaymentConfirmation `json:"confirmation"`
	Id           string              `json:"id"`
}

func (c Client) CreatePayment(amount int, email, description, idempotenceKey string) (string, string, error) {
	amountValue := fmt.Sprintf("%.2f", float64(amount)/100.0)
	reqBody := PaymentRequest{
		Amount: YooAmount{Value: amountValue, Currency: PaymentCurrency},
		Confirmation: YooConfirmation{Type: "embedded"},
		Capture:     true,
		Description: description,
		Receipt: Receipt{
			Customer: YooCustomer{Email: email},
			Items: []YooItem{
				{
					Description: description,
					Quantity:    1,
					Amount: YooAmount{Value: amountValue, Currency: PaymentCurrency},
					VatCode:        1,
					PaymentMode:    "full_payment",
					PaymentSubject: "commodity",
				},
			},
			TaxSystemCode: 2,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.SetBasicAuth(c.shopID, c.secretKey)
	req.Header.Set("Idempotence-Key", idempotenceKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBytes, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("payment creation failed: %s; %s", resp.Status, string(responseBytes))
	}

	var paymentResp PaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&paymentResp)
	if err != nil {
		return "", "", err
	}

	return paymentResp.Id, paymentResp.Confirmation.ConfirmationToken, nil
}

func (c Client) GetPayment(id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	url := fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}

	req.SetBasicAuth(c.shopID, c.secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("failed to get payment status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (c Client) GetPaymentStatus(id string) (string, error) {
	data, err := c.GetPayment(id)
	if err != nil {
		return "", err
	}

	status, ok := data["status"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response")
	}
	return status, nil
}

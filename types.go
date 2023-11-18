package main

import (
	"encoding/json"
	"time"
)

type PaymentLinkRequest struct {
	Order struct {
		Amount        int    `json:"amount"`
		InvoiceNumber string `json:"invoice_number"`
	} `json:"order"`
	Payment struct {
		PaymentDueDate int `json:"payment_due_date"`
	} `json:"payment"`
}

func (p *PaymentLinkRequest) ToJSON() ([]byte, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

type PaymentLinkResponse struct {
	Message  []string `json:"message"`
	Response struct {
		Order struct {
			Amount        string `json:"amount"`
			InvoiceNumber string `json:"invoice_number"`
			Currency      string `json:"currency"`
			SessionID     string `json:"session_id"`
		} `json:"order"`
		Payment struct {
			PaymentMethodTypes []string `json:"payment_method_types"`
			PaymentDueDate     int      `json:"payment_due_date"`
			TokenID            string   `json:"token_id"`
			URL                string   `json:"url"`
			ExpiredDate        string   `json:"expired_date"`
		} `json:"payment"`
		AdditionalInfo struct {
			DokuCheckout bool `json:"doku_checkout"`
		} `json:"additional_info"`
		Headers struct {
			RequestID string    `json:"request_id"`
			Signature string    `json:"signature"`
			Date      time.Time `json:"date"`
			ClientID  string    `json:"client_id"`
		} `json:"headers"`
	} `json:"response"`
}

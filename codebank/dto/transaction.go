package dto

import "time"

type Transaction struct {
	ID              string    `json:"transaction_id"`
	Name            string    `json:"-"`
	Number          string    `json:"credit_card_number"`
	ExpirationMonth int32     `json:"-"`
	ExpirationYear  int32     `json:"-"`
	CVV             int32     `json:"amount"`
	Amount          float64   `json:"store"`
	Store           string    `json:"description"`
	Description     string    `json:"-"`
	CreatedAt       time.Time `json:"payment_date"`
}

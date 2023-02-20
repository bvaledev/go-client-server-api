package entity

import "github.com/google/uuid"

type Quotation struct {
	ID  string  `json:"id"`
	Bid float64 `json:"bid"`
}

func NewQuotation(bid float64) *Quotation {
	return &Quotation{
		ID:  uuid.New().String(),
		Bid: bid,
	}
}

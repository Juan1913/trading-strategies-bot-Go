package models

type Order struct {
	ID        string
	Symbol    string
	Price     float64
	Quantity  float64
	Side      string
	Status    string
	Timestamp int64
}

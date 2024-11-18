package models

type Balance struct {
	Asset  string
	Free   float64
	Locked float64
}

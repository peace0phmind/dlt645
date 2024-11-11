package dlt645

import "github.com/shopspring/decimal"

type Value struct {
	Name  string
	Unit  string
	Value decimal.Decimal
}

type Client interface {
	ReadAddress() (string, error)
	Read(addr string, dic DIC) (*Value, error)
	BatchRead(addr string, dics []DIC) ([]*Value, error)
}

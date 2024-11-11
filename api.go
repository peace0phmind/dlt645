package dlt645

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type Value struct {
	Name  string
	Unit  string
	Value decimal.Decimal
	Err   error
}

func (v *Value) String() string {
	if v.Err != nil {
		return fmt.Sprintf("%s: %s", v.Name, v.Err)
	}
	return fmt.Sprintf("%s: %s%s", v.Name, v.Value, v.Unit)
}

type Client interface {
	ReadAddress() (string, error)
	Read(addr string, dic DIC) []*Value
	BatchRead(addr string, dics []DIC) []*Value
}

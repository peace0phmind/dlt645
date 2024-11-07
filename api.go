package dlt645

import "github.com/shopspring/decimal"

type Config struct {
	Protocol P
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.Protocol = PV2007
	return cfg
}

type Value struct {
	Name  string
	Unit  string
	Value decimal.Decimal
}

type Client interface {
	Read(addr string, dic DIC) (*Value, error)
	BatchRead(addr string, dics []DIC) ([]*Value, error)
}

type Packager interface {
	Encode(f *Frame) ([]byte, error)
	Decode(df []byte) (*Frame, error)
}

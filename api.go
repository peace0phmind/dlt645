package dlt645

import "github.com/shopspring/decimal"

type DLT645Config struct {
	Port string
}

type DLT645Value struct {
	Name  string
	Unit  string
	Value decimal.Decimal
}

type DLT645 interface {
	Read(dic DIC) (*DLT645Value, error)
	BatchRead(dic DIC) ([]*DLT645Value, error)
}

type Packager interface {
	Encode(f *Frame) ([]byte, error)
	Decode(df []byte) (*Frame, error)
}

// Transporter specifies the transport layer.
type Transporter interface {
	Send(request []byte) (response []byte, err error)
	CheckState() error
}

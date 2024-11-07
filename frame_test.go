package dlt645

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrame_SetAddress(t *testing.T) {
	tests := []struct {
		name            string
		address         string
		expectedAddress [6]byte
		isCompressed    bool
		expected        error
	}{
		{
			name:            "Valid_Address",
			address:         "1234567890",
			expectedAddress: [6]byte{0x90, 0x78, 0x56, 0x34, 0x12, 0x00},
			isCompressed:    false,
			expected:        nil,
		},
		{
			name:            "Valid_Address_1",
			address:         "123456789",
			expectedAddress: [6]byte{0x89, 0x67, 0x45, 0x23, 0x01, 0x00},
			isCompressed:    false,
			expected:        nil,
		},
		{
			name:            "Valid_compress_Address",
			address:         "123456789",
			expectedAddress: [6]byte{0x89, 0x67, 0x45, 0x23, 0xA1, 0xAA},
			isCompressed:    true,
			expected:        nil,
		},
		{
			name:     "Address_Exceeds_12_Bytes",
			address:  "1234567890123",
			expected: errors.New("address length is less than or equal to 12 bytes"),
		},
		{
			name:     "Address_With_Non_Numeric_Character",
			address:  "12345X789012",
			expected: errors.New("non-numeric character found at 12345X789012 index 5"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := &Frame{}
			err := frame.SetAddress(tt.address, tt.isCompressed)
			assert.Equal(t, tt.expected, err)
			if tt.expected == nil {
				assert.Equal(t, tt.expectedAddress, frame.Address)
				assert.Equal(t, tt.address, frame.GetAddress())
			}
		})
	}
}

func TestFrame_bcd(t *testing.T) {
	tests := []struct {
		value uint64
		buf   []byte
	}{
		{
			value: 1234567890,
			buf:   []byte{0x90, 0x78, 0x56, 0x34, 0x12, 0x00},
		},
		{
			value: 123456789,
			buf:   []byte{0x89, 0x67, 0x45, 0x23, 0x01, 0x00, 0x00},
		},
		{
			value: 12345678,
			buf:   []byte{0x78, 0x56, 0x34, 0x12},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.value), func(t *testing.T) {
			assert.Equal(t, tt.buf, uintToBcd(tt.value, len(tt.buf)))
			assert.Equal(t, tt.value, bcdToUint(tt.buf, len(tt.buf)))
		})
	}
}

func MustNewFromString(s string) (ret decimal.Decimal) {
	ret, _ = decimal.NewFromString(s)
	return ret
}

func TestFrame_GetValue(t *testing.T) {
	tests := []struct {
		buf              []byte
		dic              DIC
		exp              *Value
		expValueWithUnit string
	}{
		{
			buf: []byte{0x78, 0x56, 0x34, 0x12},
			dic: DICTotalActiveEnergy,
			exp: &Value{
				Name:  DICTotalActiveEnergy.Name(),
				Unit:  DICTotalActiveEnergy.Unit(),
				Value: MustNewFromString("123456.78"),
			},
			expValueWithUnit: "123456.78kWh",
		},
		{
			buf: []byte{0x01, 0x23},
			dic: DICPhaseAVoltage,
			exp: &Value{
				Name:  DICPhaseAVoltage.Name(),
				Unit:  DICPhaseAVoltage.Unit(),
				Value: MustNewFromString("230.1"),
			},
			expValueWithUnit: "230.1V",
		},
		{
			buf: []byte{0x23, 0x01, 0x00},
			dic: DICActiveConstant,
			exp: &Value{
				Name:  DICActiveConstant.Name(),
				Unit:  DICActiveConstant.Unit(),
				Value: MustNewFromString("123"),
			},
			expValueWithUnit: "123imp/kWh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.exp.Name, func(t *testing.T) {
			frame := &Frame{}
			value := frame.GetValue(tt.buf, tt.dic, PV2007)
			assert.Equal(t, tt.exp, value)
			assert.Equal(t, tt.expValueWithUnit, value.Value.String()+value.Unit)
		})
	}
}

func TestFrame_NewReadFrame(t *testing.T) {
	tests := []struct {
		addr     string
		dic      DIC
		protocol P
		expBytes []byte
		expErr   error
	}{
		{
			addr:     "1234567890",
			dic:      DICTotalActiveEnergy,
			protocol: PV2007,
			expBytes: []byte{0x68, 0x90, 0x78, 0x56, 0x34, 0x12, 0x0, 0x68, 0x11, 0x4, 0x33, 0x33, 0x33, 0x33, 0x55, 0x16},
			expErr:   nil,
		},
		{
			addr:     "12345678",
			dic:      DICPhaseAVoltage,
			protocol: PV2007,
			expBytes: []byte{0x68, 0x78, 0x56, 0x34, 0x12, 0x0, 0x0, 0x68, 0x11, 0x4, 0x33, 0x34, 0x34, 0x35, 0xc9, 0x16},
			expErr:   nil,
		},
		{
			addr:     "1234567",
			dic:      DICPhaseAVoltage,
			protocol: PV1997,
			expBytes: []byte{0x68, 0x67, 0x45, 0x23, 0x1, 0x0, 0x0, 0x68, 0x11, 0x2, 0x44, 0xe9, 0xe0, 0x16},
			expErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			f, err := NewReadFrame(tt.addr, tt.dic, tt.protocol)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expBytes, f.Bytes())
		})
	}
}

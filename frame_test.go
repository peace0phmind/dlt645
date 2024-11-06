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
			frame := &Frame{}
			assert.Equal(t, tt.buf, frame.uintToBcd(tt.value, len(tt.buf)))
			assert.Equal(t, tt.value, frame.bcdToUint(tt.buf, len(tt.buf)))
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
		exp              *DLT645Value
		expValueWithUnit string
	}{
		{
			buf: []byte{0x78, 0x56, 0x34, 0x12},
			dic: DICTotalActiveEnergy,
			exp: &DLT645Value{
				Name:  DICTotalActiveEnergy.Name(),
				Unit:  DICTotalActiveEnergy.Unit(),
				Value: MustNewFromString("123456.78"),
			},
			expValueWithUnit: "123456.78kWh",
		},
		{
			buf: []byte{0x01, 0x23},
			dic: DICPhaseAVoltage,
			exp: &DLT645Value{
				Name:  DICPhaseAVoltage.Name(),
				Unit:  DICPhaseAVoltage.Unit(),
				Value: MustNewFromString("230.1"),
			},
			expValueWithUnit: "230.1V",
		},
	}

	for _, tt := range tests {
		t.Run(tt.exp.Name, func(t *testing.T) {
			frame := &Frame{}
			value := frame.GetValue(tt.buf, tt.dic)
			assert.Equal(t, tt.exp, value)
			assert.Equal(t, tt.expValueWithUnit, value.Value.String()+value.Unit)
		})
	}
}

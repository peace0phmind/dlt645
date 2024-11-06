package dlt645

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"strings"
)

const (
	MAX_ADDRESS_LENGTH = 12
	DATA_MASK          = 0x33
)

var (
	PRE_BYTES = []byte{0xFE, 0xFE, 0xFE, 0xFE}
)

type Code byte

func NewCode(cc CCode) Code {
	return Code(cc.Val())
}

func (c Code) IsResponse() bool {
	return c&(1<<7) != 0
}

// 应答是否有错误
func (c Code) HasError() bool {
	return c&(1<<6) != 0
}

// 是否有后续数据帧
func (c Code) HasMore() bool {
	return c&(1<<5) != 0
}

type Frame struct {
	Start      byte    `value:"0x68"` // 帧起始符
	Address    [6]byte // 地址域
	FrameStart byte    `value:"0x68"` // 帧起始符
	C          Code    // 控制码
	L          byte    // 数据域长度
	Data       []byte  // 数据域
	CS         byte    // 校验码
	End        byte    `value:"0x16"` // 帧结束符
}

func (f *Frame) SetBroadcastAddress() {
	for i := range f.Address {
		f.Address[i] = 0x99
	}
}

func (f *Frame) SetAddress(address string, isCompressed bool) error {
	addrLen := len(address)

	if addrLen > MAX_ADDRESS_LENGTH {
		return fmt.Errorf("address length is less than or equal to %d bytes", MAX_ADDRESS_LENGTH)
	}

	for i := range f.Address {
		f.Address[i] = 0
	}

	for i := 0; i < addrLen; i++ {
		if address[i] < '0' || address[i] > '9' {
			return fmt.Errorf("non-numeric character found at %s index %d", address, i)
		}

		addrIndex := (addrLen - 1 - i) / 2
		if (addrLen-i)%2 == 0 {
			f.Address[addrIndex] |= (address[i] - '0') << 4
		} else {
			f.Address[addrIndex] |= address[i] - '0'
		}
	}

	if isCompressed && addrLen < MAX_ADDRESS_LENGTH {
		for i := addrLen + 1; i <= MAX_ADDRESS_LENGTH; i++ {
			addrIndex := (i - 1) / 2
			if i%2 == 0 {
				f.Address[addrIndex] |= 0xa0
			} else {
				f.Address[addrIndex] |= 0x0a
			}
		}
	}

	return nil
}

func (f *Frame) GetAddress() string {
	address := f.bcdToUint(f.Address[:], MAX_ADDRESS_LENGTH/2)
	return fmt.Sprintf("%d", address)
}

func (f *Frame) decimalDigits(value uint64) int {
	if value == 0 {
		return 1 // 特殊情况，0的位数是1
	}

	return int(math.Floor(math.Log10(float64(value)))) + 1
}

func (f *Frame) uintToBcd(value uint64, bufSize int) []byte {
	valueDigits := f.decimalDigits(value)
	if bufSize*2 < valueDigits {
		panic(fmt.Errorf("buffer size is less than %d bytes", valueDigits))
	}

	buf := make([]byte, bufSize)

	for i := 0; i < valueDigits; i++ {
		nibble := byte(value % 10)
		value = value / 10
		if i%2 != 0 {
			buf[i/2] |= nibble << 4
		} else {
			buf[i/2] |= nibble
		}
	}

	return buf
}

func (f *Frame) bcdToUint(data []byte, len int) (ret uint64) {
	base := uint64(0)

	for i := len*2 - 1; i >= 0; i-- {
		nibble := byte(0)
		if i%2 == 0 {
			nibble = data[i/2] & 0x0F
		} else {
			nibble = (data[i/2] >> 4) & 0x0F
		}

		// fix address compression prefix
		if nibble == 0x0A {
			nibble = 0
		}

		if base > 0 || nibble != 0 {
			ret = ret*base + uint64(nibble)
			if base == 0 {
				base = 10
			}
		}
	}

	return ret
}

func (f *Frame) DataAddMask() {
	for i := range f.Data {
		f.Data[i] += DATA_MASK
	}
}

func (f *Frame) DataCleanMask() {
	for i := range f.Data {
		f.Data[i] -= DATA_MASK
	}
}

func (f *Frame) GetValue(buf []byte, dic DIC) *DLT645Value {
	ret := &DLT645Value{}
	ret.Name = dic.Name()
	ret.Unit = dic.Unit()

	dotIndex := strings.Index(dic.Format(), ".")
	value := f.bcdToUint(buf, dic.Size())
	if dotIndex == -1 {
		ret.Value = decimal.NewFromUint64(value)
	} else {
		exp := int32(dic.Size()*2 - dotIndex)
		ret.Value = decimal.New(int64(value), -exp)
	}

	return ret
}

func (f *Frame) CalcCS() {
	f.CS = f.Start
	for _, a := range f.Address {
		f.CS += a
	}
	f.CS += f.FrameStart
	f.CS += byte(f.C)
	f.CS += f.L
	for _, d := range f.Data {
		f.CS += d
	}
}

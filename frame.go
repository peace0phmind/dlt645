package dlt645

import (
	"errors"
	"fmt"
	"github.com/expgo/factory"
	"github.com/shopspring/decimal"
	"strings"
)

const (
	MAX_ADDRESS_LENGTH = 12
	DATA_MASK          = 0x33
	FRAME_HEADER_LEN   = 10
	FrameStartByte     = 0x68
	FrameEndByte       = 0x16
)

var (
	PRE_BYTES = []byte{0xFE, 0xFE, 0xFE, 0xFE}
)

type Code byte

func NewCode(cc C) Code {
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
	Start   byte    `value:"0x68"` // 帧起始符
	Address [6]byte // 地址域
	AddrEnd byte    `value:"0x68"` // 帧起始符
	C       Code    // 控制码
	L       byte    // 数据域长度
	Data    []byte  // 数据域
	CS      byte    // 校验码
	End     byte    `value:"0x16"` // 帧结束符
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
	address := bcdToUint(f.Address[:], MAX_ADDRESS_LENGTH/2)
	return fmt.Sprintf("%d", address)
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

func (f *Frame) GetValue(buf []byte, dic DIC, protocol P) *Value {
	ret := &Value{}
	ret.Name = dic.Name()
	ret.Unit = dic.Unit()

	dotIndex := strings.Index(dic.Format(protocol), ".")
	value := bcdToUint(buf, dic.Size(protocol))
	if dotIndex == -1 {
		ret.Value = decimal.NewFromUint64(value)
	} else {
		exp := int32(dic.Size(protocol)*2 - dotIndex)
		ret.Value = decimal.New(int64(value), -exp)
	}

	return ret
}

func (f *Frame) _CalcCS() (cs byte) {
	cs = f.Start
	for _, a := range f.Address {
		cs += a
	}
	cs += f.AddrEnd

	cs += byte(f.C)
	cs += f.L
	for _, d := range f.Data {
		cs += d
	}

	return cs
}

func (f *Frame) CalcCS() {
	f.CS = f._CalcCS()
}

func (f *Frame) Bytes() (ret []byte) {
	ret = append(ret, f.Start)
	ret = append(ret, f.Address[:]...)
	ret = append(ret, f.AddrEnd, byte(f.C), f.L)
	ret = append(ret, f.Data...)
	ret = append(ret, f.CS, f.End)
	return ret
}

func (f *Frame) CheckStartError() error {
	if f.Start != FrameStartByte {
		return fmt.Errorf("invalid start byte: %x", f.Start)
	}

	if f.AddrEnd != FrameStartByte {
		return fmt.Errorf("invalid frame start byte: %x", f.AddrEnd)
	}

	return nil
}

func (f *Frame) CheckEndError() error {
	if f.CS != f._CalcCS() {
		return errors.New("cs error")
	}

	if f.End != FrameEndByte {
		return fmt.Errorf("invalid frame end byte: %x", f.End)
	}

	if f.C.HasError() {
		if f.L == 1 {
			return fmt.Errorf("frame has error: %d", f.Data[0])
		} else {
			return errors.New("frame has error, but data length not equals 1")
		}
	}

	return nil
}

func NewReadFrame(addr string, dic DIC, protocol P) (*Frame, error) {
	f := factory.New[Frame]()
	f.C = NewCode(CRD)
	if err := f.SetAddress(addr, false); err != nil {
		return nil, err
	}

	f.Data = dic.Code(protocol)
	f.L = byte(len(f.Data))
	f.DataAddMask()
	f.CalcCS()

	return f, nil
}

func NewFrameByRespHeader(header []byte) *Frame {
	if len(header) != FRAME_HEADER_LEN {
		panic(errors.New("header buffer length is not equal to 10"))
	}

	f := &Frame{
		Start:   header[0],
		Address: [6]byte{header[1], header[2], header[3], header[4], header[5], header[6]},
		AddrEnd: header[7],
		C:       Code(header[8] & 0x7f),
		L:       header[9],
	}

	return f
}

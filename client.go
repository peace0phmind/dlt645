package dlt645

import (
	"github.com/shopspring/decimal"
	"strings"
)

type client struct {
	Protocol    P
	transporter Transporter
}

func NewClient(transporter Transporter) Client {
	return &client{
		Protocol:    PV2007,
		transporter: transporter,
	}
}

func (mb *client) readFrame() (*Frame, error) {
	respBuf := make([]byte, FRAME_HEADER_LEN)
	_, err := mb.transporter.Read(respBuf)
	if err != nil {
		return nil, err
	}

	f := NewFrameByRespHeader(respBuf)
	if err = f.CheckStartError(); err != nil {
		return nil, err
	}

	if f.L > 0 {
		respData := make([]byte, f.L)
		_, err = mb.transporter.Read(respData)
		if err != nil {
			return nil, err
		}
		f.Data = respData
	}

	var endBuf [2]byte
	_, err = mb.transporter.Read(endBuf[:])
	if err != nil {
		return nil, err
	}
	f.CS = endBuf[0]
	f.End = endBuf[1]

	if err = f.CheckEndError(); err != nil {
		return nil, err
	}

	f.DataCleanMask()

	return f, nil
}

func (mb *client) Read(addr string, dic DIC) (*Value, error) {
	f, err := NewReadFrame(addr, dic, mb.Protocol)
	if err != nil {
		return nil, err
	}

	_, err = mb.transporter.Write(f.Bytes())
	if err != nil {
		return nil, err
	}

	respFrame, err1 := mb.readFrame()
	if err1 != nil {
		return nil, err1
	}

	return mb.getValue(respFrame.Data, dic), nil
}

func (mb *client) BatchRead(addr string, dics []DIC) ([]*Value, error) {
	return nil, nil
}

func (mb *client) getValue(buf []byte, dic DIC) *Value {
	ret := &Value{}
	ret.Name = dic.Name()
	ret.Unit = dic.Unit()

	dotIndex := strings.Index(dic.Format(mb.Protocol), ".")
	value := bcdToUint(buf, dic.Size(mb.Protocol))
	if dotIndex == -1 {
		ret.Value = decimal.NewFromUint64(value)
	} else {
		exp := int32(dic.Size(mb.Protocol)*2 - dotIndex)
		ret.Value = decimal.New(int64(value), -exp)
	}

	return ret
}

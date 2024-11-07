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

func (c *client) readFrame() (*Frame, error) {
	respBuf := make([]byte, FRAME_HEADER_LEN)
	_, err := c.transporter.Read(respBuf)
	if err != nil {
		return nil, err
	}

	f := NewFrameByRespHeader(respBuf)
	if err = f.CheckStartError(); err != nil {
		return nil, err
	}

	if f.L > 0 {
		respData := make([]byte, f.L)
		_, err = c.transporter.Read(respData)
		if err != nil {
			return nil, err
		}
		f.Data = respData
	}

	var endBuf [2]byte
	_, err = c.transporter.Read(endBuf[:])
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

func (c *client) Read(addr string, dic DIC) (*Value, error) {
	f, err := NewReadFrame(addr, dic, c.Protocol)
	if err != nil {
		return nil, err
	}

	_, err = c.transporter.Write(f.Bytes())
	if err != nil {
		return nil, err
	}

	respFrame, err1 := c.readFrame()
	if err1 != nil {
		return nil, err1
	}

	return c.getValue(respFrame.Data, dic), nil
}

func (c *client) BatchRead(addr string, dics []DIC) ([]*Value, error) {
	return nil, nil
}

func (c *client) getValue(buf []byte, dic DIC) *Value {
	ret := &Value{}
	ret.Name = dic.Name()
	ret.Unit = dic.Unit()

	dotIndex := strings.Index(dic.Format(c.Protocol), ".")
	value := bcdToUint(buf, dic.Size(c.Protocol))
	if dotIndex == -1 {
		ret.Value = decimal.NewFromUint64(value)
	} else {
		exp := int32(dic.Size(c.Protocol)*2 - dotIndex)
		ret.Value = decimal.New(int64(value), -exp)
	}

	return ret
}

package dlt645

import (
	"bytes"
	"errors"
	"github.com/expgo/factory"
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

func (c *client) writeFrame(f *Frame) error {
	var buf bytes.Buffer

	for i := 0; i < PRE_BYTE_LEN; i++ {
		_ = buf.WriteByte(PRE_BYTE)
	}

	_, err := buf.Write(f.Bytes())
	if err != nil {
		return err
	}

	_, err = c.transporter.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *client) readFrame() (*Frame, error) {
	respBuf := make([]byte, FRAME_HEADER_LEN+PRE_BYTE_LEN)
	_, err := c.transporter.Read(respBuf)
	if err != nil {
		return nil, err
	}

	f, err := NewFrameByRespHeader(respBuf)
	if err != nil {
		return nil, err
	}

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

	return f, nil
}

func (c *client) ReadAddress() (string, error) {
	f := factory.New[Frame]()
	f.C = NewCode(CRDA)
	if err := f.SetAddress("", true); err != nil {
		return "", err
	}
	f.CalcCS()

	if err := c.writeFrame(f); err != nil {
		return "", err
	}

	respFrame, err := c.readFrame()
	if err != nil {
		return "", err
	}

	return respFrame.GetAddress(), nil
}

func (c *client) getValue(buf []byte, dic DIC) (rets []*Value) {
	code := dic.Code(c.Protocol)

	if !bytes.Equal(buf[:len(code)], code) {
		return c.getErrorValues(dic, errors.New("dic code not equals"))
	}

	buf = buf[len(code):]
	_, dics := dic.CheckBlock(c.Protocol)

	for _, vDIC := range dics {
		v := &Value{}
		v.Name = vDIC.Name()
		v.Unit = vDIC.Unit()

		dotIndex := strings.Index(vDIC.Format(c.Protocol), ".")
		value := bcdToUint(buf, vDIC.Size(c.Protocol))
		if dotIndex == -1 {
			v.Value = decimal.NewFromUint64(value)
		} else {
			exp := int32(dic.Size(c.Protocol)*2 - dotIndex)
			v.Value = decimal.New(int64(value), -exp)
		}

		rets = append(rets, v)
		buf = buf[vDIC.Size(c.Protocol):]
	}

	return rets
}

func (c *client) getErrorValues(dic DIC, err error) (rets []*Value) {
	_, dics := dic.CheckBlock(c.Protocol)
	for _, vDIC := range dics {
		rets = append(rets, &Value{
			Name: vDIC.Name(),
			Unit: vDIC.Unit(),
			Err:  err,
		})
	}

	return rets
}

func (c *client) Read(addr string, dic DIC) []*Value {
	f, err := NewReadFrame(addr, dic, c.Protocol)
	if err != nil {
		return c.getErrorValues(dic, err)
	}

	if err1 := c.writeFrame(f); err1 != nil {
		return c.getErrorValues(dic, err1)
	}

	respFrame, err1 := c.readFrame()
	if err1 != nil {
		return c.getErrorValues(dic, err1)
	}

	return c.getValue(respFrame.Data, dic)
}

func (c *client) BatchRead(addr string, dics []DIC) (values []*Value) {
	for _, dic := range dics {
		values = append(values, c.Read(addr, dic)...)
	}
	return values
}

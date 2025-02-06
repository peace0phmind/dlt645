package dlt645

import (
	"errors"
	"github.com/expgo/factory"
	"github.com/tarm/serial"
	"io"
)

type SerialTransporter struct {
	baseTransporter
	conf *serial.Config
	port io.ReadWriteCloser
}

func NewSerialTransport(conf *serial.Config) *SerialTransporter {
	return factory.NewBeforeInit[SerialTransporter](func(ret *SerialTransporter) {
		ret.baseTransporter.addr = conf.Name
		ret.conf = conf
	})
}

func (t *SerialTransporter) Open() (err error) {
	if !t.running.CompareAndSwap(false, true) {
		return nil
	}

	if t.state == StateConnected {
		return nil
	}

	t.setState(StateConnecting, nil)

	t.port, err = serial.OpenPort(t.conf)
	if err != nil {
		t.L.Warnf("Open serial %s failed: %v", t.conf.Name, err)
		t.setState(StateDisconnected, err)
		return err
	}

	t.setState(StateConnected, nil)

	return err
}

func (t *SerialTransporter) Close() (err error) {
	defer func() {
		t.setState(StateConnectClosed, err)
		t.port = nil
	}()

	_ = t.baseTransporter.Close()

	if t.port == nil {
		return nil
	}

	return t.port.Close()
}

func (t *SerialTransporter) Write(data []byte) (n int, err error) {
	if t.port == nil || t.state == StateDisconnected {
		return 0, errors.New("serial transporter not connected")
	}

	defer func() {
		if err != nil {
			t.setState(StateDisconnected, err)
		}
	}()

	return t.port.Write(data)
}

func (t *SerialTransporter) Read(buf []byte) (n int, err error) {
	if t.port == nil || t.state == StateDisconnected {
		return 0, errors.New("serial transporter not connected")
	}

	defer func() {
		if err != nil {
			t.setState(StateDisconnected, err)
		}
	}()

	return t.port.Read(buf)
}

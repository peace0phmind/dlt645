package dlt645

import (
	"errors"
	"github.com/expgo/factory"
	"net"
	"time"
)

type TcpTransporter struct {
	baseTransporter
	conn net.Conn
}

func NewTcpTransport(addr string) *TcpTransporter {
	return factory.NewBeforeInit[TcpTransporter](func(ret *TcpTransporter) {
		ret.baseTransporter.addr = addr
	})
}

func (t *TcpTransporter) Open() (err error) {
	if !t.running.CompareAndSwap(false, true) {
		return nil
	}

	if t.state == StateConnected {
		return nil
	}

	t.setState(StateConnecting, nil)
	dialer := net.Dialer{Timeout: 3 * time.Second}
	t.conn, err = dialer.Dial("tcp", t.addr)
	if err != nil {
		t.L.Warnf("DialTCP %s failed: %v", t.addr, err)
		t.setState(StateDisconnected, err)
		return err
	}

	t.setState(StateConnected, nil)

	return err
}

func (t *TcpTransporter) Close() (err error) {
	defer func() {
		t.setState(StateConnectClosed, err)
		t.conn = nil
	}()

	_ = t.baseTransporter.Close()

	if t.conn == nil {
		return nil
	}

	return t.conn.Close()
}

func (t *TcpTransporter) Write(data []byte) (n int, err error) {
	if t.conn == nil || t.state == StateDisconnected {
		return 0, errors.New("tcp transporter not connected")
	}

	defer func() {
		if err != nil {
			t.setState(StateDisconnected, err)
		}
	}()

	err = t.conn.SetWriteDeadline(time.Now().Add(t.WriteTimeout))
	if err != nil {
		return 0, err
	}

	return t.conn.Write(data)
}

func (t *TcpTransporter) Read(buf []byte) (n int, err error) {
	if t.conn == nil || t.state == StateDisconnected {
		return 0, errors.New("tcp transporter not connected")
	}

	defer func() {
		if err != nil {
			t.setState(StateDisconnected, err)
		}
	}()

	err = t.conn.SetReadDeadline(time.Now().Add(t.ReadTimeout))
	if err != nil {
		return 0, err
	}

	return t.conn.Read(buf)
}

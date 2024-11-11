package dlt645

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTcpClient_ReadAddress(t *testing.T) {
	transport := NewTcpTransport("127.0.0.1:3330")
	defer transport.Close()

	c := NewClient(transport)

	err := transport.Open()
	assert.NoError(t, err)

	addr, err := c.ReadAddress()
	assert.NoError(t, err)

	t.Logf("addr: %s", addr)
}

func TestTcpClient_Read(t *testing.T) {
	transport := NewTcpTransport("127.0.0.1:3330")
	defer transport.Close()

	c := NewClient(transport)

	err := transport.Open()
	assert.NoError(t, err)

	v, err := c.Read("240727263614", DICTotalActiveEnergy)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	v, err = c.Read("240727263614", DICPositiveTotalActiveEnergy)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	v, err = c.Read("240727263614", DICNegativeTotalActiveEnergy)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	v, err = c.Read("240727263614", DICVoltage)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	v, err = c.Read("240727263614", DICCurrent)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	v, err = c.Read("240727263614", DICActivePower)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	v, err = c.Read("240727263614", DICReactivePower)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	v, err = c.Read("240727263614", DICFrequency)
	assert.NoError(t, err)
	t.Logf("value: %+v", v)

	_, err = c.Read("240727263614", DICLineVoltage)
	assert.EqualError(t, err, "frame has error: 无请求数据")
}

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

	v := c.Read("240727263614", DICTotalActiveEnergy)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICPositiveTotalActiveEnergy)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICNegativeTotalActiveEnergy)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICVoltage)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICCurrent)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICActivePower)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICReactivePower)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICFrequency)
	t.Logf("value: %+v", v)

	v = c.Read("240727263614", DICLineVoltage)
	t.Logf("value: %+v", v)
}

func TestTcpClient_BatchRead(t *testing.T) {
	transport := NewTcpTransport("127.0.0.1:3330")
	defer transport.Close()

	c := NewClient(transport)

	err := transport.Open()
	assert.NoError(t, err)

	v := c.BatchRead("240727263614", []DIC{DICTotalActiveEnergy, DICPositiveTotalActiveEnergy, DICNegativeTotalActiveEnergy,
		DICVoltage, DICCurrent, DICActivePower, DICReactivePower, DICFrequency, DICLineVoltage})
	t.Logf("value: %+v", v)
}

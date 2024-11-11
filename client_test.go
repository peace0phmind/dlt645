package dlt645

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTcpClient_ReadAddress(t *testing.T) {
	transport := NewTcpTransport("127.0.0.1:3330")

	c := NewClient(transport)

	err := transport.Open()
	assert.NoError(t, err)

	addr, err := c.ReadAddress()
	assert.NoError(t, err)

	t.Logf("addr: %s", addr)
}

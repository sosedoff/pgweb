package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientTakeTunnel(t *testing.T) {
	tunnel := &Tunnel{}

	oldClient := &Client{tunnel: tunnel}
	newClient := &Client{}

	newClient.TakeTunnel(oldClient)
	assert.Equal(t, tunnel, newClient.tunnel)
	assert.Nil(t, oldClient.tunnel)

	// Closing the old client must not tear down the tunnel it gave away
	assert.NoError(t, oldClient.Close())
	assert.Equal(t, tunnel, newClient.tunnel)

	// Nothing to take from a client without a tunnel, or from nil
	other := &Client{}
	newClient.TakeTunnel(other)
	newClient.TakeTunnel(nil)
	assert.Equal(t, tunnel, newClient.tunnel)
}

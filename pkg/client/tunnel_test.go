package client

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"testing"

	"github.com/sosedoff/pgweb/pkg/shared"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestMakeConfigWithKeyData(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	assert.NoError(t, err)

	block, err := ssh.MarshalPrivateKey(priv, "")
	assert.NoError(t, err)
	pemData := string(pem.EncodeToMemory(block))

	// Key contents take precedence over a (bogus) path
	cfg, err := makeConfig(&shared.SSHInfo{User: "u", Key: "/does/not/exist", KeyData: pemData})
	assert.NoError(t, err)
	assert.Len(t, cfg.Auth, 1)

	_, err = makeConfig(&shared.SSHInfo{User: "u", KeyData: "not a key"})
	assert.Error(t, err)
}

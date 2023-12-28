package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()

	assert.Equal(t, privKeyLen, len(privKey.Bytes()))

	pubKey := privKey.Public()
	assert.Equal(t, pubKeyLen, len(pubKey.Bytes()))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "f444a774398387299a73aebbde923f0d6d350d7ed41b2e7c628c5ffafb5847e0"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "136a1e56adf29e0b7cfab089d82ecc147450b8b3"
	)

	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	address := privKey.Public().Address()

	assert.Equal(t, addressStr, address.String())
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("hello world")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(msg, pubKey))

	// Test invalid message
	assert.False(t, sig.Verify([]byte("goodbye world"), pubKey))

	// Test invalid public key
	privKey2 := GeneratePrivateKey()
	pubKey2 := privKey2.Public()
	assert.False(t, sig.Verify(msg, pubKey2))
}

func TestPublicKeyAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()
	fmt.Println(address.String())
	assert.Equal(t, addressLen, len(address.Bytes()))
}

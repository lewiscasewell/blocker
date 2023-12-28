package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
)

const (
	privKeyLen = 64
	pubKeyLen  = 32
	seedLen    = 32
	addressLen = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {

	if len(seed) != seedLen {
		panic("invalid seed length")
	}

	return &PrivateKey{key: ed25519.NewKeyFromSeed(seed)}

}

func GeneratePrivateKey() *PrivateKey {
	seed := make([]byte, seedLen)

	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(message []byte) Signature {
	return Signature{
		value: ed25519.Sign(p.key, message),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	return &PublicKey{key: p.key.Public().(ed25519.PublicKey)}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Address() Address {
	return Address{value: p.key[:addressLen]}
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(message []byte, publicKey *PublicKey) bool {
	return ed25519.Verify(publicKey.key, message, s.value)
}

type Address struct {
	value []byte
}

func (a *Address) Bytes() []byte {
	return a.value
}

func (a *Address) String() string {
	return hex.EncodeToString(a.value)
}

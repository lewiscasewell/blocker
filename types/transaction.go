package types

import (
	"crypto/sha256"

	pb "github.com/golang/protobuf/proto"
	"github.com/lewiscasewell/blocker/crypto"
	"github.com/lewiscasewell/blocker/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		// Verify the signature
		sig := crypto.SignatureFromBytes(input.Signature)
		pubKey := crypto.PublicKeyFromBytes(input.PublicKey)

		if input.Signature != nil {
			input.Signature = nil
		}
		if !sig.Verify(HashTransaction(tx), pubKey) {
			return false
		}
	}

	return true
}

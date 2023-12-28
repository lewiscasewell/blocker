package util

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/lewiscasewell/blocker/proto"
)

func RandomHash() []byte {
	hash := make([]byte, 32)
	rand.Read(hash)
	return hash
}

func RandomBlock() *proto.Block {
	height, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}

	header := &proto.Header{
		Version:      1,
		Height:       int32(height.Int64()),
		PreviousHash: RandomHash(),
		MerkleRoot:   RandomHash(),
		Timestamp:    time.Now().UnixNano(),
	}

	return &proto.Block{
		Header: header,
	}
}

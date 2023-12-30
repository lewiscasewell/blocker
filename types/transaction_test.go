package types

import (
	"testing"

	"github.com/lewiscasewell/blocker/crypto"
	"github.com/lewiscasewell/blocker/proto"
	"github.com/lewiscasewell/blocker/util"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivKey.Public().Address()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddress := toPrivKey.Public().Address()

	input := &proto.TxInput{
		PrevTxHash:  util.RandomHash(),
		PrevTxIndex: 0,
		PublicKey:   fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress.Bytes(),
	}
	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress.Bytes(),
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	sig := SignTransaction(fromPrivKey, tx)

	input.Signature = sig.Bytes()

	assert.True(t, VerifyTransaction(tx))
}

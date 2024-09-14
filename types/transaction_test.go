package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wal99d/blockr/crypto"
	"github.com/wal99d/blockr/proto"
	"github.com/wal99d/blockr/util"
)

func TestNewTransaction(t *testing.T) {

	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	input := &proto.TxInput{

		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output := &proto.TxOutput{
		Amount:  5,
		Address: toPrivKey.Public().Bytes(),
	}

	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromPrivKey.Public().Bytes(),
	}

	tx := &proto.Transaction{

		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output, output2},
	}

	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()

	assert.True(t, VerifyTransaction(tx))
	fmt.Printf("%+v\n", tx)

}

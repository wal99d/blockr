package types

import (
	"crypto/sha256"

	pb "github.com/golang/protobuf/proto"
	"github.com/wal99d/blockr/crypto"
	"github.com/wal99d/blockr/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
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
		sig := crypto.SignatureFromByte(input.Signature)
		pubKey := crypto.PublicKeyFrombytes(input.PublicKey)
		//FIXME: make sure we don't run into problems after verification cause we have set the signature to nil.
		input.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
	}
	return true
}

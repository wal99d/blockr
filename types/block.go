package types

import (
	"crypto/sha256"

	pb "github.com/golang/protobuf/proto"
	"github.com/wal99d/blockr/crypto"
	"github.com/wal99d/blockr/proto"
)

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {

	return pk.Sign(HashBlock(b))
}

// HashBlock returns the sha256 of the block's header
func HashBlock(block *proto.Block) []byte {

	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}
	sha := sha256.Sum256(b)
	return sha[:]
}

package node

import (
	"encoding/hex"

	"github.com/wal99d/blockr/proto"
)

type BlockStorer interface{
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}

type Chain struct{
	blockStore BlockStorer
}

func NewChain(bs BlockStorer) *Chain {
	return &Chain{
		blockStore: bs,
	}
}

func (c *Chain) AddBlock(b *proto.Block) error {
	//validation
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error ) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error){
	return nil, nil
}

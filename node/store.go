package node

import (
	"fmt"

	"github.com/wal99d/blockr/proto"
	"go.uber.org/zap/internal/pool"
)


type BlockStorer interface{
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}

type MemoryBlockStore struct{
	blocks map[string]*proto.Block
}

func NewMemoryBlockStore() *MemoryBlockStore{
	return &MemoryBlockStore{
		blocks: make(map[string]*proto.Block),
	}
}

func (s *MemoryBlockStore) Get(hash string) (*proto.Block, error){
	block, ok := s.blocks[hash]
	if !ok{
		return nil, fmt.Errorf("block with hash [%s] doesn't exist", hash)
	}
	return block, nil
}

package common

import (
	"sync"
	"time"
)

type (
	BlockHeader struct {
		Index      string `json:"index"`
		Validator  string `json:"validator"`
		TimeStamp  string `json:"timestamp"`
		MerkleRoot string `json:"merkleRoot"`
	}
)

type BlockChain struct {
	ChainMetaData struct {
		Name      string `json:"name"`
		StartedAt string `json:"startedAt"`
	}
	Blocks []*Block
	Mu     sync.RWMutex
}

// NewBlockChain creates a new blockchain with genesis block
func NewBlockChain(name string) *BlockChain {
	chain := &BlockChain{
		ChainMetaData: struct {
			Name      string `json:"name"`
			StartedAt string `json:"startedAt"`
		}{
			Name:      name,
			StartedAt: time.Now().Format(time.RFC3339),
		},
		Blocks: make([]*Block, 0),
	}

	// Create and add genesis block
	genesisWallet := NewUserAccount("Genesis")
	gexTx := NewTx(genesisWallet.PublicKey, "0x1d6b85...", 10000)
	genesisWallet.SignTx(gexTx)

	genesisBlock := NewBlock([]*Tx{gexTx}, "")
	genesisBlock.SetValidator(genesisWallet)
	genesisBlock.SetMerkleTree()
	genesisBlock.Hash = genesisBlock.ProduceHash().Hex()
	genesisBlock.PrevHash = "0"
	genesisBlock.ValidatorSignature = NewValidator(genesisWallet).SignBlock(genesisBlock)

	chain.Mu.Lock()
	chain.Blocks = append(chain.Blocks, genesisBlock)
	chain.Mu.Unlock()

	return chain
}

// AddBlock adds a new block to the chain (thread-safe)
func (c *BlockChain) AddBlock(block *Block) (bool, error) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	//// Validate the block
	//if block == nil {
	//	return false, fmt.Errorf("nil block")
	//}
	//
	//if c.blockExists(block.Hash) {
	//	return false, fmt.Errorf("block with hash %s already exists", block.Hash)
	//}
	//
	//// For non-genesis blocks, check previous hash
	//if len(c.Blocks) > 0 && block.PrevHash != c.Blocks[len(c.Blocks)-1].Hash {
	//	return false, fmt.Errorf("previous hash mismatch")
	//}

	c.Blocks = append(c.Blocks, block)
	return true, nil
}

// GetBlockCount returns the number of blocks (thread-safe)
func (c *BlockChain) GetBlockCount() int {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return len(c.Blocks)
}

// GetLatestBlock returns the last block (thread-safe)
func (c *BlockChain) GetLatestBlock() *Block {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	if len(c.Blocks) == 0 {
		return nil
	}
	return c.Blocks[len(c.Blocks)-1]
}

// GetLatestHash returns the hash of the last block (thread-safe)
func (c *BlockChain) GetLatestHash() string {
	if block := c.GetLatestBlock(); block != nil {
		return block.Hash
	}
	return ""
}

// GetBlockByHash finds a block by hash (thread-safe)
func (c *BlockChain) GetBlockByHash(hash string) *Block {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	for _, block := range c.Blocks {
		if block.Hash == hash {
			return block
		}
	}
	return nil
}

// GetAllBlocks returns a copy of all blocks (thread-safe)
func (c *BlockChain) GetAllBlocks() []*Block {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	blocks := make([]*Block, len(c.Blocks))
	copy(blocks, c.Blocks)
	return blocks
}

// blockExists checks if a block exists (internal use, assumes lock is held)
func (c *BlockChain) blockExists(hash string) bool {
	for _, block := range c.Blocks {
		if block.Hash == hash {
			return true
		}
	}
	return false
}

// VerifyChain validates the entire chain integrity
func (c *BlockChain) VerifyChain() bool {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	if len(c.Blocks) == 0 {
		return false
	}

	// Check genesis block
	if c.Blocks[0].PrevHash != "0" {
		return false
	}

	// Check subsequent blocks
	for i := 1; i < len(c.Blocks); i++ {
		if c.Blocks[i].PrevHash != c.Blocks[i-1].Hash {
			return false
		}
	}

	return true
}

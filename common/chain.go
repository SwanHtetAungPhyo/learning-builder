package common

import (
	"fmt"
	"sync"
	"time"
)

type (
	BlockChain struct {
		ChainMetaData struct {
			Name      string `json:"name"`
			StartedAt string `json:"startedAt"`
		}
		Blocks []*Block     `json:"blocks"`
		mu     sync.RWMutex `json:"-"`
	}
	BlockHeader struct {
		Index      string `json:"index"`
		Validator  string `json:"validator"`
		TimeStamp  string `json:"timestamp"`
		MerkleRoot string `json:"merkleRoot"`
	}
)

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

	genesisWallet := NewUserAccount("Genesis")
	gexTx := NewTx(genesisWallet.PublicKey, "0x1d6b85...", 10000)
	genesisWallet.SignTx(gexTx)
	genesisBlock := NewBlock([]*Tx{gexTx}, "")
	genesisBlock.SetValidator(genesisWallet)
	genesisBlock.SetMerkleTree()
	genesisBlock.Hash = genesisBlock.ProduceHash().Hex()
	genesisBlock.PrevHash = "0"

	genesisValidator := NewValidator(genesisWallet)
	genesisBlock.SetValidator(genesisWallet)
	genesisValidator.ProduceBlock([]*Tx{gexTx}, genesisBlock.Hash)
	genesisBlock.Hash = genesisBlock.ProduceHash().Hex()
	genesisBlock.ValidatorSignature = genesisValidator.SignBlock(genesisBlock)
	genesisBlock.PrevHash = "0"
	if _, err := chain.AddBlock(genesisBlock); err != nil {
		panic(fmt.Sprintf("Failed to create genesis block: %v", err))
	}
	return chain
}

func (c *BlockChain) AddBlock(block *Block) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.checkHashExistence(block.Hash) {
		return false, fmt.Errorf("block with hash %s already exists", block.Hash)
	}

	if len(c.Blocks) > 0 && !c.checkHashExistence(block.PrevHash) {
		return false, fmt.Errorf("previous hash %s not found in chain", block.PrevHash)
	}

	index := c.findInsertIndex(block.Hash)
	c.Blocks = append(c.Blocks[:index], append([]*Block{block}, c.Blocks[index:]...)...)

	return true, nil
}

// findInsertIndex finds the position to insert a new block to maintain sorted order
func (c *BlockChain) findInsertIndex(hash string) int {
	left, right := 0, len(c.Blocks)
	for left < right {
		mid := (left + right) / 2
		if c.Blocks[mid].Hash < hash {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}

// checkHashExistence now assumes blocks are sorted
func (c *BlockChain) checkHashExistence(hash string) bool {
	left, right := 0, len(c.Blocks)
	for left < right {
		mid := (left + right) / 2
		if c.Blocks[mid].Hash == hash {
			return true
		} else if c.Blocks[mid].Hash < hash {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return false
}
func (c *BlockChain) GetBlockByIndex(hash string) *Block {
	c.mu.Lock()
	defer c.mu.Unlock()
	left, right := 0, len(c.Blocks)
	for left < right {
		index := (left + right) / 2
		if c.Blocks[index].Hash == hash {
			return c.Blocks[index]
		} else if c.Blocks[index].Hash > hash {
			right = index
		} else {
			left = index + 1
		}
	}
	return nil
}

func (c *BlockChain) GetLatestHash() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.Blocks) == 0 {
		return ""
	}
	return c.Blocks[len(c.Blocks)-1].Hash
}

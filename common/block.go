package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"time"
)

type (
	Block struct {
		BlockHeader        *BlockHeader `json:"blockHeader"`
		Hash               string       `json:"hash"`
		PrevHash           string       `json:"prevHash"`
		ValidatorSignature string       `json:"validatorSignature"`
		Txs                []*Tx        `json:"txs"`
	}
)

// NewBlock
func NewBlock(txs []*Tx, prevHash string) *Block {
	if txs == nil {
		txs = make([]*Tx, 0)
	}

	block := &Block{
		Txs:      txs,
		PrevHash: prevHash,
		BlockHeader: &BlockHeader{
			TimeStamp: time.Now().Format(time.RFC3339),
		},
	}
	block.Hash = block.ProduceHash().Hex()
	return block
}

// ValidateStructure
func (b *Block) ValidateStructure() error {
	if b.BlockHeader == nil {
		return errors.New("block header is nil")
	}
	if b.Hash == "" {
		return errors.New("block hash is empty")
	}
	if len(b.Txs) == 0 {
		return errors.New("block has no transactions")
	}
	return nil
}
func (b *Block) SetValidator(validator *UserAccount) {
	b.BlockHeader.Validator = validator.PublicKey
}
func (b *Block) ProduceHash() common.Hash {
	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, b.BlockHeader.Index)
	buffer.Write([]byte(b.PrevHash))
	buffer.Write([]byte(b.BlockHeader.TimeStamp))
	buffer.Write([]byte(b.BlockHeader.Validator))

	for _, tx := range b.Txs {
		buffer.Write(tx.HashTx().Bytes())
	}

	return crypto.Keccak256Hash(buffer.Bytes())
}
func (b *Block) SetMerkleTree() {
	txList := make([]merkletree.Content, len(b.Txs))
	for i, tx := range b.Txs {
		txList[i] = tx
	}

	treeMr, err := merkletree.NewTree(txList)
	if err != nil {
		log.Println("Cannot produce merkle tree:")
		log.Panic(err.Error())
		return
	}

	b.BlockHeader.MerkleRoot = treeMr.Root.String()
}

func (b *Block) VerifyBlockByMerkle() bool {
	var txList []merkletree.Content
	for _, tx := range b.Txs {
		txList = append(txList, tx)
	}
	treeMr, err := merkletree.NewTree(txList)
	if err != nil {
		log.Println("Cannot produce merkle tree:", err.Error())
		return false
	}
	rootHash := treeMr.Root.String()
	if rootHash != b.BlockHeader.MerkleRoot {
		log.Println("Merkle root mismatch:", rootHash, b.BlockHeader.MerkleRoot)
		return false
	}
	_, err = treeMr.VerifyTree()
	if err != nil {
		log.Println("Cannot verify merkle tree:", err.Error())
		return false
	}
	return true
}

func (b *Block) VerifyBlockBySignature() (bool, error) {
	sigBytes, err := hex.DecodeString(b.ValidatorSignature)
	if err != nil || len(sigBytes) != 65 {
		return false, errors.New("invalid signature")
	}

	pubKeyBytes, err := hex.DecodeString(b.BlockHeader.Validator)
	if err != nil || len(pubKeyBytes) != 65 {
		return false, errors.New("invalid uncompressed pubkey")
	}
	if len(pubKeyBytes) != 65 || pubKeyBytes[0] != 0x04 {
		return false, errors.New("public key must be uncompressed (65 bytes, 0x04 prefix)")
	}
	// Get the EXACT same hash used during signing
	messageHash := b.ProduceHash()

	// Only adjust V if your signatures use 27/28
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	return crypto.VerifySignature(pubKeyBytes, messageHash[:], sigBytes[:64]), nil
}

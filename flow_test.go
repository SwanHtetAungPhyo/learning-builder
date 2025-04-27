package main

// import (
// 	"encoding/hex"
// 	"github.com/SwanHtetAungPhyo/learning/common"
// 	"github.com/SwanHtetAungPhyo/learning/validator/producing"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// 	"time"
// )

// func TestUserAccountCreation(t *testing.T) {
// 	user := common.NewUserAccount("User1")
// 	assert.NotNil(t, user)
// 	assert.NotEmpty(t, user.PublicKey)
// 	assert.NotEmpty(t, user.GetPrivateKey())
// 	assert.Equal(t, "User1", user.Name)

// 	pubKeyBytes, err := hex.DecodeString(user.PublicKey)
// 	assert.Nil(t, err)
// 	assert.True(t, len(pubKeyBytes) == 64 || len(pubKeyBytes) == 65)
// }

// func TestTransactionCreation(t *testing.T) {
// 	tx := common.NewTx("User1", "User2", 100)
// 	assert.Equal(t, "User1", tx.From)
// 	assert.Equal(t, "User2", tx.To)
// 	assert.Equal(t, 100, tx.Amount)
// 	assert.NotEmpty(t, tx.Timestamp)

// 	// Verify timestamp is recent
// 	txTime, err := time.Parse(time.RFC3339, tx.Timestamp)
// 	assert.Nil(t, err)
// 	assert.WithinDuration(t, time.Now(), txTime, 5*time.Second)
// }

// func TestSignTransaction(t *testing.T) {
// 	user := common.NewUserAccount("User1")
// 	tx := common.NewTx(user.PublicKey, "User2", 100)
// 	signedTx := user.SignTx(tx)

// 	assert.NotEmpty(t, signedTx.Signature)
// 	assert.Equal(t, 130, len(signedTx.Signature))

// 	// Verify signature
// 	hash := tx.HashTx()
// 	sigBytes, err := hex.DecodeString(signedTx.Signature)
// 	assert.Nil(t, err)

// 	pubKeyBytes, err := hex.DecodeString(user.PublicKey)
// 	assert.Nil(t, err)

// 	valid := crypto.VerifySignature(pubKeyBytes, hash[:], sigBytes[:64])
// 	assert.True(t, valid)
// }

// func TestValidatorSigningBlock(t *testing.T) {
// 	user := common.NewUserAccount("User1")
// 	validator := producing.NewValidator(user)

// 	tx1 := common.NewTx(user.PublicKey, "User2", 100)
// 	tx2 := common.NewTx("User2", user.Name, 50)
// 	tx1 = user.SignTx(tx1)
// 	tx2 = user.SignTx(tx2)

// 	block := validator.ProduceBlock([]*common.Tx{tx1, tx2}, "")
// 	assert.NotNil(t, block.BlockHeader.MerkleTree)

// 	preSignHash := block.ProduceHash()
// 	postSignHash := block.ProduceHash()
// 	assert.Equal(t, preSignHash, postSignHash)

// 	assert.NotEmpty(t, block.ValidatorSignature)
// 	assert.NotEmpty(t, block.BlockHeader.Validator)

// 	isValid, err := block.VerifyBlockBySignature()
// 	assert.Nil(t, err)
// 	assert.True(t, isValid)
// 	assert.True(t, block.VerifyBlockByMerkle())
// }
// func TestValidatorProduceAndVerifyBlock(t *testing.T) {
// 	user1 := common.NewUserAccount("User1")
// 	user2 := common.NewUserAccount("User2")
// 	validator := producing.NewValidator(user1)

// 	tx1 := common.NewTx(user1.PublicKey, user2.PublicKey, 50)
// 	tx2 := common.NewTx(user2.PublicKey, user1.PublicKey, 30)
// 	tx1 = user1.SignTx(tx1)
// 	tx2 = user2.SignTx(tx2)

// 	block := validator.ProduceBlock([]*common.Tx{tx1, tx2}, "")

// 	// Verify Merkle tree
// 	assert.True(t, block.VerifyBlockByMerkle())

// 	// Verify signature
// 	isValid, err := block.VerifyBlockBySignature()
// 	assert.Nil(t, err)
// 	assert.True(t, isValid)

// 	// Verify transactions
// 	assert.Equal(t, 2, len(block.Txs))
// 	assert.Equal(t, tx1.Hash, block.Txs[0].Hash)
// 	assert.Equal(t, tx2.Hash, block.Txs[1].Hash)
// }

// func TestBlockChainAddBlock(t *testing.T) {
// 	blockchain := common.NewBlockChain("Test Blockchain")
// 	genesisBlock := blockchain.Blocks[0]

// 	assert.Equal(t, "0", genesisBlock.PrevHash)
// 	assert.NotEmpty(t, genesisBlock.Hash)
// 	assert.NotEmpty(t, genesisBlock.ValidatorSignature)

// 	user := common.NewUserAccount("User1")
// 	validator := producing.NewValidator(user)

// 	tx1 := common.NewTx(user.Name, "User2", 50)
// 	tx2 := common.NewTx("User2", user.Name, 30)
// 	tx1 = user.SignTx(tx1)
// 	tx2 = user.SignTx(tx2)

// 	block := validator.ProduceBlock([]*common.Tx{tx1, tx2}, genesisBlock.Hash)
// 	block.SetValidator(user)
// 	block.ValidatorSignature = validator.SignBlock(block)
// 	assert.Equal(t, genesisBlock.Hash, block.PrevHash)
// 	assert.NotEmpty(t, block.Hash)
// 	assert.NotEmpty(t, block.ValidatorSignature)

// 	added, err := blockchain.AddBlock(block)
// 	assert.Nil(t, err)
// 	assert.True(t, added)

// 	assert.Equal(t, 2, len(blockchain.Blocks))
// 	assert.Equal(t, genesisBlock.Hash, blockchain.Blocks[1].PrevHash)
// 	assert.Equal(t, block.Hash, blockchain.Blocks[1].Hash)

// 	_, err = blockchain.AddBlock(block)
// 	assert.NotNil(t, err)
// }

// func TestBlockChainGetLatestHash(t *testing.T) {
// 	blockchain := common.NewBlockChain("Test Blockchain")
// 	genesisHash := blockchain.GetLatestHash()
// 	assert.Equal(t, blockchain.Blocks[0].Hash, genesisHash)

// 	user := common.NewUserAccount("User1")
// 	validator := producing.NewValidator(user)

// 	tx1 := common.NewTx(user.Name, "User2", 50)
// 	tx2 := common.NewTx("User2", user.Name, 30)
// 	tx1 = user.SignTx(tx1)
// 	tx2 = user.SignTx(tx2)

// 	block := validator.ProduceBlock([]*common.Tx{tx1, tx2}, genesisHash)
// 	_, err := blockchain.AddBlock(block)
// 	if err != nil {
// 		return
// 	}

// 	latestHash := blockchain.GetLatestHash()
// 	assert.Equal(t, block.Hash, latestHash)
// 	assert.NotEqual(t, genesisHash, latestHash)
// }

// func TestBlockChainGetBlockByIndex(t *testing.T) {
// 	blockchain := common.NewBlockChain("Test Blockchain")
// 	genesisBlock := blockchain.Blocks[0]

// 	foundBlock := blockchain.GetBlockByIndex(genesisBlock.Hash)
// 	assert.Equal(t, genesisBlock, foundBlock)

// 	notFound := blockchain.GetBlockByIndex("nonexistent")
// 	assert.Nil(t, notFound)
// }

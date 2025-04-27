package common

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

type Validator struct {
	Wallet *UserAccount
	Stake  int64 `json:"stake"`
}

func NewValidator(wallet *UserAccount) *Validator {
	validator := &Validator{
		Wallet: wallet,
		Stake:  0,
	}
	return validator
}
func (v *Validator) AddStake(amount int) {
	v.Stake += int64(amount)
}
func (v *Validator) SubtractStake(amount int) {
	v.Stake -= int64(amount)
}
func (v *Validator) GetStake() int64 {
	return v.Stake
}
func (v *Validator) SignBlock(block *Block) string {

	validatorPrivateKey := Must[*ecdsa.PrivateKey](crypto.HexToECDSA(v.Wallet.privateKey))
	blockHash := block.ProduceHash()
	sigByte := Must[[]byte](crypto.Sign(blockHash[:], validatorPrivateKey))
	if len(sigByte) != 65 {
		log.Panic("Invalid signature length:", len(sigByte))
		return " "
	}
	return hex.EncodeToString(sigByte)
}
func (v *Validator) VerifyBlock(block *Block) bool {
	panic("Impl me")
}

func (v *Validator) ProduceBlock(tx []*Tx, prevHash string) *Block {
	newBlockProposal := NewBlock(tx, prevHash)
	newBlockProposal.SetMerkleTree()
	return newBlockProposal
}

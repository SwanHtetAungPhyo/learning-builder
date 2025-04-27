package producing

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Validator struct {
	Wallet *common.UserAccount
	Stake  int64 `json:"stake"`
}

func NewValidator(wallet *common.UserAccount) *Validator {
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
func (v *Validator) SignBlock(block *common.Block) string {
	validatorPrivateKey := common.Must[*ecdsa.PrivateKey](crypto.HexToECDSA(v.Wallet.GetPrivateKey()))
	blockHash := block.ProduceHash()

	// Sign the hash (Ethereum adds the recovery ID automatically)
	sigBytes := common.Must[[]byte](crypto.Sign(blockHash[:], validatorPrivateKey))

	// Ensure the signature is 65 bytes [R || S || V]
	if len(sigBytes) != 65 {
		panic("invalid Ethereum signature length")
	}

	// Store the uncompressed public key (required for VerifySignature)
	pubKey := crypto.FromECDSAPub(&validatorPrivateKey.PublicKey)
	block.BlockHeader.Validator = hex.EncodeToString(pubKey) // 65 bytes (0x04...)

	return hex.EncodeToString(sigBytes)
}

func (v *Validator) ProduceBlock(tx []*common.Tx, prevHash string) *common.Block {
	newBlockProposal := common.NewBlock(tx, prevHash)
	newBlockProposal.SetMerkleTree()
	newBlockProposal.SetValidator(v.Wallet)
	return newBlockProposal
}

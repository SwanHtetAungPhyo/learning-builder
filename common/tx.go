package common

import (
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goccy/go-json"
	"strconv"
	"strings"
	"time"
)

type (
	Tx struct {
		From      string `json:"from"`
		To        string `json:"to"`
		Signature string `json:"signature"`
		Amount    int    `json:"amount"`
		Timestamp string `json:"timestamp"`
		PrevHash  string `json:"prevHash"`
		Hash      string `json:"hash"`
	}

	UserAccount struct {
		Name       string `json:"name"`
		PublicKey  string `json:"publicKey"`
		privateKey string `json:"_"`
		CreatedAt  string `json:"createdAt"`
		Balance    int    `json:"balance"`
	}
)

var _ merkletree.Content = (*Tx)(nil)

func NewTx(from string, to string, amount int) *Tx {
	tx := &Tx{
		From:      from,
		To:        to,
		Amount:    amount,
		Signature: "",
		Timestamp: time.Now().Format(time.RFC3339),
		PrevHash:  "adffdsafads",
	}
	tx.Hash = tx.HashTx().Hex()
	return tx
}

func (t *Tx) messageReconstruction() string {
	jsonData := Must[[]byte](json.Marshal(t))
	return string(jsonData)
}
func (t *Tx) MessageToSign() string {
	var msgToSign strings.Builder
	for _, field := range []string{t.From, t.To, strconv.Itoa(t.Amount), t.Timestamp, t.PrevHash, t.Hash} {
		msgToSign.WriteString(field)
	}
	return msgToSign.String()
}

func (t *Tx) HashTx() common.Hash {
	if t == nil {
		return common.Hash{}
	}
	return crypto.Keccak256Hash([]byte(t.MessageToSign()))
}

// CalculateHash This is for merkle tree
func (t *Tx) CalculateHash() ([]byte, error) {
	return t.HashTx().Bytes(), nil
}

func (t *Tx) Equals(other merkletree.Content) (bool, error) {
	otherTx, ok := other.(*Tx)
	if !ok {
		return false, nil
	}
	return t.From == otherTx.From && t.To == otherTx.To && t.Amount == otherTx.Amount, nil
}

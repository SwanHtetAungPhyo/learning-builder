package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goccy/go-json"
	"strconv"
	"strings"
)

type (
	Tx struct {
		From      string `json:"from"`
		To        string `json:"to"`
		Signature string `json:"signature"`
		Amount    int    `json:"amount"`
	}

	UserAccount struct {
		Name       string `json:"name"`
		PublicKey  string `json:"publicKey"`
		privateKey string `json:"_"`
		CreatedAt  string `json:"createdAt"`
		Balance    int    `json:"balance"`
	}
)

func NewTx(from string, to string, amount int) *Tx {
	return &Tx{
		From:      from,
		To:        to,
		Amount:    amount,
		Signature: "",
	}
}

func (t *Tx) messageReconstruction() string {
	jsonData := Must[[]byte](json.Marshal(t))
	return string(jsonData)
}
func (t *Tx) MessageToSign() string {
	var msgToSign strings.Builder
	msgToSign.WriteString(t.From)
	msgToSign.WriteString(t.To)
	msgToSign.WriteString(strconv.Itoa(t.Amount))
	return msgToSign.String()
}

func (t *Tx) HashTx() common.Hash {
	if t == nil {
		return common.Hash{}
	}
	return crypto.Keccak256Hash([]byte(t.MessageToSign()))
}

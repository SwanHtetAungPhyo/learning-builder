package main

import (
	"errors"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

var _ merkletree.Content = (*Tx)(nil)

type Tx struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Signature string `json:"signature"`
	Amount    int    `json:"amount"`
	Timestamp string `json:"timestamp"`
	PrevHash  string `json:"prevHash"`
	Hash      string `json:"hash"`
}

func (t Tx) CalculateHash() ([]byte, error) {
	txHash := crypto.Keccak256Hash([]byte(t.String()))
	return txHash[:], nil
}

func (t Tx) Equals(other merkletree.Content) (bool, error) {
	otherTx, ok := other.(Tx)
	if !ok {
		return false, errors.New("type mismatch: not a Tx")
	}
	return t.From == otherTx.From && t.To == otherTx.To && t.Amount == otherTx.Amount, nil
}
func (t Tx) String() string {
	var txBuilder strings.Builder
	for _, field := range []string{t.From, t.To, string(rune(t.Amount)), t.Timestamp, t.PrevHash, t.Hash, t.Signature, "\n"} {
		txBuilder.WriteString(field)
	}
	return txBuilder.String()
}

func main() {
	txList := []merkletree.Content{
		Tx{From: "Alice", To: "Bob", Amount: 10, Timestamp: "2025-04-26T10:00:00", Signature: "sig1", PrevHash: "0xabc", Hash: "0xdef"},
		Tx{From: "Bob", To: "Charlie", Amount: 5, Timestamp: "2025-04-26T10:01:00", Signature: "sig2", PrevHash: "0xdef", Hash: "0xghi"},
	}

	tree, err := merkletree.NewTree(txList)
	if err != nil {
		panic(err.Error())
	}
	root := tree.MerkleRoot()
	fmt.Println(root)

	valid, err := tree.VerifyTree()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(valid)
}

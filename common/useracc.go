package common

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func NewUserAccount(name string) *UserAccount {
	userAcc := &UserAccount{Name: name}
	userAcc.produceWalletKeys()
	userAcc.Balance = 0
	userAcc.CreatedAt = time.Now().Format(time.RFC3339)
	return userAcc
}

func (u *UserAccount) AddBalance(amount int) {
	u.Balance += amount
}
func (u *UserAccount) SubtractBalance(amount int) {
	u.Balance -= amount
}

func (u *UserAccount) produceWalletKeys() {
	key := Must[*ecdsa.PrivateKey](crypto.GenerateKey())

	pubKeyString := hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey))
	privateKeyString := hex.EncodeToString(crypto.FromECDSA(key))
	u.PublicKey = pubKeyString
	u.privateKey = privateKeyString
}

func (u *UserAccount) SignTx(tx *Tx) *Tx {
	if tx == nil {
		return nil
	}

	toECDSA := Must[*ecdsa.PrivateKey](crypto.HexToECDSA(u.privateKey))
	messageHash := tx.HashTx()
	sigByte := Must[[]byte](crypto.Sign(messageHash[:], toECDSA))

	if len(sigByte) != 65 {
		log.Println("Invalid signature length:", len(sigByte))
		return nil
	}

	tx.Signature = hex.EncodeToString(sigByte)
	return tx
}
func (u *UserAccount) CommunicateWithRPC(tx *Tx) {
	//newRpcClient := NewHttpClientWithTimeout(time.Second*10, RpcUrl)

	restyClient := resty.New()
	req := restyClient.R()
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(tx)
	post, err := req.Post("http://localhost:8545/submit")
	if err != nil {
		log.Println(err.Error())
		return
	}
	if post.StatusCode() != 200 {
		log.Println(post.StatusCode(), post.Status())
		return
	}

	log.Println(post.StatusCode(), post.Status(), post.Body())

	//req, err := http.NewRequest("POST", RpcUrl, bytes.NewBuffer(txData))
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept", "application/json")
	//if err != nil {
	//	panic(err.Error())
	//	return
	//}
	//resp, err := newRpcClient.httpClient.Post("http://localhost:8545/submit", "application/json", req.Body)
	//log.Println(resp.StatusCode, resp.Status, resp.Body)
	//if err != nil {
	//	panic(err.Error())
	//	return
	//}
	//err = resp.Body.Close()
	//if err != nil {
	//	panic(err.Error())
	//	return
	//}
	//log.Println(resp.StatusCode, resp.Status, resp.Body)
}

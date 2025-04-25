package main

import (
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goccy/go-json"
)

const (
	rabbitMQURL  = "amqp://guest:guest@localhost:5672/"
	exchangeName = "transactions"
)

func main() {
	client := common.NewRabbitMQClient(rabbitMQURL).
		Connect().
		CreateChannel()
	defer client.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go startValidator(client, "validator1", done, &wg)
	wg.Add(1)
	go startValidator(client, "validator2", done, &wg)

	<-signalChan
	log.Println("Shutting down...")
	close(done)
	wg.Wait()
	log.Println("Shutdown completed successfully")
}

func startValidator(client *common.RabbitMQClient, queueName string, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	msgs := client.ConsumeMsgWithKey(queueName, exchangeName, queueName+"key")
	messageChan := make(chan []byte)

	go func() {
		defer close(messageChan)
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					return
				}
				messageChan <- d.Body
			case <-done:
				return
			}
		}
	}()

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}

			var tx common.Tx
			if err := json.Unmarshal(msg, &tx); err != nil {
				log.Printf("[%s] Invalid JSON: %v", queueName, err)
				continue
			}

			if VerifyTx(&tx) {
				log.Printf("[%s] ✅ Verified TX: %s → %s", queueName, tx.From, tx.To)
			} else {
				log.Printf("[%s] ❌ Unverified TX: %s → %s", queueName, tx.From, tx.To)
			}

		case <-done:
			return
		}
	}
}

func VerifyTx(tx *common.Tx) bool {
	messageHash := tx.HashTx()
	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}

	publicKeyBytes, err := hex.DecodeString(tx.From)
	if err != nil {
		return false
	}

	if len(sigBytes) != 65 {
		return false
	}

	return crypto.VerifySignature(publicKeyBytes, messageHash[:], sigBytes[:64])
}

// client/main.go
package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	rabbitMQURL  = "amqp://guest:guest@localhost:5672/"
	exchangeName = "transactions"
	queueName    = "2061fcaf013131a753bac07e10cdf46eae95cb96bbbfcdbd7564667fc350db62"

	serverAddress = "127.0.0.1:8081"
	batchSize     = 10
	batchTimeout  = 500 * time.Millisecond
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

	validTxs := make(chan *common.Tx, 100)
	blocksToPropose := make(chan *common.Block, 10)
	validator := common.NewValidator(common.NewUserAccount(uuid.New().String()))
	wg.Add(1)
	go startValidator(client, done, validTxs, &wg)

	wg.Add(1)
	go batchTxs(done, validator, validTxs, blocksToPropose, &wg)

	//wg.Add(1)

	//go DialAndPropose(done, blocksToPropose, &wg)

	<-signalChan
	log.Println("Shutting down client...")
	close(done)
	wg.Wait()
	log.Println("Client shutdown completed")
}

func startValidator(client *common.RabbitMQClient, done <-chan struct{}, validTxs chan<- *common.Tx, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(validTxs)

	var routingKey = common.RoutingKeyCalculator(queueName)
	msgs := client.ConsumeMsgWithKey(queueName, exchangeName, routingKey)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				return
			}

			var tx common.Tx
			if err := json.Unmarshal(d.Body, &tx); err != nil {
				log.Printf("Invalid JSON: %v", err)
				continue
			}

			if VerifyTx(&tx) {
				validTxs <- &tx
				log.Printf("✅ Verified TX: %s → %s", tx.From, tx.To)
			} else {
				log.Printf("❌ Unverified TX: %s → %s", tx.From, tx.To)
			}

		case <-done:
			return
		}
	}
}

func batchTxs(done <-chan struct{}, validator *common.Validator, validTxs <-chan *common.Tx, blocks chan<- *common.Block, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(blocks)

	var batch []*common.Tx
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()

	for {
		select {
		case tx, ok := <-validTxs:
			if !ok {
				return
			}

			batch = append(batch, tx)
			if len(batch) >= batchSize {
				proposeBlock(batch, validator, blocks)
				batch = nil
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(batchTimeout)
			}

		case <-done:
			if len(batch) > 0 {
				proposeBlock(batch, validator, blocks)
			}
			return
		}
	}
}

func proposeBlock(txs []*common.Tx, validator *common.Validator, blocks chan<- *common.Block) {
	block := validator.ProduceBlock(txs, "FDSFDSFDSFDFSD")
	blocks <- block
	log.Printf("📦 Proposed block with %d TXs (Hash: %s)", len(txs), block.Hash)
}

func DialAndPropose(done <-chan struct{}, blocks <-chan *common.Block, wg *sync.WaitGroup) {
	defer wg.Done()

	var conn net.Conn

	for {
		select {
		case <-done:
			return
		default:
			var err error
			conn, err = net.Dial("tcp", serverAddress)
			err = conn.SetDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return
			}
			if err == nil {
				err := conn.Close()
				if err != nil {
					return
				}
				break
			}
			log.Printf("Connection failed, retrying...: %v", err)
			time.Sleep(1 * time.Second)
		}
		if conn != nil {
			break
		}
	}

	encoder := json.NewEncoder(conn)
	encoder.SetIndent("", "")

	for {
		select {
		case block, ok := <-blocks:
			if !ok {
				return
			}

			if block.BlockHeader == nil {
				block.BlockHeader = &common.BlockHeader{
					TimeStamp: time.Now().Format(time.RFC3339),
				}
			}

			if err := encoder.Encode(block); err != nil {
				log.Printf("Failed to send block: %v", err)
				conn.Close()
				conn = nil
				for conn == nil {
					select {
					case <-done:
						return
					default:
						newConn, err := net.Dial("tcp", serverAddress)
						if err == nil {
							conn = newConn
							encoder = json.NewEncoder(conn)
							encoder.SetIndent("", "")
							break
						}
						time.Sleep(1 * time.Second)
					}
				}
			}

		case <-done:
			return
		}
	}
}
func VerifyTx(tx *common.Tx) bool {
	messageHash := tx.HashTx()
	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil || len(sigBytes) != 65 {
		return false
	}

	publicKeyBytes, err := hex.DecodeString(tx.From)
	if err != nil {
		return false
	}

	return crypto.VerifySignature(publicKeyBytes, messageHash[:], sigBytes[:64])
}

func CloseConnection(conn net.Conn) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Printf("Connection close error: %v", err)
		}
	}
}

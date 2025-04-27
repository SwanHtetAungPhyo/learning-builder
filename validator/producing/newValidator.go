package producing

//import (
//	"github.com/SwanHtetAungPhyo/learning/common"
//	"github.com/goccy/go-json"
//	"github.com/google/uuid"
//	"log"
//	"os"
//	"os/signal"
//	"sync"
//	"syscall"
//)
//
//var txsToBeInNewBlocks = make(chan *common.Tx)
//var proposalChannel = make(chan *common.Block)
//
//func main() {
//	client := common.NewRabbitMQClient(rabbitMQURL).
//		Connect().
//		CreateChannel()
//	defer client.Close()
//	signalChan := make(chan os.Signal, 1)
//	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
//	validatorWallet := common.NewUserAccount(uuid.New().String())
//	validator := common.NewValidator(validatorWallet)
//	var wg sync.WaitGroup
//	done := make(chan struct{})
//	wg.Add(1)
//	go StartValidator(client, "validator1", done, &wg, validator)
//
//}
//
//func StartValidator(client *common.RabbitMQClient, queueName string, done <-chan struct{}, wg *sync.WaitGroup, validator *common.Validator) {
//	defer wg.Done()
//
//	msgs := client.ConsumeMsgWithKey(queueName, exchangeName, queueName+"key")
//	messageChan := make(chan []byte)
//
//	go func() {
//		defer close(messageChan)
//		for {
//			select {
//			case d, ok := <-msgs:
//				if !ok {
//					return
//				}
//				messageChan <- d.Body
//			case <-done:
//			}
//		}
//	}()
//
//	for {
//		select {
//		case msg, ok := <-messageChan:
//			if !ok {
//				return
//			}
//
//			var tx common.Tx
//			if err := json.Unmarshal(msg, &tx); err != nil {
//				log.Printf("[%s] Invalid JSON: %v", queueName, err)
//				continue
//			}
//
//			if VerifyTx(&tx) {
//				log.Printf("[%s] ✅ Verified TX: %s → %s", queueName, tx.From, tx.To)
//				txsToBeInNewBlocks <- &tx
//			} else {
//				log.Printf("[%s] ❌ Unverified TX: %s → %s", queueName, tx.From, tx.To)
//			}
//		case verifiedTx := <-txsToBeInNewBlocks:
//			var txBatch []*common.Tx
//			batchSize := 10
//			if len(txBatch) < batchSize {
//				txBatch = append(txBatch, verifiedTx)
//			}
//			log.Printf("Producing block with %d transactions", len(txBatch))
//			proposalChannel <- validator.ProduceBlock(txBatch)
//		//case proposal := <-proposalChannel:
//		//	proposal
//		case <-done:
//			return
//		}
//	}
//}

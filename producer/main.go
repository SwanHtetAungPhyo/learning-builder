package main

import (
	"github.com/SwanHtetAungPhyo/learning/common"
	"log"
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

	// Create exchange and queues
	client.CreateExchange(exchangeName, "direct")
	client.CreateQueue("validator1")
	client.CreateQueue("validator2")
	client.BindQueueToExchange("validator1", exchangeName, "validator1key")
	client.BindQueueToExchange("validator2", exchangeName, "validator2key")

	alice := common.NewUserAccount("Alice")
	bob := common.NewUserAccount("Bob")
	alice.AddBalance(100)

	//for i := 0; i < 10; i++ {
	tx := common.NewTx(alice.PublicKey, bob.PublicKey, 10)
	tx = alice.SignTx(tx)
	alice.SubtractBalance(10)
	if tx != nil {
		alice.CommunicateWithRPC(tx)
		return
	}
	log.Println("Helll ")
	//if tx == nil {
	//	log.Println("Signing transaction failed")
	//	continue
	//}
	//
	//	txBytes := common.Must[[]byte](json.Marshal(tx))
	//
	//	client.SendMsgJson(txBytes, exchangeName, "validator1key")
	//	client.SendMsgJson(txBytes, exchangeName, "validator2key")
	//
	//	log.Printf("✅ Transaction %d sent to both queues", i)
	//	time.Sleep(500 * time.Millisecond)
	//}
}

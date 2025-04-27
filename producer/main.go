package main

import (
	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	rabbitMQURL  = "amqp://guest:guest@localhost:5672/"
	exchangeName = "transactions"
)

func main() {
	// Initialize accounts
	alice := common.NewUserAccount("Alice")
	bob := common.NewUserAccount("Bob")
	alice.AddBalance(100)

	// Initial batch of transactions
	for i := 0; i < 10; i++ {
		amount := 10*i + 1
		if alice.Balance < amount {
			log.Println("Insufficient balance for initial transactions")
			break
		}

		tx := common.NewTx(alice.PublicKey, bob.PublicKey, amount)
		tx = alice.SignTx(tx)
		if tx != nil {
			alice.SubtractBalance(amount)
			alice.CommunicateWithRPC(tx)
			log.Printf("Sent initial transaction %d: %d", i, amount)
		} else {
			log.Println("Failed to sign transaction")
		}
		time.Sleep(100 * time.Millisecond)
	}
	c := cron.New()
	_, err := c.AddFunc("@every 1s", func() {
		amount := 10
		if alice.Balance < amount {
			log.Println("Insufficient balance for recurring transaction")
			alice.AddBalance(100)
			return
		}

		tx := common.NewTx(alice.PublicKey, bob.PublicKey, amount)
		tx = alice.SignTx(tx)
		if tx != nil {
			alice.SubtractBalance(amount)
			alice.CommunicateWithRPC(tx)
		} else {
			log.Println("Failed to sign recurring transaction")
		}
	})

	if err != nil {
		log.Fatal("Error adding cron function:", err)
	}

	c.Start()
	defer c.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Println("Cron job started. Press Ctrl+C to stop...")
	<-sigChan
	log.Println("Shutting down...")
}

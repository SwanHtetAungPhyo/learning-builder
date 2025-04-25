package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"sync"

	"github.com/SwanHtetAungPhyo/learning/common"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "validator1"
)

var consumerChannel = make(chan []byte)

func main() {
	client := common.NewRabbitMQClient(rabbitMQURL).
		Connect().
		CreateChannel()

	client.CreateQueue(queueName)

	// Register consumer
	msgs := common.MustEnhance(func() (<-chan amqp.Delivery, error) {
		return client.ConsumeMsg(queueName), nil
	}, "Registering consumer failed")

	// Start receiving messages
	go func() {
		for d := range msgs {
			consumerChannel <- d.Body
		}
	}()

	// Start message processing
	var wg sync.WaitGroup
	wg.Add(1)
	outChan := ProcessingMessage(consumerChannel, &wg)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for processedMsg := range outChan {
			fmt.Printf("Processed message: %s\n", processedMsg)
		}
	}()

	<-sigChan
	close(consumerChannel)
	wg.Wait()
	client.CloseConnection()
}

func ProcessingMessage(consumerChannel chan []byte, wg *sync.WaitGroup) <-chan []byte {
	defer wg.Done()
	type Message struct {
		Length          int     `json:"length"`
		SizeInMB        float64 `json:"SizeInMB"`
		OriginalMessage []byte  `json:"originalMessage"`
	}
	outChan := make(chan []byte)
	go func() {
		for message := range consumerChannel {
			msgToOut := &Message{
				Length:          len(string(message)),
				SizeInMB:        float64(len(message)) / (1024 * 1024),
				OriginalMessage: message,
			}
			outChan <- common.MustEnhance(func() ([]byte, error) {
				return json.Marshal(msgToOut)
			}, "Marshaling message failed")
		}
		close(outChan)
	}()
	return outChan
}

package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"sync"
)

var (
	balance    int
	mu         sync.Mutex
	wg         sync.WaitGroup
	deHistory  []int
	subHistory []int
)

func add(in chan int) {
	defer wg.Done()
	for i := range in {
		log.Println("Deposit:\t", i)
		mu.Lock()
		balance += i
		deHistory = append(deHistory, i) // Log deposit history
		mu.Unlock()
	}
}

func subtract(sub chan int) {
	defer wg.Done()
	for i := range sub {
		log.Println("Withdraw:\t", i)
		mu.Lock()
		if balance < i {
			log.Println("Cannot subtract")
			continue
		}
		balance -= i
		subHistory = append(subHistory, i) // Log subtraction history
		mu.Unlock()
	}
}

func main() {
	var in = make(chan int)
	var sub = make(chan int)
	c := cron.New()

	wg.Add(4)
	go add(in)
	go subtract(sub)

	// Sending deposits every 1 second
	var senderIn sync.WaitGroup
	var sendSub sync.WaitGroup
	senderIn.Add(1)
	go func() {
		defer senderIn.Done()
		_, err := c.AddFunc("@every 1s", func() {
			select {
			case in <- 1000:
			default:
				log.Println("Deposit channel is closed, skipping")
			}
		})
		if err != nil {
			log.Fatal("Error adding cron function for deposit:", err)
			return
		}
	}()

	// Sending withdrawals every 2 seconds
	sendSub.Add(1)
	go func() {
		defer sendSub.Done()
		_, err := c.AddFunc("@every 2s", func() {
			select {
			case sub <- 100:
			default:
				log.Println("Withdraw channel is closed, skipping")
			}
		})
		if err != nil {
			log.Fatal("Error adding cron function for withdrawal:", err)
			return
		}
	}()

	// Print balance and history every 10 seconds
	var printing sync.WaitGroup
	printing.Add(1)
	go func() {
		defer printing.Done()
		_, err := c.AddFunc("@every 10s", func() {
			mu.Lock()
			fmt.Println("Balance:\t", balance)
			for idx, deposit := range deHistory {
				fmt.Println("Deposit:\t", idx, ":", deposit)
			}
			for idx, sub := range subHistory {
				fmt.Println("Withdraw:\t", idx, ":", sub)
			}
			mu.Unlock()
		})
		if err != nil {
			log.Fatal("Error adding cron function for printing:", err)
			return
		}
	}()

	// Start the cron scheduler
	c.Start()

	// Handle OS interrupt signal (Ctrl+C) to stop the program
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	<-osChan

	// Wait for all goroutines to finish
	printing.Wait()
	senderIn.Wait()
	sendSub.Wait()

	// Close channels and stop cron scheduler
	close(in)
	close(sub)

	// Wait for all goroutines to finish before stopping the cron job
	wg.Wait()
	c.Stop()

	fmt.Println("Final balance:", balance)
}

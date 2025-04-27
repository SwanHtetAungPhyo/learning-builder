package main

//
//import (
//	"fmt"
//	"log"
//	"sync"
//)
//
//var (
//	balance int
//	mu      sync.Mutex
//	wg      sync.WaitGroup
//	cond    sync.Cond // To control the interleaving of operations
//)
//
//func add(in chan int) {
//	defer wg.Done()
//	for i := range in {
//		mu.Lock() // Lock the mutex before modifying balance
//		balance += i
//		log.Println("Add:", i, "New Balance:", balance)
//		mu.Unlock()   // Unlock the mutex after modification
//		cond.Signal() // Notify that the add operation is done
//		cond.Wait()   // Wait for the signal to subtract
//	}
//}
//
//func subtract(sub chan int) {
//	defer wg.Done()
//	for i := range sub {
//		mu.Lock() // Lock the mutex before modifying balance
//		balance -= i
//		log.Println("Subtract:", i, "New Balance:", balance)
//		mu.Unlock()   // Unlock the mutex after modification
//		cond.Signal() // Notify that the subtract operation is done
//		cond.Wait()   // Wait for the signal to add again
//	}
//}
//
//func main() {
//	in := make(chan int)
//	sub := make(chan int)
//
//	// Initialize the Cond with a Mutex to use it safely
//	cond.L = &mu
//
//	// Start add and subtract goroutines
//	wg.Add(2) // We need to wait for both goroutines to finish
//	go add(in)
//	go subtract(sub)
//
//	// Goroutines to send data to channels
//	var sender sync.WaitGroup
//	sender.Add(2)
//
//	go func() {
//		defer sender.Done()
//		// Sending data for additions
//		for i := 0; i < 5; i++ {
//			in <- 1000
//			cond.Wait() // Wait until subtract can happen
//		}
//		close(in) // Close the add channel
//	}()
//
//	go func() {
//		defer sender.Done()
//		// Sending data for subtractions
//		for i := 0; i < 5; i++ {
//			sub <- 500
//			cond.Wait() // Wait until add can happen again
//		}
//		close(sub) // Close the subtract channel
//	}()
//
//	// Start the sender goroutines, wait until they're done
//	sender.Wait()
//
//	// Send a signal to start adding after all subtractions are done
//	cond.Signal()
//
//	// Wait for the add and subtract goroutines to finish
//	wg.Wait()
//
//	fmt.Println("Final balance:", balance)
//}

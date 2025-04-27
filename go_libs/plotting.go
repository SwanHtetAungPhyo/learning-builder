package main

import (
	"fmt"
	"github.com/guptarohit/asciigraph"
	"math/rand"
	"time"
)

func main() {
	var digitChan = make(chan int)
	var data []float64
	go func() {
		for i := range digitChan {
			num := float64(i)
			data = append(data, num)
			if len(data) > 100 {
				data = data[1:]
			}

			fmt.Print("\033[H\033[2J")

			fmt.Printf("%s\n", asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(100)))
		}
	}()

	go func() {
		for _ = range 10000 {
			digitChan <- rand.Intn(20)
			time.Sleep(900 * time.Millisecond)
		}
		close(digitChan)
	}()
	select {}
}

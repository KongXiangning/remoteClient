package main

import (
	"fmt"
	"time"
)

func main() {
	unbufChan := make(chan int)
	sign := make(chan byte, 2)
	small := make(chan int, 1)
	lager := make(chan int, 1)

	go func() {
		for i := 0; i < 10; i++ {
			if i < 5 {
				small <- i
			} else {
				lager <- i
			}
			select {
			case unbufChan <- small:
			case unbufChan <- lager:

			}
			fmt.Printf("The %d select is selected\n", i)
			time.Sleep(time.Second)
		}
		close(unbufChan)
		fmt.Println("The channel is closed.")
		sign <- 0
	}()

	go func() {
	loop:
		for {
			select {
			case e, ok := <-unbufChan:
				if !ok {
					fmt.Println("Closed channel.")
					break loop
				}
				fmt.Printf("e: %d\n", e)
				time.Sleep(2 * time.Second)
			}
		}
		sign <- 1
	}()
	<-sign
	<-sign
}

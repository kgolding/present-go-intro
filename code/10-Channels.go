package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := generator()

	fmt.Println("Waiting for first item")
	firstItem := <-intChan
	fmt.Println("First item is", firstItem)

	for item := range intChan {
		fmt.Println("Received item", item)
	}
	fmt.Println("intChan has closed")
}

func generator() chan int {
	myChan := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(time.Second)
			myChan <- i
			if i == 5 {
				close(myChan)
				return
			}
			i++
		}
	}()
	return myChan
}

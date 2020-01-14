package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := generator()

	timeOut := time.After(time.Second * 6)

	for {
		select {
		case item, ok := <-intChan:
			if !ok {
				fmt.Println("generator finished")
				return
			}
			fmt.Println("Received item", item)

		case <-timeOut:
			fmt.Println("timed out")
			return
		}
	}
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

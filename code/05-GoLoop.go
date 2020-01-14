package main

import (
	"fmt"
	"time"
)

func main() {
	a := []string{"Zero", "One", "Two", "Three", "Four", "Five"}

	for i, s := range a {
		go output(i, s)
	}
	time.Sleep(time.Millisecond)
}

func output(i int, s string) {
	fmt.Println("Index", i, "is", s)
}

package main

import (
	"fmt"
	"sync"
)

func main() {
	a := []string{"Zero", "One", "Two", "Three", "Four", "Five"}
	var wg sync.WaitGroup
	for i, s := range a {
		wg.Add(1)
		go output(i, s, &wg)
	}
	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("All done")
}

func output(i int, s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Index", i, "is", s)
}

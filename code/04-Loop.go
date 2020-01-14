package main

import (
	"fmt"
)

func main() {
	a := []string{"Zero", "One", "Two", "Three", "Four", "Five"}

	for i, s := range a {
		fmt.Println("Index", i, "is", s)
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.88.253:5678")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	log.Println("Connected!")

	go func() {
		b := make([]byte, 255)
		for {
			n, err := conn.Read(b)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Print(string(b[:n]))
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		conn.Write([]byte(in))
	}
}

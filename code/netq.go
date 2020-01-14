package main

import (
	"log"
	"net"
	"time"
)

func main() {
	localAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:5678")
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	go sender(conn)

	log.Println("Listening")
	for {
		buffer := make([]byte, 256)
		numberBytes, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error receiving", err)
		}
		log.Println("Received from", remoteAddr, ":", string(buffer[:numberBytes]))
	}
}

func sender(conn *net.UDPConn) {
	broadcastAddr, _ := net.ResolveUDPAddr("udp", "192.168.5.255:5678")
	for {
		time.Sleep(time.Second)
		question := "5 + 5"
		_, err := conn.WriteToUDP([]byte(question), broadcastAddr)
		if err != nil {
			log.Fatal("Error sending", err)
		}
		log.Println("Sent:", question)
	}
}

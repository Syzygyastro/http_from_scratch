package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Error resolving the server address %s\n", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Error dialing the UDP address %s\n", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("=========================================")
	fmt.Println("Connection established to localhost:42069!")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading line %s/n", err.Error())
			os.Exit(1)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Line sent: %s", line)
	}

}

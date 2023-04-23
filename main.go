package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	var port string
	// Check in port
	if len(os.Args[1:]) == 0 {
		port = ":8989"
	} else if len(os.Args[1:]) == 1 {
		// checkong our massiv "s"
		port = checkPort(os.Args[1:][0])
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	l, err := net.Listen("tcp", port)

	if err != nil {
		log.Printf("Could not start server: %s. \n", err.Error())
		return
	}

	log.Println("Server is running on port " + port)
	var chat = NewChat()
	defer func() {
		os.Remove("history.txt")
		l.Close()
	}()

	chat.history, err = os.Create("history.txt")
	if err != nil {
		log.Printf("Could read a file: %s. \n", err.Error())
		return
	}
	go broadcast(chat)

	for {
		conn, err := l.Accept() 
		if err != nil {
			log.Printf("Could not accept connection: %s. \n", err.Error())
			return
		}
		go handleConnection(conn, chat)
	}
}
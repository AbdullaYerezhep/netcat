package main

import (
	"net"
	"os"
	"sync"
)

type chat struct {
	sync.Mutex
	users map[net.Conn]string
	connections int 
	msg chan message
	user chan user
	history *os.File
}

type message struct {
	username string
	time string
	text string
}

type user struct {
	conn     net.Conn
	username string
}

func NewChat() *chat {
	return &chat{
		users: make(map[net.Conn]string),
		connections: 0,
		msg: make(chan message),
		user: make(chan user),
		history: &os.File{},
	}
}

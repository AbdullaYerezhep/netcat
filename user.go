package main

import (
	"bufio"
	"fmt"
	"net"
)


func handleConnection(c net.Conn, chat *chat ) {
	if chat.isFull(c) {
		return
	}
	chat.connections++
	
	chat.welcome(c)

	username := chat.username(c) 
	
	chat.Mutex.Lock()
	chat.users[c] = username
	chat.Mutex.Unlock()

	newUser := user{conn: c, username: username}
	

	chat.user <- newUser

	chat.readHistory(c)

	input := bufio.NewScanner(c)
	for	input.Scan() {
		currentTime := generateTime()
		if !checkSymbol(input.Text()) {
			_, err := fmt.Fprintln(c, "Error in text output\noutput text again")
			chat.check(err, c)
		}
		text := spacedeletion(input.Text())
		msg := message{username: username, time: currentTime, text: text + "\n"}
		chat.history.WriteString(msg.time + "[" + msg.username+ "]:" + msg.text)
		chat.msg <- msg
	}

	defer func(){
		chat.connections--
		c.Close()
		delete(chat.users, c)
		msg := message{username: username, text:"", time:""}
		chat.history.WriteString(msg.username + " has left out chat...\n")
		chat.msg <- msg
	}()
}

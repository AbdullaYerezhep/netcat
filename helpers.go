package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

func generateTime() string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s]", currentTime)
}

func spacedeletion(s string) string {
	var result string
	var count int
	for _, v := range s {
		if v != ' ' {
			if len(result) == 0 {
				count = 0
			}
			if count > 0 {
				result += " "
				count = 0
			}
			result += string(v)
		} else {
			count++
			continue
		}
	}
	return result
}


func checkPort(s string) string {
	var result string
	// the port check for the length according to the diphtholt it should be 4
	if len(s) != 4 {
		fmt.Println("port number length should be equal to 4, try again")
		os.Exit(0)
	}
	for _, v := range s {
		// numbers check
		if v < '0' || v > '9' {
			fmt.Println("There must be numbers")
			os.Exit(0)
			break
		} else {
			result += string(v)
		}
	}
	return ":" + result
}



func checkSymbol(s string) bool {
	s = strings.TrimSuffix(s, "\n")
	rxmsg := regexp.MustCompile("^[\u0400-\u04FF\u0020-\u007F]+$")
	if !rxmsg.MatchString(s) {
		return false
	}
	return true
}



func (ch *chat) welcome(c net.Conn) {
	welcomeTxt, err := ioutil.ReadFile("welcome.txt")
	if err != nil {
		log.Println("Could not read a file!")
	}
	
	fmt.Fprintf(c, string(welcomeTxt) + "\n")
	
}



func (ch *chat) isFull(c net.Conn) bool {
	if ch.connections >= 10 {
		fmt.Fprintf(c, "Chat is full, please try again later.\n")
		c.Close()
		return true
	}
	return false
}



func (ch *chat) username(c net.Conn) string {
	var username string
	scanner := bufio.NewScanner(c)
	for {
		_, err := fmt.Fprint(c, "[ENTER YOUR NAME]: ")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		if scanner.Scan() {
			if !checkSymbol(scanner.Text()) {
				_, err := fmt.Fprint(c, "Forbidden symbols were used\ntry again\n")
				if err != nil {
					return ""
				}
				continue
			}
			if spacedeletion(scanner.Text()) == "" {
				_, err := fmt.Fprintln(c, "No empty username is allowed")
				if err != nil {
					return ""
				}
				continue
			}
			if !ch.repeatName(spacedeletion(scanner.Text())) {
				_, err := fmt.Fprintln(c, "Such name already exists!")
				if err != nil {
					return ""
				}
				continue
			}
			username = spacedeletion(scanner.Text())
			break
		}
	}
	return username
}



func (ch *chat) repeatName(s string) bool {
	for _, name := range ch.users {
		if name == s {
			return false
		}
	}
	return true
}



func (ch *chat) readHistory(c net.Conn) {
	history, err := os.ReadFile("history.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(history) != 0 {
		_, err := fmt.Fprint(c, string(history))
		ch.check(err, c)
	}
}


func (ch *chat) checkBrodcaster(err error, c net.Conn) {
	log.Println("Could not output to the user")
	delete(ch.users, c)
	left := message{username: ch.users[c], text:"", time:""}
	ch.msg <- left
}

func (ch *chat) check(err error, c net.Conn) {
	if err != nil {
		log.Println("Could not output to the user")
		delete(ch.users, c)
		left := message{username: ch.users[c], text:"", time:""}
		ch.msg <- left
		return
	}
}

package main

import "fmt"

func broadcast(ch *chat) {
	for {
		select {
			case user := <- ch.user:
				ch.Lock()
				msg := user.username + " has joined our chat...\n"
				ch.history.WriteString(msg)
				for k := range ch.users {
					if k != user.conn {
						_, err := fmt.Fprintf(k, msg)
						if err != nil {
							ch.checkBrodcaster(err, k)
							continue
						}
					}
				}
				ch.Unlock()
			case msg := <- ch.msg:
				ch.Lock()
				for k := range ch.users {
						if msg.time == "" {
							msg := msg.username + " has left out chat...\n"
							
							_, err := fmt.Fprintf(k, msg)
							if err != nil {
								ch.checkBrodcaster(err, k)
								continue
							}
						}else{
							if len(msg.text) != 0 {
								msg :=  msg.time + "[" + msg.username+ "]:" + msg.text
	
								_, err := fmt.Fprintf(k, msg)
								if err != nil {
									ch.checkBrodcaster(err, k)
									continue
								}
							}
						}
				}
				ch.Unlock()
		}
	}
}
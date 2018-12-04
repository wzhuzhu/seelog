package seelog

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

//  websocket客户端
type client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

// 客户端管理
type clientManager struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

var manager = clientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
}

func (manager *clientManager) start() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[seelog] error:%+v", err)
		}
	}()

	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				conn.socket.Close()
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (c *client) write() {

	for msg := range c.send {
		_, err := c.socket.Write(msg)
		if err != nil {
			fmt.Println("write msg failed. ", err)
			break
		}
	}
	log.Println("web socket closed")
}

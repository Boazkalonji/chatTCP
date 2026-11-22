package main

import (
	"net"
	"sync"
)

type GlobalHub struct {
	clients    map[net.Conn]*string
	cadenas    sync.Mutex
	braoadcast chan string
}

func (hub *GlobalHub) DiffusionMessage() {

	message := <-hub.braoadcast

	hub.cadenas.Lock()

	for conn, _ := range hub.clients {

		_, err := conn.Write([]byte(message))

		if err != nil {

			conn.Close()
			delete(hub.clients, conn)

		}

	}
	hub.cadenas.Unlock()

}

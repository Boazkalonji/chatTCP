package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	db := InitDB()
	defer db.Close()

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		fmt.Printf("une erreur est survenu au moment de l'ecoute %v", err)
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Println("🟢 SERVEUR DE CHAT EN LIGNE SUR LE PORT 8080...")

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("une erreur survenu au moment de la connexion %v", err)
		}

		fmt.Printf("📥 Nouveau client connecté : %s\n", conn.RemoteAddr().String())

		conn.SetWriteDeadline(time.Now().Add(300 * time.Second))

		go gererClient(conn, db)

	}

}

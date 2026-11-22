package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
)

var HubPrincipal = GlobalHub{
	clients:    make(map[net.Conn]*string),
	braoadcast: make(chan string, 100),
}

func main() {

	db := InitDB()
	defer db.Close()

	certificat, err := tls.LoadX509KeyPair("cert.pem", "key.pem")

	if err != nil {
		fmt.Printf("une erreur est survenu au moment du telechargement de clé  %v", err)
	}

	MaConfig := &tls.Config{

		Certificates: []tls.Certificate{certificat},
		MinVersion:   tls.VersionTLS12,
	}

	listener, err := tls.Listen("tcp", "127.0.0.1:8080", MaConfig)

	if err != nil {
		fmt.Printf("une erreur est survenu au moment de l'ecoute %v", err)
		os.Exit(1)
	}

	defer listener.Close()

	go HubPrincipal.DiffusionMessage()

	fmt.Println("🟢 SERVEUR DE CHAT EN LIGNE SUR LE PORT 8080...")

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("une erreur survenu au moment de la connexion %v", err)
		}

		fmt.Printf("📥 Nouveau client connecté : %s\n", conn.RemoteAddr().String())

		go gererClient(conn, db)

	}

}

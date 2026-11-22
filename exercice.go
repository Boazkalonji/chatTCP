package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Requete struct {
	Action string `json:"action"`
	User   string `json:"user"`
}

type Reponse struct {
	Statut  string `json:"status"`
	Message string `json:"message"`
}

func gererClientT(conn net.Conn) {

	defer func() {

		_, err := conn.Write([]byte(string("LE SERVEUR EST FERMER ")))

		if err != nil {
			fmt.Printf("UNE ERREUR AU MOMENT DE FERMER LE SERVEUR")
			conn.Close()
		}

		conn.Close()
	}()

	for {

		//conn.SetReadDeadline(time.Now().Add(15 * time.Second))

		cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		canalResultat := make(chan bool)

		go func() {

			var requete Requete
			_ = json.NewDecoder(conn).Decode(&requete)
			fmt.Printf("LA REQUET DU CLIENT EST : ACTION %v - DONNEE %v", requete.Action, requete.Donnee)
			canalResultat <- true

			switch requete.Action {

			case "ping":

				json.NewEncoder(conn).Encode(Reponse{
					Statut:  "OK",
					Message: "pong",
				})

			case "echo":

				json.NewEncoder(conn).Encode(Reponse{
					Statut:  "OK",
					Message: requete.Donnee,
				})

			case "quit":

				json.NewEncoder(conn).Encode(Reponse{
					Statut:  "BYEBYE",
					Message: "AU REVOIR",
				})

			default:
				json.NewEncoder(conn).Encode(Reponse{
					Statut:  "ERREUR",
					Message: "	inconue",
				})

			}

		}()

		select {

		case <-cxt.Done():
			fmt.Println("⏰ [CONTEXTE] Le decodage  a pris trop de temps.")
			conn.Write([]byte("Erreur : Temps de traitement dépassé.\n"))
			break

		case status := <-canalResultat:

			if status != true {

				fmt.Printf("✅ REQUETE DECODE AVEC SUCCES \n")
				conn.Write([]byte("✅ Message enregistré !\n"))

			}

		}

	}

}

func mainT() {

	listener, err := net.Listen("tcp", "127.0.0.1:9090")

	if err != nil {
		fmt.Printf("une erreur est survenu au moment de l'ecoute %v", err)
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Println("🟢 SERVEUR ET EN LIGNE SUR LE PORT 9090...")

	StopChan := make(chan os.Signal, 1)
	signal.Notify(StopChan, os.Interrupt, syscall.SIGTERM)

	go func(l net.Listener) {

		for {

			conn, err := l.Accept()

			if err != nil {
				fmt.Printf("une erreur survenu au moment de la connexion %v", err)
			}

			fmt.Printf("📥 Nouveau client connecté : %s\n", conn.RemoteAddr().String())

			go gererClient(conn)

		}

	}(listener)

	<-StopChan
	fmt.Printf("LE SERVEUR EST COUPER")
	listener.Close()

}

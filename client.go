package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"
)

func gererClient(conn net.Conn, db *sql.DB) {

	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(200 * time.Second))

	pseudoBuf := make([]byte, 1024)
	nPseudo, err := conn.Read(pseudoBuf)

	if err != nil {
		fmt.Printf("UNE ERREUR AU MOMENT DE LA LECTURE DE LA DONNE %v", err)
		return
	}
	pseudo := strings.TrimSpace(string(pseudoBuf[:nPseudo]))

	for {

		messageBuf := make([]byte, 2048)
		nMessage, err := conn.Read(messageBuf)

		if err != nil {
			fmt.Printf("UNE ERREUR AU MOMENT DE LA ECRITURE DU MESSAGE  %v", err)

			if messageBuf != nil {

				break

			}
		}

		message := strings.TrimSpace(string(messageBuf[:nMessage]))

		

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		canalResultat := make(chan bool, 2)

		go func() {

			dateDuJour := time.Now().Format("2006-01-02 15:04:05")

			_, err := db.Exec("INSERT INTO messages (auteur, contenu, date) VALUES (?, ?, ?)", pseudo, message, dateDuJour)
			if err != nil {
				fmt.Println("Erreur d'insertion BDD :", err)
				canalResultat <- false
				return
			}
			canalResultat <- true

		}()

		select {

		case <-ctx.Done():

			fmt.Println("⏰ [CONTEXTE] L'enregistrement a pris trop de temps.")
			conn.Write([]byte("Erreur : Temps de traitement dépassé.\n"))
			return

		case success := <-canalResultat:
			if !success {
				conn.Write([]byte("❌ Erreur lors de l'enregistrement.\n"))
				continue
			}
			fmt.Printf("✅ Message de %s enregistré en BDD.\n", pseudo)
			conn.Write([]byte("✅ Message enregistré !\n"))
		}

	}

}

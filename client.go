package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Action     string "`json:action`"
	User       string "`json:user`"
	MotDePasse string "`json:motDePasse`"
}

func Authentification(conn net.Conn, donneRecu Login, db *sql.DB) (string, bool) {

	conn.SetReadDeadline(time.Now().Add(15 * time.Second))

	if donneRecu.Action != "LOGIN" {

		conn.Write([]byte("VEUILLEZ CHOISIR LA BONNE OPTION OU CREER VOUS UN COMPT"))
		return "", false

	}

	var mdpHash string
	err := db.QueryRow("SELECT mdp FROM users WHERE user = ?", donneRecu.User).Scan(&mdpHash)

	if err != nil {
		fmt.Printf("une erreur au moment de la recuperation du mdp dans la base %v", err)
		return "", false

	}

	errHash := bcrypt.CompareHashAndPassword([]byte(mdpHash), []byte(donneRecu.MotDePasse))

	if errHash != nil {
		fmt.Printf("Une erreur au moment du hashage %v", errHash)
		return "", false

	}

	_, resultat := CreeSsion(donneRecu.User, conn)

	if resultat == false {
		fmt.Printf("ERREUR AU MOMENT DE L'AJOUT DE L'USER DANS SESSIONACTIVE %v", resultat)
	}

	conn.SetReadDeadline(time.Time{})
	conn.Write([]byte("CONNEXION REUSSIE"))

	return donneRecu.User, true

}

func gererClient(conn net.Conn, db *sql.DB) {

	defer conn.Close()

	var donneRecu Login

	SaisieClient := json.NewDecoder(conn)
	err := SaisieClient.Decode(&donneRecu)

	if err != nil {
		fmt.Printf("une erreur au moment de la recuperation du json %v", err)
	}

	if donneRecu.Action == "REGISTER" {

		_, err := bcrypt.GenerateFromPassword([]byte(donneRecu.MotDePasse), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("une erreur au moment du hashage du mdp %v", err)
		}

		dateDuJour := time.Now().Format("2006-01-02 15:04:05")

		_, err = db.Exec("INSERT INTO users (user, mdp, date) VALUES (?, ?, ?)", donneRecu.User, donneRecu.MotDePasse, dateDuJour)
		if err != nil {
			fmt.Println("Erreur d'insertion BDD :", err)
			return

		}

	} else if donneRecu.Action == "AUTH_TOKEN" {

		var session Session

		TokenEnvoyer := json.NewDecoder(conn)
		err = TokenEnvoyer.Decode(&session)

		resultat := VerfiierToken(session.Token, conn)

		if resultat == true {

			conn.Write([]byte("CONNEXION PAR TOKEN RECU"))

		}

	}

	for {

		_, resultat := Authentification(conn, donneRecu, db)

		if resultat == false {
			conn.Write([]byte("ERREUR DE CONNEXION "))
			break

		}

	}

	pseudo := strings.TrimSpace(string(donneRecu.User))

	go func(c net.Conn) {

		HubPrincipal.cadenas.Lock()

		HubPrincipal.clients[c] = &pseudo
		HubPrincipal.cadenas.Unlock()

	}(conn)

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

			HubPrincipal.braoadcast <- fmt.Sprintf("Voici le messesage du client ", pseudo, message)

			canalResultat <- true

		}()

		go func() {

			for {

				msg := make([]byte, 2048)

				_, err := conn.Read(msg)

				if err != nil {
					fmt.Print("UNE EREUR SE PRODUIT %v", err)
				}

				os.Stdout.Write([]byte(msg))

			}

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

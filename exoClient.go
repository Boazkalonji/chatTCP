package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Requete struct {
	Action string `json:"action"`
	Donnee string `json:"donnee"`
}

type Reponse struct {
	Statut  string `json:"status"`
	Message string `json:"message"`
}

func Tmain() {
	call, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("❌ Erreur de connexion au serveur : %v\n", err)
		os.Exit(1)
	}
	defer call.Close()
	fmt.Println("🟢 CONNECTÉ AU SERVEUR (127.0.0.1:9090)")

	// 1. Goroutine d'écoute en arrière-plan (la boucle for est DEDANS)
	go func(c net.Conn) {
		decoder := json.NewDecoder(c)
		for {
			var response Reponse
			err := decoder.Decode(&response)
			if err != nil {
				fmt.Println("\n🔴 Connexion fermée par le serveur (timeout ou arrêt).")
				os.Exit(0)
			}
			fmt.Printf("\n📩 [SERVEUR] Statut : %s | Message : %s\n", response.Statut, response.Message)
			fmt.Print("👉 Action (ping / echo / quit) : ")
		}
	}(call)

	// 2. Boucle principale d'envoi depuis le clavier
	lecteurClavier := bufio.NewReader(os.Stdin)
	encoder := json.NewEncoder(call)

	for {
		fmt.Print("👉 Action (ping / echo / quit) : ")
		actionSaisie, _ := lecteurClavier.ReadString('\n')
		actionPropre := strings.TrimSpace(actionSaisie) // Enlève le \r\n de Windows !

		var donneePropre string
		if actionPropre == "echo" {
			fmt.Print("✏️  Donnée à envoyer : ")
			donneeSaisie, _ := lecteurClavier.ReadString('\n')
			donneePropre = strings.TrimSpace(donneeSaisie)
		}

		// Envoi de la structure Requete au serveur
		err := encoder.Encode(Requete{
			Action: actionPropre,
			Donnee: donneePropre,
		})
		if err != nil {
			fmt.Printf("❌ Erreur d'envoi : %v\n", err)
			break
		}

		if actionPropre == "quit" {
			break
		}
	}
}

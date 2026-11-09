package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("UNE ERREUR AU MOMENT DE LA RECEPTION DE L'APPEL, %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	pseudo := make([]byte, 1024)

	fmt.Printf("VEUILLEZ SAISIR VOTRE PSEUDO SVP !!! \n")
	n, _ := os.Stdin.Read(pseudo)
	conn.Write([]byte(string(pseudo[:n])))

	for {
		fmt.Printf("VEUILLEZ SAISIR VOTRE MESSAGE SVP !!! \n")
		texte := make([]byte, 1024)
		n, _ = os.Stdin.Read(texte)

		conn.Write([]byte(string(texte[:n])))
		


		response := make([]byte, 2048)

		n, _ = conn.Read(response)

		fmt.Printf("LA REPONSE DU SERVEUR EST : %s", string(response[:n]))
	}
}

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Login struct {
	Action     string `json:"action"`
	User       string `json:"user"`
	MotDePasse string `json:"motDePasse"`
}

type Session struct {
	User       string    `json:"user"`
	Expiration time.Time `json:"expiration"`
	Token      string    `json:"token"`
}

func main() {

	ConigClient := &tls.Config{

		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:8080", ConigClient)
	if err != nil {
		fmt.Printf("UNE ERREUR AU MOMENT DE LA RECEPTION DE L'APPEL, %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("VEUILLEZ CHOISIR L'OPTION QUI VOUS CONVIENT  \n 1. LOGIN \n 2. REGISTER \n 3. AUTH_TOKEN \n")
	LecteurClavier := bufio.NewReader(os.Stdin)
	action, _ := LecteurClavier.ReadString('\n')
	action = strings.TrimSpace(action)

	switch action {
	case "1": //LOGIN

		fmt.Printf("VEUILLEZ SAISIR VOTRE NOM(OU PSEUDO) \n")
		user, _ := LecteurClavier.ReadString('\n')
		userFinal := strings.TrimSpace(user)

		fmt.Printf("VEUILLEZ SAISIR VOTRE PASSWORD  \n")
		mdp, _ := LecteurClavier.ReadString('\n')
		mdpFinal := strings.TrimSpace(mdp)

		Jsonlogin := json.NewEncoder(conn)
		Jsonlogin.Encode(Login{

			Action:     "LOGIN",
			User:       userFinal,
			MotDePasse: mdpFinal,
		})

	case "2": //REGISTER

		fmt.Printf("VEUILLEZ SAISIR VOTRE NOM(OU PSEUDO) \n")
		user, _ := LecteurClavier.ReadString('\n')
		userFinal := strings.TrimSpace(user)

		fmt.Printf("VEUILLEZ SAISIR VOTRE PASSWORD  \n")
		mdp, _ := LecteurClavier.ReadString('\n')
		mdpFinal := strings.TrimSpace(mdp)

		Jsonlogin := json.NewEncoder(conn)
		Jsonlogin.Encode(Login{

			Action:     "REGISTER",
			User:       userFinal,
			MotDePasse: mdpFinal,
		})

	case "3": //AUTH_TOKEN

		fmt.Printf("VEUILLEZ SAISIR VOTRE NOM(OU PSEUDO) \n")
		user, _ := LecteurClavier.ReadString('\n')
		userFinal := strings.TrimSpace(user)

		Jsonlogin := json.NewEncoder(conn)

		//var TokenGenere Session
		//recuToken := json.NewDecoder(conn)
		//recuToken.Decode(&TokenGenere)
		//TokenFinalPropre := TokenGenere.Token

		Jsonlogin.Encode(Session{
			User: userFinal,
			//Token: TokenFinalPropre,
		})
	}

	for {
		fmt.Printf("VEUILLEZ SAISIR VOTRE MESSAGE SVP !!! \n")
		texte := make([]byte, 1024)
		n, _ := os.Stdin.Read(texte)

		conn.Write([]byte(string(texte[:n])))

		response := make([]byte, 2048)

		n, _ = conn.Read(response)

		fmt.Printf("LA REPONSE DU SERVEUR EST : %s", string(response[:n]))
	}
}

package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

var SessionActives = make(map[string]Session)
var MuSession sync.Mutex

func CreeSsion(user string, conn net.Conn) (string, bool) {

	octets := make([]byte, 32)
	_, err := rand.Read(octets)

	if err != nil {
		fmt.Printf("Une erreur au moment de la generation du Token %v", err)
		return "", false
	}
	token := hex.EncodeToString(octets)
	fmt.Printf(token)

	MuSession.Lock()
	SessionActives[token] = Session{
		User:       user,
		Expiration: time.Now().Add(2 * time.Hour),
	}
	MuSession.Unlock()

	TokenENVOIE := json.NewEncoder(conn)

	TokenENVOIE.Encode(Session{
		User:       user,
		Expiration: SessionActives[token].Expiration,
		Token:      token,
	})

	return token, true

}

func VerfiierToken(tokenRecu string, conn net.Conn) bool {

	MuSession.Lock()
	session, existe := SessionActives[tokenRecu]
	MuSession.Unlock()

	if existe != false || time.Now().After(session.Expiration) {

		conn.Write([]byte("VOTRE SESSIONS A EXPIRER , CONNECTER VOUS "))
		return false

	}

	comparaison := subtle.ConstantTimeCompare([]byte(tokenRecu), []byte(SessionActives[tokenRecu].Token)) == 1
	fmt.Printf("", comparaison)
	return true

}

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

type Session struct {
	User       string    `json:"user"`
	Expiration time.Time `json:"expiration"`
	Token      string    `json:"token"`
}

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

func VerfiierToken(user string, tokenRecu string) bool {

	MuSession.Lock()
	defer MuSession.Unlock()

	session, existe := SessionActives[tokenRecu]
	if !existe {
		return false
	}

	comparaison := subtle.ConstantTimeCompare([]byte(session.Token), []byte(tokenRecu)) == 1
	fmt.Printf("", comparaison)
	return true

}

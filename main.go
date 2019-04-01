package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"time"
	"github.com/gofrs/uuid"
	"errors"
)

type user struct {
	uuid       uuid.UUID
	name       []byte
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

var userList []user

/*
	reqLogger logs the attributes
	1. Time
	2. HTTP verb
	3. Requested URI
	4. Remote Address
	of an incoming request
*/
func reqLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request Time %s ", time.Now())
		log.Printf("Method %s ", r.Method)
		log.Printf("Request URI %s ", r.RequestURI)
		log.Printf("Remote address %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

/*
	Crypto Utilities
	___________________________________________________________________________
*/

/*
	generateKeyPair() generate a private and public key pair based on ECDSA
*/
func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, err
}

/*
	signTx() signs the payload with the provided private key
*/
func signTx(priv *ecdsa.PrivateKey, payload []byte) (*big.Int, *big.Int, error) {
	r, s, err := ecdsa.Sign(rand.Reader, priv, payload)
	if err != nil {
		return nil, nil, err
	}
	return r, s, err
}

/*
	verifyTx() verify the payload against the provided public key
*/
func verifyTx(pub *ecdsa.PublicKey, payload []byte, r, s *big.Int) bool {
	flag := ecdsa.Verify(pub, payload, r, s)
	return flag
}

/*
	User Utilities
	___________________________________________________________________________
*/

/*
	createUser() creates a user and append it to userList slice
*/
func createUser(name []byte) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	privKey, pubKey, err := generateKeyPair()
	if err != nil {
		return err
	}
	u := user{uuid: uuid, name: name, privateKey: privKey, publicKey: pubKey}
	userList = append(userList, u)
	return nil
}

/*
	getPrivateKey() returns private key with provided uuid
*/
func getPrivateKey(uuid uuid.UUID) (*ecdsa.PrivateKey, error) {
	for _, u := range userList {
		if u.uuid == uuid {
			return u.privateKey, nil
		}
	}
	return nil, errors.New("user not found")
}

/*
	getPublicKey() returns public key with provided uuid
*/
func getPublicKey(uuid uuid.UUID) (*ecdsa.PublicKey, error) {
	for _, u := range userList {
		if u.uuid == uuid {
			return u.publicKey, nil
		}
	}
	return nil, errors.New("user not found")
}

/*
	getUserName() returns name with provided uuid
*/
func getUserName(uuid uuid.UUID) ([]byte, error) {
	for _, u := range userList {
		if u.uuid == uuid {
			return u.name, nil
		}
	}
	return nil, errors.New("user not found")
}

/*
	indexHandler serves '/' endpoint
*/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

/*
	dashboardHandler serves '/dashboard' endpoint
*/
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/dashboard.html")
}

func main() {
	http.HandleFunc("/", reqLogger(indexHandler))
	http.HandleFunc("/dashboard", reqLogger(dashboardHandler))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./assets/css"))))
	log.Printf("App running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
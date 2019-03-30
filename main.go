package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"time"
)

func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, err
}

func signTx(priv *ecdsa.PrivateKey, payload []byte) (*big.Int, *big.Int, error) {
	r, s, err := ecdsa.Sign(rand.Reader, priv, payload)
	if err != nil {
		return nil, nil, err
	}
	return r, s, err
}

func verifyTx(pub *ecdsa.PublicKey, payload []byte, r, s *big.Int) bool {
	flag := ecdsa.Verify(pub, payload, r, s)
	return flag
}

func reqLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request Time %s ", time.Now())
		log.Printf("Method %s ", r.Method)
		log.Printf("Request URI %s ", r.RequestURI)
		log.Printf("Remote address %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/dashboard.html")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	priv, pub, err := generateKeyPair()
	if err != nil {
		log.Print(err)
	}
	log.Print("private key ", priv)
	log.Print("public key ", pub)

	payload := []byte("abcdefghijklmnopqrstuvwxyz")
	p, q, err := signTx(priv, payload)
	if err != nil {
		log.Printf("cannot sign payload")
	}
	res := verifyTx(pub, payload, p, q)
	log.Print("res ", res)

}

func main() {
	http.HandleFunc("/", reqLogger(indexHandler))
	http.HandleFunc("/dashboard", reqLogger(dashboardHandler))
	http.HandleFunc("/test", reqLogger(testHandler))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./assets/css"))))
	log.Printf("App running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

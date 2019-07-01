package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

type user struct {
	UUID       uuid.UUID `json:"uuid"`
	Name       string    `json:"name"`
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

var userList []user

type transaction struct {
	timeStamp int64
	txMessage []byte
	prevHash  []byte
	currHash  []byte
}

type block struct {
	Tx []*transaction
}

var blk block

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
func createUser(name string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	privKey, pubKey, err := generateKeyPair()
	if err != nil {
		return err
	}
	u := user{UUID: uuid, Name: name, privateKey: privKey, publicKey: pubKey}

	payload, _ := json.Marshal(u)
	log.Printf(string(payload))

	userList = append(userList, u)
	return nil
}

/*
	getPrivateKey() returns private key with provided uuid
*/
func getPrivateKey(uuid uuid.UUID) (*ecdsa.PrivateKey, error) {
	for _, u := range userList {
		if u.UUID == uuid {
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
		if u.UUID == uuid {
			return u.publicKey, nil
		}
	}
	return nil, errors.New("user not found")
}

/*
	getUserName() returns name with provided uuid
*/
func getUserName(uuid uuid.UUID) (string, error) {
	for _, u := range userList {
		if u.UUID == uuid {
			return u.Name, nil
		}
	}
	return "", errors.New("user not found")
}

/*
	Transaction Utilities
	___________________________________________________________________________
*/

/*
	createCoin() creates a payload for creating Tx and updates the values if Tx is successful
*/
func createCoin(sender *uuid.UUID, receiver *uuid.UUID, amount int) ([]byte, error) {
	if sender == &userList[0].UUID && receiver == nil {
		// goofy created a coin
		uuid, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		/*
			Update coin value in user account
		*/
		message := [][]byte{[]byte("Goofy created"), []byte(strconv.Itoa(amount)), []byte("goofy coins with uuid"), []byte(uuid.String())}
		payload := bytes.Join(message, []byte(" "))
		return payload, nil
	}

	/*
		need info of coin object, which coin sender is transferring to receiver
	*/
	senderName, err := getUserName(*sender)
	if err != nil {
		return nil, err
	}
	receiverName, err := getUserName(*receiver)
	if err != nil {
		return nil, err
	}
	message := [][]byte{[]byte(senderName), []byte("paid"), []byte(receiverName), []byte(strconv.Itoa(amount)), []byte("goofy coins")}
	payload := bytes.Join(message, []byte(" "))
	return payload, nil
}

/*
	createTx() appends the payload to Tx slice
*/
func createTx(payload []byte, prevHash []byte) {
	Tx := &transaction{timeStamp: time.Now().Unix(), txMessage: payload, prevHash: prevHash}
	timestamp := []byte(strconv.FormatInt(Tx.timeStamp, 10))
	txData := bytes.Join([][]byte{timestamp, Tx.txMessage, Tx.prevHash}, []byte{})
	hash := sha256.Sum256(txData)
	Tx.currHash = hash[:]
	blk.Tx = append(blk.Tx, Tx)
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

/*
	Exposed APIs
	___________________________________________________________________________
*/

func userAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		type payload struct {
			UserName string `json:"userName"`
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(""))
		}

		var data payload
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(""))
		}

		err = createUser(data.UserName)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(""))
		}
	} else if r.Method == "GET" {
		payload, err := json.Marshal(userList)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(""))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// func testHandler(w http.ResponseWriter, r *http.Request) {
// 	type payload struct {
// 		UserName string `json:"userName"`
// 	}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Print(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(""))
// 	}

// 	var data payload
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		log.Print(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(""))
// 	}

// 	err = createUser([]byte(data.UserName))
// 	if err != nil {
// 		log.Print(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(""))
// 	}

// 	for _, user := range userList {
// 		log.Print(user.uuid)
// 		log.Printf("%s", user.name)
// 		log.Printf("%x", user.privateKey)
// 		log.Printf("%x", user.publicKey)
// 	}
// }

func main() {
	http.HandleFunc("/", reqLogger(indexHandler))
	http.HandleFunc("/dashboard", reqLogger(dashboardHandler))
	http.HandleFunc("/api/user", reqLogger(userAPI))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./assets/css"))))
	log.Printf("App running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

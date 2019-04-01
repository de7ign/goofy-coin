package main

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	_, _, err := generateKeyPair()
	if err != nil {
		t.Error("cannot generate keypair")
	}
}

func TestSignVerify(t *testing.T) {
	priv, pub, err := generateKeyPair()
	if err != nil {
		t.Error("cannot generate kepair")
	}
	payload := []byte("abcdefghijklmnopqrstuvwxyz")
	r, s, err := signTx(priv, payload)
	if err != nil {
		t.Error("cannot sign payload")
	}
	res := verifyTx(pub, payload, r, s)
	if !res {
		t.Error("verification failed")
	}
}

func TestUserUtilities(t *testing.T) {
	err := createUser([]byte("goofy"))
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser([]byte("alice"))
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser([]byte("bob"))
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser([]byte("claire"))
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser([]byte("dave"))
	if err != nil {
		t.Error("cannot create user")
	}

	priv, err := getPrivateKey(userList[0].uuid)
	if err != nil {
		t.Error(err)
	}
	if userList[0].privateKey != priv {
		t.Error("Private key doesnt match")
	}

	pub, err := getPublicKey(userList[0].uuid)
	if err != nil {
		t.Error(err)
	}
	if userList[0].publicKey != pub {
		t.Error("Public key doesnt match")
	}
}
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

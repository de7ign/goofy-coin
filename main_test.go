package main

import (
	"bytes"
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
	err := createUser("goofy")
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser("alice")
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser("bob")
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser("claire")
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser("dave")
	if err != nil {
		t.Error("cannot create user")
	}

	priv, err := getPrivateKey(userList[0].UUID)
	if err != nil {
		t.Error(err)
	}
	if userList[0].privateKey != priv {
		t.Error("Private key doesnt match")
	}

	pub, err := getPublicKey(userList[0].UUID)
	if err != nil {
		t.Error(err)
	}
	if userList[0].publicKey != pub {
		t.Error("Public key doesnt match")
	}
}

func TestTransaction(t *testing.T) {
	/*
		creating user for transaction purpose
	*/
	err := createUser("goofy")
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser("alice")
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser("bob")
	if err != nil {
		t.Error("cannot create user")
	}
	err = createUser("claire")
	if err != nil {
		t.Error("cannot create user")
	}

	payload, err := createCoin(&userList[0].UUID, nil, 10)
	if err != nil {
		t.Error("cannot create payload")
	}
	createTx(payload, nil)
	payload, err = createCoin(&userList[1].UUID, &userList[2].UUID, 10)
	if err != nil {
		t.Error("cannot create payload")
	}
	createTx(payload, blk.Tx[len(blk.Tx)-1].currHash)
	payload, err = createCoin(&userList[2].UUID, &userList[3].UUID, 10)
	if err != nil {
		t.Error("cannot create payload")
	}
	createTx(payload, blk.Tx[len(blk.Tx)-1].currHash)

	for i, ele := range blk.Tx {
		if i != 0 && bytes.Compare(blk.Tx[i].prevHash, blk.Tx[i-1].currHash) != 0 {
			t.Error("error in hashing")
		}
		t.Logf("time stamp : %d", ele.timeStamp)
		t.Logf("tx message : %s", ele.txMessage)
		t.Logf("prevhash   : %x", ele.prevHash)
		t.Logf("currhash   : %x", ele.currHash)
	}
}

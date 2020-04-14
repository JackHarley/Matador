package main

/**
 * Matador Multi-User Messaging Encryption System
 *
 * Copyright (c) 2020, Jack Harley, jackpharley.com
 * All Rights Reserved
 */

import (
	"bufio"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type friendsContainer map[string]friend
type friend struct {
	Name           string
	PKCS1PublicKey string
	Fingerprint    string
}

var fCont friendsContainer

func init() {
	fCont = make(map[string]friend)
}

func friendsSetup() {
	attemptLoadFriends()
}

func saveFriends() {
	f, _ := json.Marshal(fCont)
	e := ioutil.WriteFile("friends.json", f, 0644)
	if e != nil {
		log.Fatal("Failed to save friends to file: " + e.Error())
	}
}

func attemptLoadFriends() {
	f, _ := ioutil.ReadFile("friends.json")
	json.Unmarshal(f, &fCont)
}

func addFriend() {
	if len(os.Args) < 3 {
		fmt.Println("You must provide a path to the public key file for the friend you are adding, " +
			"please try again")
		os.Exit(1)
	}

	publicIn, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println("Unable to find public key at the specified path, please try again.")
		os.Exit(1)
	}

	block, _ := pem.Decode(publicIn)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		fmt.Println("Public key file provided is not valid and/or is corrupt, please try again.")
		os.Exit(1)
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Public key file provided is not valid and/or is corrupt, please try again.")
		os.Exit(1)
	}

	fmt.Println("Key opened successfully.")

	stdin := bufio.NewReader(os.Stdin)
	fmt.Printf("Please enter the name of the person who owns this key: ")
	name, err := stdin.ReadString('\n')
	name = strings.TrimSpace(name)

	fingerprint := sha256.Sum256(x509.MarshalPKCS1PublicKey(publicKey))
	friend := friend{
		Name:           name,
		PKCS1PublicKey: base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(publicKey)),
		Fingerprint:    hex.EncodeToString(fingerprint[:]),
	}

	fCont[friend.Fingerprint] = friend
	saveFriends()
}

func listFriends() {
	if len(fCont) == 0 {
		fmt.Println("You currently have no friends added, please use the \"add <public key path>\" " +
			"command to add one")
		os.Exit(0)
	}

	fmt.Println("The following people are able to decrypt any messages that you currently create:")
	for _, f := range fCont {
		fmt.Printf("- \"%s\", Fingerprint: %s\n", f.Name, f.Fingerprint)
	}
}

func deleteFriend() {
	if len(os.Args) < 3 {
		fmt.Println("You must provide a fingerprint or fingerprint prefix to delete, please try again")
		os.Exit(1)
	}

	fingerprint := os.Args[2]
	if f, ok := fCont[fingerprint]; ok {
		delete(fCont, fingerprint)
		saveFriends()
		fmt.Printf("Deleted friend named %s with fingerprint: %s\n", f.Name, f.Fingerprint)
		return
	}

	for _, f := range fCont {
		if strings.HasPrefix(f.Fingerprint, fingerprint) {
			delete(fCont, f.Fingerprint)
			saveFriends()
			fmt.Printf("Deleted friend named %s with fingerprint: %s\n", f.Name, f.Fingerprint)
			return
		}
	}

	fmt.Printf("Failed to find friend with fingerprint matching %s, please try again\n", fingerprint)
}

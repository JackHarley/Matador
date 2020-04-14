package main

/**
 * Matador Multi-User Messaging Encryption System
 *
 * Copyright (c) 2020, Jack Harley, jackpharley.com
 * All Rights Reserved
 */

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// initIdentify initializes a new identity, prompting for a name on the CLI and generating a
// public private keypair
func initIdentity() {
	if _, err := os.Stat("private_key.pem"); err == nil {
		fmt.Println("You already have an identity generated, please delete your private_key.pem " +
			"file and run the init command again if you would like to create a fresh identity.")
		os.Exit(1)
	}

	stdin := bufio.NewReader(os.Stdin)

	fmt.Printf("Please enter the name you would like to be known as and press enter: ")
	name, err := stdin.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read your name, exiting, please restart the program and try again.")
		os.Exit(1)
	}

	name = strings.TrimSpace(name)
	nameForFile := strings.Replace(name, " ", "", -1)
	nameForFile = strings.ToLower(nameForFile)
	keySize := 2048
	randSource := rand.Reader

	// generate key
	fmt.Printf("Generating key...")
	key, err := rsa.GenerateKey(randSource, keySize)
	if err != nil {
		fmt.Println("Failed to generate keys, please restart the program and try again.")
		os.Exit(1)
	}
	fmt.Printf("done!\n")

	fmt.Printf("Saving keys...")

	// save private key
	privateOut, err := os.Create("private_key.pem")
	defer privateOut.Close()
	if err != nil {
		fmt.Println("Unable to save private key, please restart the program and try again.")
		os.Exit(1)
	}
	privateKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	err = pem.Encode(privateOut, privateKeyPem)
	if err != nil {
		fmt.Println("Unable to save private key, please restart the program and try again.")
		os.Exit(1)
	}

	// save public key
	publicOut, err := os.Create(nameForFile + ".pub.pem")
	defer publicOut.Close()
	if err != nil {
		fmt.Println("Unable to save public key, please restart the program and try again.")
		os.Exit(1)
	}
	publicKeyPem := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	}
	err = pem.Encode(publicOut, publicKeyPem)
	if err != nil {
		fmt.Println("Unable to save public key, please restart the program and try again.")
		os.Exit(1)
	}
	fmt.Printf("done!\n\n")

	fmt.Println("Identity generated succesfully, please make sure to restrict access to this " +
		"directory from unauthorized users.")
	fmt.Printf("You should send your %s.pub.pem key to all other members of your group.\n", nameForFile)
}

func getPrivateKey() *rsa.PrivateKey {
	privateIn, err := ioutil.ReadFile("private_key.pem")
	if err != nil {
		fmt.Println("No private key found, or unable to read it. Please run matador init if you " +
			"need to generate an identity.")
		os.Exit(1)
	}

	block, _ := pem.Decode(privateIn)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		fmt.Println("Private key file is corrupt, please restore it from a backup or " +
			"generate a new identity.")
		os.Exit(1)
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Private key file is corrupt, please restore it from a backup or " +
			"generate a new identity.")
		os.Exit(1)
	}

	return key
}

func getPublicKey() *rsa.PublicKey {
	privateKey := getPrivateKey()
	publicKey := privateKey.Public()
	rsaPublic := publicKey.(*rsa.PublicKey)
	return rsaPublic
}

func getPublicFingerprint() string {
	publicKey := getPublicKey()
	fingerprint := sha256.Sum256(x509.MarshalPKCS1PublicKey(publicKey))
	return hex.EncodeToString(fingerprint[:])
}

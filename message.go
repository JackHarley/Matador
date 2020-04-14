package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

/**
 * Matador Multi-User Messaging Encryption System
 *
 * Copyright (c) 2020, Jack Harley, jackpharley.com
 * All Rights Reserved
 */

type EncryptedMessage struct {
	Nonce                []byte
	Ciphertext           []byte
	EncryptedSessionKeys []EncryptedSessionKey
}
type EncryptedSessionKey struct {
	PublicKeyFingerprint string
	Ciphertext           []byte
}

func encryptMessage() {
	if len(fCont) == 0 {
		fmt.Println("You currently have no friends added, you must add at least one friend before " +
			"encrypting a message, please use the \"add <public key path>\" command to add one")
		os.Exit(0)
	}

	fmt.Printf("Please start typing your message to encrypt, you can use as many lines as you need. " +
		"End your message with a line with nothing but a single . on it:\n\n")
	scanner := bufio.NewScanner(os.Stdin)
	message := ""
	for scanner.Scan() {
		if scanner.Text() == "." {
			break
		}
		message += scanner.Text() + "\n"
	}

	fmt.Printf("\nEncrypting message... (this may take some time if your friends list is large)\n")

	// generate a session key
	sessionKey := make([]byte, 32) // 32 bytes = 256 bit (AES256)
	rand.Read(sessionKey)

	// prepare cipher
	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		fmt.Println("Failed to create AES blockcipher to encrypt message, please try again.")
		os.Exit(1)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Failed to create AES-GCM cipher to encrypt message, please try again.")
		os.Exit(1)
	}
	nonce := make([]byte, 12)
	rand.Read(nonce)

	// encrypt the message
	em := EncryptedMessage{
		Nonce:                nonce,
		Ciphertext:           aesgcm.Seal(nil, nonce, []byte(message), nil),
		EncryptedSessionKeys: make([]EncryptedSessionKey, 0, len(fCont)),
	}

	// encrypt the session key with rsa for each friend's public key
	randSource := rand.Reader
	for _, friend := range fCont {
		pkcs1PubKeyBytes, _ := base64.StdEncoding.DecodeString(friend.PKCS1PublicKey)
		publicKey, _ := x509.ParsePKCS1PublicKey(pkcs1PubKeyBytes)
		ciphertext, _ := rsa.EncryptPKCS1v15(randSource, publicKey, sessionKey)
		esk := EncryptedSessionKey{
			PublicKeyFingerprint: friend.Fingerprint,
			Ciphertext:           ciphertext,
		}
		em.EncryptedSessionKeys = append(em.EncryptedSessionKeys, esk)
	}

	asn1OutBytes, _ := asn1.Marshal(em)
	outPem := &pem.Block{
		Type:  "MATADOR ENCRYPTED MESSAGE",
		Bytes: asn1OutBytes,
	}

	fmt.Printf("Encryption complete, send the following text to your friends:\n\n")
	pem.Encode(os.Stdout, outPem)
}

func decryptMessage() {
	fmt.Printf("Please paste the message to decrypt below and press Enter. You should include " +
		"the BEGIN and END lines:\n\n")

	// read input
	scanner := bufio.NewScanner(os.Stdin)
	data := ""
	for scanner.Scan() {
		data += scanner.Text() + "\n"
		if strings.Contains(scanner.Text(), "-----END MATADOR ENCRYPTED MESSAGE-----") {
			break
		}
	}

	fmt.Printf("\nAttempting to decrypt message...\n")

	// parse
	inPem, _ := pem.Decode([]byte(data))
	if inPem == nil || inPem.Type != "MATADOR ENCRYPTED MESSAGE" {
		fmt.Printf("Unable to recognise the input message, please try again and make sure to " +
			"include the full message including the BEGIN and END lines.")
		os.Exit(1)
	}

	em := EncryptedMessage{}
	_, err := asn1.Unmarshal(inPem.Bytes, &em)
	if err != nil {
		fmt.Printf("The message appears to be corrupted, please ask the sender to re-encrypt " +
			"it and send it again.")
		os.Exit(1)
	}

	// load our private key
	fingerprint := getPublicFingerprint()
	for _, esk := range em.EncryptedSessionKeys {
		if fingerprint == esk.PublicKeyFingerprint {
			privateKey := getPrivateKey()

			// decrypt session key
			sessionKey := make([]byte, 32)
			randSource := rand.Reader
			err := rsa.DecryptPKCS1v15SessionKey(randSource, privateKey, esk.Ciphertext, sessionKey)
			if err != nil {
				fmt.Printf("Failed to decrypt session key, please try again.")
				os.Exit(1)
			}

			// prepare cipher for decrypting message
			block, err := aes.NewCipher(sessionKey)
			if err != nil {
				fmt.Println("Failed to create AES blockcipher to decrypt message, please try again.")
				os.Exit(1)
			}
			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				fmt.Println("Failed to create AES-GCM cipher to decrypt message, please try again.")
				os.Exit(1)
			}

			// decrypt the message
			message, err := aesgcm.Open(nil, em.Nonce, em.Ciphertext, nil)
			if err != nil {
				fmt.Println("Failed to decrypt message, please try again.")
				os.Exit(1)
			}

			fmt.Printf("Message decrypted successfully, printing to console:\n\n%s", string(message))
			return
		}
	}

	fmt.Printf("Your private key has not been granted access to this message, make sure the sender " +
		"has you added as a friend with your public key.")
	os.Exit(1)
}

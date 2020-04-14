package main

/**
 * Matador Multi-User Messaging Encryption System
 *
 * Copyright (c) 2020, Jack Harley, jackpharley.com
 * All Rights Reserved
 */

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		helpText()
		os.Exit(1)
	}

	friendsSetup()

	command := os.Args[1]

	switch command {
	case "init":
		initIdentity()
	case "add":
		addFriend()
	case "list":
		listFriends()
	case "delete":
		deleteFriend()
	case "encrypt":
		encryptMessage()
	case "decrypt":
		decryptMessage()
	}

}

func helpText() {
	fmt.Println("Welcome to Matador")
	fmt.Println("")
	fmt.Println("Usage: matador <command>")
	fmt.Println("")
	fmt.Println("Please use one of the following commands:")
	fmt.Println("    init - Create an identity for yourself and generate associated encryption keys")
	fmt.Println("    add <public key path> - Add a friend to your secret group, this friend will be " +
		"able to read any messages you send")
	fmt.Println("    list - List friends in your secret group")
	fmt.Println("    delete <fingerprint> - Delete a friend from your secret group, they will not be " +
		"able to read any future messages you send")
	fmt.Println("    encrypt - Encrypt a secret message that only your friends will be able to read")
	fmt.Println("    decrypt - Decrypt a secret message sent from someone who has you on their " +
		"friends list")
	fmt.Println()
}

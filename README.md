# Matador

Matador is a multi-user encryption system designed to allow for end-to-end encrypted messaging between a small group of users over an untrusted platform.

The encrypted messages generated consist entirely of characters from the ASCII character set and therefore are suitable for posting on Facebook, in a WhatsApp chat, etc.

Matador requires that users exchange their public keys before commencing messaging, and supports removal of a user from the group (each user must perform the removal indidvidually). Public keys do not need to be communicated over a secure channel, compromise of a public key will not compromise the security of the messages.

Matador was designed in April 2020 by Jack Harley for a university assignment.

Usage
---------------------------
```
Usage: matador <command>

Please use one of the following commands:
    init - Create an identity for yourself and generate associated encryption keys
    add <public key path> - Add a friend to your secret group, this friend will be able to read any messages you send
    list - List friends in your secret group
    delete <fingerprint> - Delete a friend from your secret group, they will not be able to read any future messages you send
    encrypt - Encrypt a secret message that only your friends will be able to read
    decrypt - Decrypt a secret message sent from someone who has you on their friends list
```

What's with the name?
---------------------------
The NSA operates a secret cryptanalysis program named BULLRUN. We hope but cannot confirm that the NSA is unable to access messages shared through Matador.

Legal
---------------------------
Absolutely no permission is granted to Trinity College Dublin students currently taking the Advanced Telecommunications module to access or utilise anything in this repository. Plagiarism is a serious offense and the projects for this module are automatically checked against online sources by TurnItIn. Even looking at the code in this repository prior to writing your own is likely to prejudice you and ultimately result in you writing very similar code. See https://en.wikipedia.org/wiki/Clean_room_design.

This repository is available on GitHub primarily to allow for evaluation of my skill by potential employers, or for those who are genuinely curious and learning a similar topic through self-study. If you are not a TCD student currently taking, or planning to take the Advanced Telecommunications module, then permission is hereby granted, free of charge to use the contents of this repository without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the contents of this repository.
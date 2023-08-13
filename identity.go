package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"log"
	"os"
	"ssb-ng/handshake"
)

// GetIdentity gives you your ed25519 keypair
// Either a new one is generated, or a previous one is reused
// The pre-existing key is assumed to be at `..ssb_keypair` file
func GetIdentity() *handshake.EdKeyPair {
	var keypair handshake.EdKeyPair = handshake.EdKeyPair{}

	if _, err := os.Stat("./.ssb_keypair"); err != nil {
		log.Println("No Identity Found. Creating New Identity")

		idFile, err := os.Create("./.ssb_keypair")
		if err != nil {
			log.Fatalf("Unexpected circumstance. %s", err)
		}
		defer idFile.Close()

		kp, err := handshake.GenEdKeyPair(rand.Reader)
		if err != nil {
			log.Fatalf("Couldn't generate Identity. %s", err)
		}

		idFile.Write(kp.Secret)
		idFile.Write(kp.Public)
		keypair.Public = kp.Public
		keypair.Secret = kp.Secret
	} else {
		idFile, err := os.Open("./.ssb_keypair")

		if err != nil {
			log.Fatalf("Unexpected circumstance. %s", err)
		}
		defer idFile.Close()

		var pubkey [ed25519.PublicKeySize]byte
		var privkey [ed25519.PrivateKeySize]byte
		// order matters
		_, err = idFile.Read(privkey[:])
		if err != nil {
			log.Fatalf("Unexpected circumstance. %s", err)
		}
		_, err = idFile.Read(pubkey[:])
		if err != nil {
			log.Fatalf("Unexpected circumstance. %s", err)
		}
		keypair.Public = pubkey[:]
		keypair.Secret = privkey[:]
	}

	return &keypair
}

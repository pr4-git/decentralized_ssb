package domain

import (
	"crypto/sha512"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/ed25519"
)

const (
	sigSize  int = ed25519.SignatureSize
	hashSize int = sha512.Size
)

type Message struct {
	// autoincrementing id of the message
	Id ulid.ULID
	// The content of the message as opaque bytes
	Content []byte
	// the hash of previous message in chain, nil implies first message
	PreviousHash []byte

	// the sha-512 hash of everything but signature
	Hash []byte
	// message signed by the private key of creator
	Signature []byte
}

func NewMessage(content []byte, previoushash []byte, seckey ed25519.PrivateKey) *Message {
	message := Message{
		Id:           ulid.Make(),
		Content:      content,
		PreviousHash: previoushash,
	}
	message.Signature = ed25519.Sign(seckey, content)

	// hash includes everything but the signature
	// should signature be hashed??
	hasher := sha512.New()
	hasher.Write(message.Id[:])
	hasher.Write(message.Content)
	hasher.Write(message.PreviousHash)
	message.Hash = hasher.Sum(nil)

	return &message
}

package handshake

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
)

// ChallangeLength is the length of the challange message in bytes
const ChallangeLength = 64

// ClientAuthLenght is the length of the clientAuth message in bytes
const ClientAuthLength = 16 + 32 + 64

// ServerAuthLenght is the length of the serverAuth message in bytes
const ServerAuthLength = 16 + 64

// MACLength is the length of MAC in bytes
const MACLength = 16

// GenEdKeyPair generates an ed25519 keypair using the passed reader
// if the reader is nil, it uses crypto/rand.Reader
// for deterministic purposes any reader can be passed, which would
// be useful in tests but insecure in production usecase
func GenEdKeyPair(r io.Reader) (*EdKeyPair, error) {
	if r == nil {
		r = rand.Reader
	}

	pubSrv, secSrv, err := ed25519.GenerateKey(r)
	if err != nil {
		return nil, err
	}

	// check low order and regnerate
	// but i ignore for now

	return &EdKeyPair{Public: pubSrv, Secret: secSrv}, nil
}

// Client performs handshake using the cryptographic identity specified in
// state, using conn in the client role
func Client(state *State, conn io.ReadWriter) (err error) {
	// send challenge
	_, err = conn.Write(state.createChallange())
	if err != nil {
		return fmt.Errorf("error sending challenge by client (Error: %v)", err)
	}

	// recv challenge
	challengeResp := make([]byte, ChallangeLength)
	_, err = io.ReadFull(conn, challengeResp)
	if err != nil {
		return fmt.Errorf("error receiving challenge by client (Error: %v)", err)
	}

	// verify challenge
	if !state.verifyChallange(challengeResp) {
		return fmt.Errorf("cannot verify server challange")
	}

	// send authentication vector
	_, err = conn.Write(state.createClientAuth())
	if err != nil {
		return fmt.Errorf("error sending client hello (Error: %v)", err)
	}

	// recv authentication vector
	boxedSig := make([]byte, ServerAuthLength)
	_, err = io.ReadFull(conn, boxedSig)
	if err != nil {
		return fmt.Errorf("error receiving server auth (Error: %v)", err)
	}

	// authanticate remote
	if !state.verifyServerAccept(boxedSig) {
		return fmt.Errorf("connot verify server accept")
	}

	state.cleanSecrets()
	return nil
}

// Server performs handshake using the cryptographic identity specified in
// state, using conn in the server role
func Server(state *State, conn io.ReadWriter) (err error) {

	// recv challenge
	challenge := make([]byte, ChallangeLength)
	_, err = io.ReadFull(conn, challenge)
	if err != nil {
		return fmt.Errorf("error receiving challenge by server (Error: %v)", err)
	}

	// verify challenge
	if !state.verifyChallange(challenge) {
		return fmt.Errorf("cannot verify client challange")
	}

	// send challenge
	_, err = conn.Write(state.createChallange())
	if err != nil {
		return fmt.Errorf("error sending challenge by server (Error: %v)", err)
	}

	// recv authentication vector
	hello := make([]byte, ClientAuthLength)
	_, err = io.ReadFull(conn, hello)
	if err != nil {
		return fmt.Errorf("error receiving client hello (Error: %v)", err)
	}

	// authanticate remote
	if !state.verifyClientAuth(hello) {
		return fmt.Errorf("connot verify client auth")
	}

	// accept
	_, err = conn.Write(state.createServerAccept())
	if err != nil {
		return fmt.Errorf("error sending server accept (Error: %v)", err)
	}

	state.cleanSecrets()
	return nil
}

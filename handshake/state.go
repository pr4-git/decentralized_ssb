package handshake

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"itsy/ssb/handshake/internal/extra25519"
	"log"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/auth"
	"golang.org/x/crypto/nacl/box"
)

// State is the state of peer during handshake process
type State struct {
	appKey [32]byte

	secHash      []byte
	localAppmac  [32]byte
	remoteAppMac []byte

	localExchange  CurveKeyPair
	local          EdKeyPair
	remoteExchange CurveKeyPair
	remotePublic   ed25519.PublicKey // logterm key

	// eph secrets
	secret, secret2, secret3 [32]byte

	hello []byte

	helloAlice, helloBob [32]byte
}

type EdKeyPair struct {
	Public ed25519.PublicKey
	Secret ed25519.PrivateKey
}

func NewKeyPair(public, secret []byte) (*EdKeyPair, error) {
	var keypair EdKeyPair
	if n := len(secret); n != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("Invalid private-key size of %db", n)
	}

	keypair.Secret = secret

	if n := len(public); n != ed25519.PublicKeySize {
		return nil, fmt.Errorf("Invalid public-key size of %db", n)
	}
	// TODO: check EdLowOrder (???)
	/*
		if lo25519.IsEdLowOrder(public){
			return nil, fmt.Errorf("Invalid keypair");
		}
	*/
	keypair.Public = public

	return &keypair, nil
}

type CurveKeyPair struct {
	Public [32]byte
	Secret [32]byte
}

func NewClientState(appKey []byte, local EdKeyPair, remotePublic ed25519.PublicKey) (*State, error) {
	state, err := newState(appKey, local)
	if err != nil {
		return state, fmt.Errorf("error creating client state (ERROR: %v)", err)
	}

	state.remotePublic = remotePublic
	if l := len(state.remotePublic); l != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid key size for remote/public of %db", l)
	}

	return state, err
}

func NewServerState(appKey []byte, local EdKeyPair) (*State, error) {
	return newState(appKey, local)
}

// newState initializes the state for both client and server
func newState(appKey []byte, local EdKeyPair) (*State, error) {
	pubKey, secKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error generating keypair for new state (Error: %v)", err)
	}

	s := State{
		remotePublic: make([]byte, ed25519.PublicKeySize),
	}

	copy(s.appKey[:], appKey)
	copy(s.localExchange.Public[:], pubKey[:])
	copy(s.localExchange.Secret[:], secKey[:])
	s.local = local

	if l := len(s.local.Public); l != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid key size for ephemeral/public of %db", l)
	}

	if l := len(s.local.Secret); l != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid key size for ephemeral/private of %db", l)
	}

	return &s, nil
}

// createChallange returns a buffer with a challange
func (s *State) createChallange() []byte {
	mac := auth.Sum(s.localExchange.Public[:], &s.appKey)
	copy(s.localAppmac[:], mac[:])

	return append(s.localAppmac[:], s.localExchange.Public[:]...)
}

func (s *State) verifyChallange(ch []byte) bool {
	mac := ch[:32]
	remoteEphPubKey := ch[32:]

	ok := auth.Verify(mac, remoteEphPubKey, &s.appKey)

	copy(s.remoteExchange.Public[:], remoteEphPubKey)
	s.remoteAppMac = mac

	// is supposed to be [32]byte in size
	var sec []byte

	// THIS might set all dst to zero on low-order
	// BEWARE!!!!
	//curve25519.ScalarMult(&sec, &s.localExchange.Secret, &s.remoteExchange.Public)
	sec, err := curve25519.X25519(s.localExchange.Secret[:], s.remoteExchange.Public[:])
	if err != nil {
		log.Fatalf("Couldn't verify challange (Error: %s)", err)
	}

	if len(sec) != 32 {
		log.Fatalf("Derived key is not 32b. Unexpected error!!")
	}

	copy(s.secret[:], sec)

	secHasher := sha256.New()
	secHasher.Write(s.secret[:])
	s.secHash = secHasher.Sum(nil)

	return ok
}

func (s *State) createClientAuth() []byte {
	var curveRemotePubKey [32]byte

	if !extra25519.PublicKeyToCurve25519(&curveRemotePubKey, s.remotePublic) {
		log.Fatalf("Couldn't convert remote to curve key when creating client auth")
	}

	var helloBob []byte

	helloBob, err := curve25519.X25519(s.localExchange.Secret[:], curveRemotePubKey[:])
	if err != nil {
		log.Fatalf("Couldn't verify challange (Error: %s)", err)
	}

	if len(helloBob) != 32 {
		log.Fatalf("Derived key is not 32b. Unexpected error!!")
	}
	copy(s.helloBob[:], helloBob)

	secHasher := sha256.New()
	secHasher.Write(s.appKey[:])
	secHasher.Write(s.secret[:])
	secHasher.Write(s.helloBob[:])
	copy(s.secret2[:], secHasher.Sum(nil))

	var sigMsg bytes.Buffer
	sigMsg.Write(s.appKey[:])
	sigMsg.Write(s.remotePublic[:])
	sigMsg.Write(s.secHash)

	sig := ed25519.Sign(s.local.Secret, sigMsg.Bytes())

	var helloBuf bytes.Buffer
	helloBuf.Write(sig[:])
	helloBuf.Write(s.local.Public[:])
	s.hello = helloBuf.Bytes()

	out := make([]byte, 0, len(s.hello)-box.Overhead)
	var nonce [24]byte
	out = box.SealAfterPrecomputation(out, s.hello, &nonce, &s.secret2)

	return out
}

var nullHello [ed25519.SignatureSize + ed25519.PublicKeySize]byte

func (s *State) verifyClientAuth(data []byte) bool {
	var cvSec [32]byte
	extra25519.PrivateKeyToCurve25519(&cvSec, s.local.Secret)

	helloBob, err := curve25519.X25519(cvSec[:], s.remoteExchange.Public[:])
	if err != nil {
		log.Fatalf("Couldn't verify challange (Error: %s)", err)
	}
	if len(helloBob) != 32 {
		log.Fatalf("Derived key is not 32b. Unexpected error!!")
	}
	copy(s.helloBob[:], helloBob)

	sechasher := sha256.New()
	sechasher.Write(s.appKey[:])
	sechasher.Write(s.secret[:])
	sechasher.Write(s.helloBob[:])
	copy(s.secret2[:], sechasher.Sum(nil))

	s.hello = make([]byte, 0, len(data)-16)

	var nonce [24]byte // always zero???
	var openOk bool
	s.hello, openOk = box.OpenAfterPrecomputation(s.hello, data, &nonce, &s.secret2)

	var sig = make([]byte, ed25519.SignatureSize)
	var public = make([]byte, ed25519.PublicKeySize)

	if openOk {
		copy(sig, s.hello[:ed25519.SignatureSize])
		copy(public[:], s.hello[ed25519.SignatureSize:])
	} else {
		copy(sig, nullHello[:ed25519.SignatureSize])
		copy(public[:], nullHello[ed25519.SignatureSize:])
	}

	// check low order here

	var sigMsg bytes.Buffer
	sigMsg.Write(s.appKey[:])
	sigMsg.Write(s.local.Public[:])
	verifyOk := ed25519.Verify(public, sigMsg.Bytes(), sig)

	copy(s.remotePublic, public)
	return openOk && verifyOk
}

func (s *State) createServerAccept() []byte {
	var curveRemotePubKey [32]byte

	if !extra25519.PublicKeyToCurve25519(&curveRemotePubKey, s.remotePublic) {
		log.Fatalf("Couldn't convert remote to curve key when creating server accept")
	}

	var helloAlice []byte

	helloAlice, err := curve25519.X25519(s.localExchange.Secret[:], curveRemotePubKey[:])
	if err != nil {
		log.Fatalf("Couldn't create server accept (Error: %s)", err)
	}

	if len(helloAlice) != 32 {
		log.Fatalf("Derived key is not 32b. Unexpected error!!")
	}
	copy(s.helloAlice[:], helloAlice)

	secHasher := sha256.New()
	secHasher.Write(s.appKey[:])
	secHasher.Write(s.secret[:])
	secHasher.Write(s.helloBob[:])
	secHasher.Write(s.helloAlice[:])
	copy(s.secret3[:], secHasher.Sum(nil))

	var sigMsg bytes.Buffer
	sigMsg.Write(s.appKey[:])
	sigMsg.Write(s.hello[:])
	sigMsg.Write(s.secHash)

	// acceptance message
	okay := ed25519.Sign(s.local.Secret, sigMsg.Bytes())

	out := make([]byte, 0, len(okay)+16)
	var nonce [24]byte
	return box.SealAfterPrecomputation(out, okay[:], &nonce, &s.secret3)
}

func (s *State) verifyServerAccept(boxedOkay []byte) bool {
	var curveLocalSec [32]byte
	extra25519.PrivateKeyToCurve25519(&curveLocalSec, s.local.Secret)

	helloAlice, err := curve25519.X25519(curveLocalSec[:], s.remoteExchange.Public[:])
	if err != nil {
		log.Fatalf("Couldn't verify server accept (Error: %s)", err)
	}
	if len(helloAlice) != 32 {
		log.Fatalf("Derived key is not 32b. Unexpected error!!")
	}
	copy(s.helloAlice[:], helloAlice)

	sechasher := sha256.New()
	sechasher.Write(s.appKey[:])
	sechasher.Write(s.secret[:])
	sechasher.Write(s.helloBob[:])
	sechasher.Write(s.helloAlice[:])
	copy(s.secret3[:], sechasher.Sum(nil))

	var nonce [24]byte // always zero???
	//sig := make([]byte, 0, len(boxedOkay)-16)
	// shadowed below ig??
	sig, openOk := box.OpenAfterPrecomputation(nil, boxedOkay, &nonce, &s.secret3)

	var sigMsg bytes.Buffer
	sigMsg.Write(s.appKey[:])
	sigMsg.Write(s.hello[:])
	sigMsg.Write(s.secHash)

	verifyOk := ed25519.Verify(s.remotePublic, sigMsg.Bytes(), sig)

	return verifyOk && openOk
}

// Clean secrets overwrites all intermediate secrets
// and copied the final secret to s.secret
func (s *State) cleanSecrets() {
	var zeros [64]byte

	copy(s.secHash, zeros[:])
	copy(s.secret[:], zeros[:])
	copy(s.helloBob[:], zeros[:])
	copy(s.helloAlice[:], zeros[:])

	h := sha256.New()
	h.Write(s.secret3[:])
	copy(s.secret[:], h.Sum(nil))
	copy(s.secret2[:], zeros[:])
	copy(s.secret3[:], zeros[:])
	copy(s.localExchange.Secret[:], zeros[:])
}

func (s *State) Remote() []byte {
	return s.remotePublic[:]
}

// GetBoxStreamEncKeys returns the encryption key and nonce suitable for boxstream
func (s *State) GetBoxstreamEncKeys() ([32]byte, [24]byte) {
	// Error before cleansecrets is called mayhaps
	// because it'd be insecure

	var enKey [32]byte
	h := sha256.New()
	h.Write(s.secret[:])
	h.Write(s.remotePublic[:])
	copy(enKey[:], h.Sum(nil))

	var nonce [24]byte
	copy(nonce[:], s.remoteAppMac)
	return enKey, nonce
}

// GetBoxStreamDecKeys returns the decryption key and nonce suitable for boxstream
func (s *State) GetBoxstreamDecKeys() ([32]byte, [24]byte) {
	// Error before cleansecrets is called mayhaps
	// because it'd be insecure

	var deKey [32]byte
	h := sha256.New()
	h.Write(s.secret[:])
	h.Write(s.local.Public[:])
	copy(deKey[:], h.Sum(nil))

	var nonce [24]byte
	copy(nonce[:], s.localAppmac[:])
	return deKey, nonce
}

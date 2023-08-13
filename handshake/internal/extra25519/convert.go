package extra25519

import (
	"crypto/sha512"

	"golang.org/x/crypto/ed25519"

	"filippo.io/edwards25519"
)

// PrivateKeyToCurve25519 converts an ed25519 private keu to a correspoding
// curve25519 private key such that the resulting curve25519 public key will equal
// the result from PublicKeyToCurve25519
func PrivateKeyToCurve25519(curve25519Private *[32]byte, privateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(privateKey[:32])
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(curve25519Private[:], digest)
}

// PublicKeyToCurve25519 converts an Ed25519 public key into curve25519
// public key that would be generated from the same private key
func PublicKeyToCurve25519(curve25519Public *[32]byte, publicKey ed25519.PublicKey) bool {
	// TODO: Check low order
	edPoint, err := new(edwards25519.Point).SetBytes(publicKey)

	if err != nil {
		return false
	}

	copy(curve25519Public[:], edPoint.BytesMontgomery())

	return true
}

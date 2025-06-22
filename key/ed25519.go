package key

import (
	"crypto/ed25519"
	"fmt"
)

type KeyEd25519 struct {
	Private ed25519.PrivateKey
}

func generateEd25519() (*KeyEd25519, error) {
	_, pri, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return makeEd25519(pri)
}

func makeEd25519(pri ed25519.PrivateKey) (*KeyEd25519, error) {
	return &KeyEd25519{
		Private: pri,
	}, nil
}

func fromBytesEd25519(priBytes []byte) (*KeyEd25519, error) {
	if len(priBytes) != 32 {
		return nil, fmt.Errorf("incorrect length of Ed25519 private key: %v", len(priBytes))
	}
	return makeEd25519(priBytes)
}

func (k *KeyEd25519) keyType() KeyType {
	return Ed25519
}

func (k *KeyEd25519) PrivateBytes() []byte {
	return k.Private
}

func (k *KeyEd25519) PublicBytes() []byte {
	return k.Private[32:]
}

func (k *KeyEd25519) Sign(message []byte, context string) ([]byte, error) {
	return k.Private.Sign(nil, message, &ed25519.Options{
		Context: context,
	})
}

func verifyEd25519(pub ed25519.PublicKey, sig []byte, message []byte, context string) error {
	return ed25519.VerifyWithOptions(pub, message, sig, &ed25519.Options{
		Context: context,
	})
}

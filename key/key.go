package key

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/multiformats/go-multibase"
	"github.com/trustbloc/did-go/doc/did"
	"github.com/trustbloc/did-go/method/key"
)

type Key interface {
	PrivateBytes() []byte
	PublicBytes() []byte
	Did() string
	Resolution() did.DocResolution
	Sign(message []byte, context string) ([]byte, error)
	WritePrivateKey(filename string) error
}

type KeyType int

const (
	Ed25519 KeyType = iota
)

type KeyEd25519 struct {
	Private ed25519.PrivateKey
}

func (k *KeyEd25519) PrivateBytes() []byte {
	return k.Private
}

func (k *KeyEd25519) PublicBytes() []byte {
	return k.Private[32:]
}

func (k *KeyEd25519) Did() string {
	prefixedKey := append([]byte{0xED, 0x01}, k.PublicBytes()...)
	str, err := multibase.Encode(multibase.Base58BTC, prefixedKey)
	if err != nil {
		panic(err)
	}
	return "did:key:" + str
}

func (k *KeyEd25519) Resolution() did.DocResolution {
	resolution, err := ResolveDid(k.Did())
	if err != nil {
		panic(err)
	}
	return *resolution
}

func (k *KeyEd25519) Sign(message []byte, context string) ([]byte, error) {
	return k.Private.Sign(nil, message, &ed25519.Options{
		Context: context,
	})
}

func (k *KeyEd25519) WritePrivateKey(filename string) error {
	bytes, err := x509.MarshalPKCS8PrivateKey(k.Private)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: bytes,
	}
	handle, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(handle, pemBlock)
	if err != nil {
		return err
	}
	return handle.Close()
}

func GenerateKey(keyType KeyType) (Key, error) {
	switch keyType {
	case Ed25519:
		{
			_, pri, err := ed25519.GenerateKey(nil)
			if err != nil {
				return nil, err
			}
			return &KeyEd25519{
				Private: pri,
			}, nil

		}
	default:
		return nil, fmt.Errorf("invalid key type: %v", keyType)
	}
}

func ReadPrivateKey(filename string) (Key, error) {
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, rest := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM data present")
	} else if len(rest) > 0 {
		return nil, fmt.Errorf("extra PEM data present")
	} else if block.Type != "ED25519 PRIVATE KEY" {
		return nil, fmt.Errorf("wrong PEM block type")
	}
	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &KeyEd25519{
		Private: pri.(ed25519.PrivateKey),
	}, nil
}

func PublicKeyFromDid(did string) ([]byte, error) {
	if !strings.HasPrefix(did, "did:key:") {
		return nil, fmt.Errorf("invalid DID format")
	}
	str := strings.TrimPrefix(did, "did:key:")
	_, data, err := multibase.Decode(str)
	if err != nil {
		return nil, fmt.Errorf("multibase decode error: %v", err)
	}
	if len(data) != 34 || data[0] != 0xED || data[1] != 0x01 {
		return nil, fmt.Errorf("not a valid ed25519 multicodec key")
	}
	return data[2:], nil
}

func ResolveDid(did string) (*did.DocResolution, error) {
	return key.New().Read(did)
}

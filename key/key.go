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
	keyType() KeyType
	PrivateBytes() []byte
	PublicBytes() []byte
	Sign(message []byte, context string) ([]byte, error)
}

type KeyType int

const (
	UnknownKeyType KeyType = iota
	Ed25519
	Bls12381
)

func prefixBytes(keyType KeyType) []byte {
	switch keyType {
	case Ed25519:
		return []byte{0xED, 0x01}
	case Bls12381:
		return []byte{0xEB, 0x01}
	default:
		panic(fmt.Errorf("invalid key type: %v", keyType))
	}
}

func blockType(keyType KeyType) string {
	switch keyType {
	case Ed25519:
		return "ED25519 PRIVATE KEY"
	case Bls12381:
		return "BLS12-381 PRIVATE KEY"
	default:
		panic(fmt.Errorf("invalid key type: %v", keyType))
	}
}

func GenerateKey(keyType KeyType) (Key, error) {
	switch keyType {
	case Ed25519:
		return generateEd25519()
	case Bls12381:
		return generateBls12381()
	default:
		return nil, fmt.Errorf("invalid key type: %v", keyType)
	}
}

func Verify(did string, sig []byte, message []byte, context string) error {
	keyType, pubBytes, err := PublicKeyFromDid(did)
	if err != nil {
		return err
	}
	switch keyType {
	case Ed25519:
		{
			verifyEd25519(pubBytes, sig, message, context)
		}
	}
	return fmt.Errorf("invalid key type: %v", keyType)
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
	} else if block.Type == blockType(Ed25519) {
		pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			keyEd25519, okay := pri.(ed25519.PrivateKey)
			if okay {
				return makeEd25519(keyEd25519)
			} else {
				return fromBytesEd25519(block.Bytes)
			}
		} else {
			return fromBytesEd25519(block.Bytes)
		}
	} else if block.Type == blockType(Bls12381) {
		return fromBytesBls12381(block.Bytes)
	}
	return nil, fmt.Errorf("wrong PEM block type")
}

func WritePrivateKey(k Key, filename string) error {
	var bytes []byte
	keyEd25519, okay := k.(*KeyEd25519)
	if okay {
		var err error
		bytes, err = x509.MarshalPKCS8PrivateKey(keyEd25519.Private)
		if err != nil {
			bytes = k.PrivateBytes()
		}
	} else {
		bytes = k.PrivateBytes()
	}
	block := &pem.Block{
		Type:  blockType(k.keyType()),
		Bytes: bytes,
	}
	handle, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(handle, block)
	if err != nil {
		return err
	}
	return handle.Close()
}

func Did(k Key) string {
	prefixedKey := append(prefixBytes(k.keyType()), k.PublicBytes()...)
	str, err := multibase.Encode(multibase.Base58BTC, prefixedKey)
	if err != nil {
		panic(err)
	}
	return "did:key:" + str
}

func ResolveDid(did string) (*did.DocResolution, error) {
	return key.New().Read(did)
}

func PublicKeyFromDid(did string) (KeyType, []byte, error) {
	if !strings.HasPrefix(did, "did:key:") {
		return UnknownKeyType, nil, fmt.Errorf("invalid DID format")
	}
	str := strings.TrimPrefix(did, "did:key:")
	_, data, err := multibase.Decode(str)
	if err != nil {
		return UnknownKeyType, nil, fmt.Errorf("multibase decode error: %v", err)
	}
	if len(data) < 2 {
		return UnknownKeyType, nil, fmt.Errorf("invalid multicodec: %x", data)
	}
	firstTwoMatch := func(x []byte, y []byte) bool {
		return x[0] == y[0] && x[1] == y[1]
	}
	if firstTwoMatch(data, prefixBytes(Ed25519)) {
		return Ed25519, data[2:], nil
	}
	if firstTwoMatch(data, prefixBytes(Bls12381)) {
		return Bls12381, data[2:], nil
	}
	return UnknownKeyType, nil, fmt.Errorf("unsupported multicodec key: %x", data[:2])
}

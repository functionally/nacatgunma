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

type Key struct {
	Private    ed25519.PrivateKey
	Public     ed25519.PublicKey
	Did        string
	Resolution did.DocResolution
}

func GenerateKey() (*Key, error) {
	_, pri, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return makeKey(pri)
}

func makeKey(pri []byte) (*Key, error) {
	pub := pri[32:]
	prefixedKey := append([]byte{0xED, 0x01}, pub...)
	str, err := multibase.Encode(multibase.Base58BTC, prefixedKey)
	if err != nil {
		return nil, err
	}
	did := "did:key:" + str
	resolution, err := key.New().Read(did)
	if err != nil {
		return nil, err
	}
	return &Key{
		Private:    pri,
		Public:     pub,
		Did:        did,
		Resolution: *resolution,
	}, nil
}

func (key *Key) Sign(message []byte) ([]byte, error) {
	return key.Private.Sign(nil, message, nil)
}

func (key *Key) WritePrivateKey(filename string) error {
	bytes, err := x509.MarshalPKCS8PrivateKey(key.Private)
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

func ReadPrivateKey(filename string) (*Key, error) {
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
	return makeKey(pri.(ed25519.PrivateKey))
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

package achain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func WritePrivateKey(filename string, key *ecdsa.PrivateKey) error {
	bytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
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

func ReadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, rest := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM data present")
	} else if len(rest) > 0 {
		return nil, fmt.Errorf("extra PEM data present")
	} else if block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("wrong PEM block type")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

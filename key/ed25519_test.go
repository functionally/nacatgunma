package key

import (
	"bytes"
	"reflect"
	"testing"
)

func TestGenerateKey_Ed25519(t *testing.T) {
	testGenerateKey(Ed25519, reflect.TypeOf(&KeyEd25519{}), t)
}

func TestDid_Ed25519(t *testing.T) {
	testDid(Ed25519, "did:key:z6M", t)
}

func TestPubFromDid_Ed25519(t *testing.T) {
	testPubFromDid(Ed25519, t)
}

func TestDidResolution_Ed25519(t *testing.T) {
	testDidResolution(Ed25519, "Ed25519VerificationKey2018", t)
}

func TestIO_Ed25519(t *testing.T) {
	testIO(
		Ed25519,
		func(kx Key, ky Key) {
			kx1, ok := kx.(*KeyEd25519)
			if !ok {
				t.Errorf("incorrect key type generated: %T", kx1)
			}
			ky1, ok := ky.(*KeyEd25519)
			if !ok {
				t.Errorf("incorrect key type read: %T", ky1)
			}
			if !bytes.Equal(kx1.Private, ky1.Private) {
				t.Error("private key does not match")
			}
		},
		t,
	)
}

func TestVerify_Ed25519(t *testing.T) {
	testVerify(Ed25519, t)
}

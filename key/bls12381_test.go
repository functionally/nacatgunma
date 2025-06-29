package key

import (
	"bytes"
	"os"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestDid(t *testing.T) {
	k0, err := GenerateKey(Bls12381)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	if d0[:12] != "did:key:zUC7" {
		t.Errorf("incorrect DID encoding: %v", d0)
	}
}

func TestPubFromDid(t *testing.T) {
	k0, err := GenerateKey(Bls12381)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	t1, p1, err := PublicKeyFromDid(d0)
	if err != nil {
		t.Error(err)
	}
	if t1 != Bls12381 {
		t.Errorf("incorrect key type: %v", t1)
	}
	if !bytes.Equal(k0.PublicBytes(), p1) {
		t.Error("public key does not match")
	}
}

func TestDidResolution(t *testing.T) {
	k0, err := GenerateKey(Bls12381)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	r1, err := ResolveDid(d0)
	if err != nil {
		t.Error(err)
	}
	if r1.DIDDocument.VerificationMethod[0].Type != "Bls12381G2Key2020" {
		t.Errorf("incorrect verification method type: %v", r1.DIDDocument.VerificationMethod[0].Type)
	}
}

func TestIO(t *testing.T) {
	k0, err := GenerateKey(Bls12381)
	if err != nil {
		t.Error(err)
	}
	k1, ok := k0.(*KeyBls12381)
	if !ok {
		t.Errorf("incorrect key type generated: %T", k1)
	}
	kFile, err := os.CreateTemp(".", "tmp-*.pem")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(kFile.Name())
	kFile.Close()
	err = WritePrivateKey(k0, kFile.Name())
	if err != nil {
		t.Error(err)
	}
	k2, err := ReadPrivateKey(kFile.Name())
	if err != nil {
		t.Error(err)
	}
	k3, ok := k2.(*KeyBls12381)
	if !ok {
		t.Errorf("incorrect key type read: %T", k1)
	}
	if k1.Private != k3.Private {
		t.Error("private key does not match")
	}
	if k1.Public != k3.Public {
		t.Error("public key does not match")
	}
}

func TestVerify(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)
	k0, err := GenerateKey(Bls12381)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	properties.Property("sign and verify", prop.ForAll(
		func(ctx string, msg string) bool {
			sig, err := k0.Sign([]byte(msg), ctx)
			if err != nil {
				return false
			}
			return Verify(d0, sig, []byte(msg), ctx) == nil
		},
		gen.AnyString(),
		gen.AnyString(),
	))
	properties.Property("sign and not verify", prop.ForAll(
		func(ctx string, msg0 string, msg1 string) bool {
			sig, err := k0.Sign([]byte(msg0), ctx)
			if err != nil {
				return false
			}
			return msg0 == msg1 || Verify(d0, sig, []byte(msg1), ctx) != nil
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.AnyString(),
	))
	properties.TestingRun(t)
}

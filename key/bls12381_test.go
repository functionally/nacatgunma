package key

import (
	"reflect"
	"testing"
)

func TestGenerateKey_Bls12381(t *testing.T) {
	testGenerateKey(Bls12381, reflect.TypeOf(&KeyBls12381{}), t)
}

func TestDid_Bls12381(t *testing.T) {
	testDid(Bls12381, "did:key:zUC7", t)
}

func TestPubFromDid_Bls12381(t *testing.T) {
	testPubFromDid(Bls12381, t)
}

func TestDidResolution_Bls12381(t *testing.T) {
	testDidResolution(Bls12381, "Bls12381G2Key2020", t)
}

func TestIO_Bls12381(t *testing.T) {
	testIO(
		Bls12381,
		func(kx Key, ky Key) {
			kx1, ok := kx.(*KeyBls12381)
			if !ok {
				t.Errorf("incorrect key type generated: %T", kx1)
			}
			ky1, ok := ky.(*KeyBls12381)
			if !ok {
				t.Errorf("incorrect key type read: %T", ky1)
			}
			if kx1.Private != ky1.Private {
				t.Error("private key does not match")
			}
			if kx1.Public != ky1.Public {
				t.Error("public key does not match")
			}
		},
		t,
	)
}

func TestVerify_Bls12381(t *testing.T) {
	testVerify(Bls12381, t)
}

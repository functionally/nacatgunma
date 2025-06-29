package key

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func testGenerateKey(keyType KeyType, kT reflect.Type, t *testing.T) {
	k0, err := GenerateKey(keyType)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(k0) != kT {
		t.Errorf("incorrect key type generated: %T", k0)
	}
}

func testDid(keyType KeyType, prefix string, t *testing.T) {
	k0, err := GenerateKey(keyType)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	if d0[:len(prefix)] != prefix {
		t.Errorf("incorrect DID encoding: %v", d0)
	}
}

func testDidResolution(keyType KeyType, keyTypeName string, t *testing.T) {
	k0, err := GenerateKey(keyType)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	r1, err := ResolveDid(d0)
	if err != nil {
		t.Error(err)
	}
	if r1.DIDDocument.VerificationMethod[0].Type != keyTypeName {
		t.Errorf("incorrect verification method type: %v", r1.DIDDocument.VerificationMethod[0].Type)
	}
}

func testPubFromDid(keyType KeyType, t *testing.T) {
	k0, err := GenerateKey(keyType)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	t1, p1, err := PublicKeyFromDid(d0)
	if err != nil {
		t.Error(err)
	}
	if t1 != keyType {
		t.Errorf("incorrect key type: %v", t1)
	}
	if !bytes.Equal(k0.PublicBytes(), p1) {
		t.Error("public key does not match")
	}
}

func testIO(keyType KeyType, equals func(kx Key, ky Key), t *testing.T) {
	k0, err := GenerateKey(keyType)
	if err != nil {
		t.Error(err)
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
	k1, err := ReadPrivateKey(kFile.Name())
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(k0) != reflect.TypeOf(k1) {
		t.Errorf("invalid key type read: %T", k1)
	}
	equals(k0, k1)
}

func testVerify(keyType KeyType, t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)
	k0, err := GenerateKey(keyType)
	if err != nil {
		t.Error(err)
	}
	d0 := Did(k0)
	properties.Property("sign and verify", prop.ForAll(
		func(ctx string, msg []byte) bool {
			if len(ctx) > 255 {
				return true
			}
			sig, err := k0.Sign(msg, ctx)
			if err != nil {
				return false
			}
			z := Verify(d0, sig, msg, ctx)
			return z == nil
		},
		genContext(),
		genBytes(0, 1000),
	))
	properties.Property("sign and not verify", prop.ForAll(
		func(ctx string, msg0 []byte, msg1 []byte) bool {
			if len(ctx) > 255 {
				return true
			}
			sig, err := k0.Sign(msg0, ctx)
			if err != nil {
				return false
			}
			return bytes.Equal(msg0, msg1) || Verify(d0, sig, msg1, ctx) != nil
		},
		genContext(),
		genBytes(0, 1000),
		genBytes(0, 1000),
	))
	properties.TestingRun(t)
}

func genBytes(minLen int, maxLen int) gopter.Gen {
	return gen.IntRange(minLen, maxLen).FlatMap(
		func(n interface{}) gopter.Gen {
			return gen.SliceOfN(
				n.(int),
				gen.UInt8()).Map(func(v []uint8) []byte { return []byte(v) })
		},
		reflect.TypeOf([]byte{}),
	)
}

func genContext() gopter.Gen {
	return genString(0, 255)
}

func genString(minLen int, maxLen int) gopter.Gen {
	return genBytes(minLen, maxLen).Map(func(v []uint8) string { return string(v) })
}

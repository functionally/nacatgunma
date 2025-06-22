package key

import (
	"crypto/rand"
	"fmt"

	bls12381 "github.com/kilic/bls12-381"
)

type KeyBls12381 struct {
	Private bls12381.Fr
	Public  bls12381.PointG1
}

func makeBls12381(pri *bls12381.Fr) (*KeyBls12381, error) {
	g1 := bls12381.NewG1()
	return &KeyBls12381{
		Private: *pri,
		Public:  *g1.MulScalar(g1.New(), g1.One(), pri),
	}, nil
}

func generateBls12381() (*KeyBls12381, error) {
	pri, err := bls12381.NewFr().Rand(rand.Reader)
	if err != nil {
		return nil, err
	}
	return makeBls12381(pri)
}

func fromBytesBls12381(priBytes []byte) (*KeyBls12381, error) {
	if len(priBytes) != 32 {
		return nil, fmt.Errorf("incorrect length of BLS-12381 private key: %v", len(priBytes))
	}
	return makeBls12381(bls12381.NewFr().FromBytes(priBytes))
}

func (k *KeyBls12381) keyType() KeyType {
	return Bls12381
}

func (k *KeyBls12381) PrivateBytes() []byte {
	return k.Private.ToBytes()
}

func (k *KeyBls12381) PublicBytes() []byte {
	return bls12381.NewG1().ToCompressed(&k.Public)
}

func (k *KeyBls12381) Sign(message []byte, context string) ([]byte, error) {
	g2 := bls12381.NewG2()
	point, err := hashToCurveBls12381(message, context)
	if err != nil {
		return nil, err
	}
	sig := g2.New()
	g2.MulScalar(sig, point, &k.Private)
	return g2.ToCompressed(sig), nil
}

func pointG1FromBytesBls12381(bytes []byte) (*bls12381.PointG1, error) {
	g1 := bls12381.NewG1()
	return g1.FromCompressed(bytes)
}

func pointG2FromBytesBls12381(bytes []byte) (*bls12381.PointG2, error) {
	g2 := bls12381.NewG2()
	return g2.FromCompressed(bytes)
}

func verifyBls12381(pub *bls12381.PointG1, sig *bls12381.PointG2, message []byte, context string) error {
	g1 := bls12381.NewG1()
	point, err := hashToCurveBls12381(message, context)
	if err != nil {
		return err
	}
	engine := bls12381.NewEngine()
	engine.AddPair(pub, point)
	engine.AddPairInv(g1.One(), sig)
	okay := engine.Check()
	if !okay {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func hashToCurveBls12381(message []byte, context string) (*bls12381.PointG2, error) {
	g2 := bls12381.NewG2()
	dst := []byte("BLS_SIG_" + context + "_XMD:SHA-256_SSWU_RO_NUL_")
	return g2.HashToCurve(message, dst)
}

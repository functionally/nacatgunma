package tgdh

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"

	bls12381 "github.com/kilic/bls12-381"
)

type Node struct {
	Private *bls12381.Fr
	Public  *bls12381.PointG1
	Left    *Node
	Right   *Node
}

func Strip(node *Node) *Node {
	return &Node{
		Private: nil,
		Public:  node.Public,
		Left:    node.Left,
		Right:   node.Right,
	}
}

func GenerateLeaf() (*Node, error) {
	pri, err := bls12381.NewFr().Rand(rand.Reader)
	if err != nil {
		return nil, err
	}
	return Leaf(pri), nil
}

func makeNode(pri *bls12381.Fr, left *Node, right *Node) *Node {
	g1 := bls12381.NewG1()
	return &Node{
		Private: pri,
		Public:  g1.MulScalar(g1.New(), g1.One(), pri),
		Left:    left,
		Right:   right,
	}
}

func Leaf(pri *bls12381.Fr) *Node {
	return makeNode(pri, nil, nil)
}

func Join(left *Node, right *Node) (*Node, error) {
	g1 := bls12381.NewG1()
	prod := g1.New()
	if left.Private != nil {
		g1.MulScalar(prod, right.Public, left.Private)
	} else if right.Private != nil {
		g1.MulScalar(prod, left.Public, right.Private)
	} else {
		return nil, fmt.Errorf("one child must have a private key")
	}
	prodHash := sha512.Sum512(g1.ToCompressed(prod))
	pri := bls12381.NewFr().FromBytes(prodHash[:])
	return makeNode(pri, left, right), nil
}

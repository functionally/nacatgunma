package tgdh

import (
	bls12381 "github.com/kilic/bls12-381"
)

type Node struct {
	Private *bls12381.Fr
	Public  *bls12381.PointG1
	Left    *Node
	Right   *Node
}

func Leaf(pri *bls12381.Fr, pub *bls12381.PointG1) *Node {
	return &Node{
		Private: pri,
		Public:  pub,
		Left:    nil,
		Right:   nil,
	}
}

/*
func Join(left *Node, right *Node) (*Node, error) {
	g1 := bls12381.NewG1()
	if left.Private != nil {

	}
}
*/

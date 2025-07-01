package tgdh

import (
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestGenerateLeaf(t *testing.T) {
	_, err := GenerateLeaf()
	if err != nil {
		t.Error(err)
	}
}

func TestRoot(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)
	leaves := make([]*Node, 8)
	for i := range leaves {
		leaf, err := GenerateLeaf()
		if err != nil {
			t.Error(err)
		}
		leaves[i] = leaf
	}
	join := func(flag bool, left *Node, right *Node) *Node {
		if flag {
			node, err := Join(left, right.Strip())
			if err != nil {
				t.Error(err)
			}
			return node
		} else {
			node, err := Join(left.Strip(), right)
			if err != nil {
				t.Error(err)
			}
			return node
		}
	}
	compute := func(flags [7]bool) *Node {
		ab := join(flags[0], leaves[0], leaves[1])
		cd := join(flags[1], leaves[2], leaves[3])
		ef := join(flags[2], leaves[4], leaves[5])
		gh := join(flags[3], leaves[6], leaves[7])
		abcd := join(flags[4], ab, cd)
		efgh := join(flags[5], ef, gh)
		return join(flags[6], abcd, efgh)
	}
	properties.Property("roots match", prop.ForAll(
		func(flags0 [7]bool, flags1 [7]bool) bool {
			node0 := compute(flags0)
			node1 := compute(flags1)
			return node0.Public == node1.Public && *node0.Private == *node1.Private
		},
		gen.ArrayOfN(7, gen.Bool(), reflect.TypeOf(true)),
		gen.ArrayOfN(7, gen.Bool(), reflect.TypeOf(true)),
	))
	properties.TestingRun(t)
}

func TestFindPath(t *testing.T) {
	A, _ := GenerateLeaf()
	B, _ := GenerateLeaf()
	AB, _ := Join(A, B.Strip())
	path, err := FindPath(A, AB.DeepStrip())
	if err != nil {
		t.Error(err)
	}
	if len(path) != 2 || path[0].Public != A.Public || path[1].Public != AB.Public {
		t.Error("incorrect path")
	}
}

func TestRecompute(t *testing.T) {
	A, _ := GenerateLeaf()
	B, _ := GenerateLeaf()
	AB, _ := Join(A, B.Strip())
	root, err := DerivePrivates(A, AB.DeepStrip())
	if err != nil {
		t.Error(err)
	}
	if *AB.Private != *root.Private {
		t.Error("incorrect recomputed private key")
	}
}

func TestJson(t *testing.T) {
	A, _ := GenerateLeaf()
	B, _ := GenerateLeaf()
	AB, _ := Join(A, B.Strip())
	root, err := DerivePrivates(A, AB.DeepStrip())
	if err != nil {
		t.Error(err)
	}
	j, err := root.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	root1, err := UnmarshalJSON(j)
	if err != nil {
		t.Error(err)
	}
	if !Equal(root1, root) {
		t.Error("deserialization does not match")
	}
}

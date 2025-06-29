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
			node, err := Join(left, Strip(right))
			if err != nil {
				t.Error(err)
			}
			return node
		} else {
			node, err := Join(Strip(left), right)
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
			return *node0.Public == *node1.Public && *node0.Private == *node1.Private
		},
		gen.ArrayOfN(7, gen.Bool(), reflect.TypeOf(true)),
		gen.ArrayOfN(7, gen.Bool(), reflect.TypeOf(true)),
	))
	properties.TestingRun(t)
}

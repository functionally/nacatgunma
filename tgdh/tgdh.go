package tgdh

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	bls12381 "github.com/kilic/bls12-381"
	"github.com/multiformats/go-multibase"
	"golang.org/x/crypto/hkdf"
)

type Node struct {
	Private *bls12381.Fr
	Public  bls12381.PointG1
	Left    *Node
	Right   *Node
}

func (root *Node) Strip() *Node {
	return &Node{
		Private: nil,
		Public:  root.Public,
		Left:    root.Left,
		Right:   root.Right,
	}
}

func (root *Node) DeepStrip() *Node {
	var left *Node
	if root.Left != nil {
		left = root.Left.DeepStrip()
	}
	var right *Node
	if root.Right != nil {
		right = root.Right.DeepStrip()
	}
	return &Node{
		Private: nil,
		Public:  root.Public,
		Left:    left,
		Right:   right,
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
		Public:  *g1.MulScalar(g1.New(), g1.One(), pri),
		Left:    left,
		Right:   right,
	}
}

func Leaf(pri *bls12381.Fr) *Node {
	return makeNode(pri, nil, nil)
}

func hashToFr(secret *bls12381.PointG1, salt []byte, info []byte) (*bls12381.Fr, error) {
	hashReader := hkdf.New(sha256.New, bls12381.NewG1().ToCompressed(secret), salt, info)
	hashBytes := make([]byte, 32)
	_, err := io.ReadFull(hashReader, hashBytes)
	if err != nil {
		return nil, err
	}
	return bls12381.NewFr().FromBytes(hashBytes), nil
}

func Join(left *Node, right *Node) (*Node, error) {
	g1 := bls12381.NewG1()
	prod := g1.New()
	if left.Private != nil {
		g1.MulScalar(prod, &right.Public, left.Private)
	} else if right.Private != nil {
		g1.MulScalar(prod, &left.Public, right.Private)
	} else {
		return nil, fmt.Errorf("one child must have a private key")
	}
	pri, err := hashToFr(prod, nil, []byte("nacatgunma-tgdh-bls12381g1"))
	if err != nil {
		return nil, err
	}
	return makeNode(pri, left, right), nil
}

func (root *Node) Remove(leaf *Node) (*Node, error) {
	if root.Public == leaf.Public {
		return nil, fmt.Errorf("cannot remove root")
	}
	path, err := FindPath(leaf, root)
	if err != nil {
		return nil, err
	}
	leaf = path[0]
	roots := []*Node{}
	for _, node := range path[1:] {
		if node.Left == leaf && node.Right != nil {
			roots = append(roots, node.Right)
		} else if node.Right == leaf && node.Left != nil {
			roots = append(roots, node.Left)
		}
		leaf = node
	}
	var start *Node
	for _, node := range roots {
		if node.Private != nil {
			start = node
			break
		}
	}
	if start == nil {
		return nil, fmt.Errorf("root does not contain any private keys")
	}
	root1 := start
	for _, node := range roots {
		if node == start {
			continue
		}
		root1, err = Join(root1, node)
		if err != nil {
			return nil, err
		}
	}
	return root1, nil
}

func (root *Node) Did() string {
	pub := bls12381.NewG1().ToCompressed(&root.Public)
	prefixedKey := append([]byte{0xEA, 0x01}, pub...)
	str, err := multibase.Encode(multibase.Base58BTC, prefixedKey)
	if err != nil {
		panic(err)
	}
	return "did:key:" + str
}

func visitPath(leaf *Node, root *Node, candidate *Node) ([]*Node, error) {
	if candidate == nil {
		return nil, fmt.Errorf("path not found")
	}
	if leaf.Public == candidate.Public {
		return []*Node{candidate}, nil
	}
	path, err := visitPath(leaf, root, candidate.Left)
	if err == nil {
		return append(path, candidate), nil
	}
	path, err = visitPath(leaf, root, candidate.Right)
	if err == nil {
		return append(path, candidate), nil
	}
	return nil, fmt.Errorf("path not found")
}

func FindPath(leaf *Node, root *Node) ([]*Node, error) {
	return visitPath(leaf, root, root)
}

func DerivePrivates(leaf *Node, root *Node) (*Node, error) {
	path, err := FindPath(leaf, root)
	if err != nil {
		return nil, err
	}
	root1 := leaf
	for _, node := range path[1:] {
		if node.Left == nil || node.Right == nil {
			return nil, fmt.Errorf("ill-formed tree")
		}
		if node.Left.Public == root1.Public {
			root1, err = Join(root1, node.Right)
			if err != nil {
				return nil, err
			}
		}
		if node.Right.Public == root1.Public {
			root1, err = Join(node.Left, root1)
			if err != nil {
				return nil, err
			}
		}
	}
	return root1, nil
}

func (root *Node) DeriveSeed(dst []byte, salt []byte, info []byte) error {
	if root.Private == nil {
		return fmt.Errorf("missing private key")
	}
	hashReader := hkdf.New(sha256.New, root.Private.ToBytes(), salt, info)
	_, err := io.ReadFull(hashReader, dst)
	return err
}

type jsonNode struct {
	Private string    `json:"private,omitempty"`
	Public  string    `json:"public"`
	Left    *jsonNode `json:"left,omitempty"`
	Right   *jsonNode `json:"right,omitempty"`
}

func (root *Node) toJSONNode() (*jsonNode, error) {
	if root == nil {
		return nil, nil
	}
	var privHex string
	if root.Private != nil {
		privBytes := root.Private.ToBytes()
		privHex = hex.EncodeToString(privBytes[:])
	}
	g1 := bls12381.NewG1()
	pubBytes := g1.ToCompressed(&root.Public)
	pubHex := hex.EncodeToString(pubBytes)
	left, err := root.Left.toJSONNode()
	if err != nil {
		return nil, err
	}
	right, err := root.Right.toJSONNode()
	if err != nil {
		return nil, err
	}
	return &jsonNode{
		Private: privHex,
		Public:  pubHex,
		Left:    left,
		Right:   right,
	}, nil
}

func (root *Node) MarshalJSON() ([]byte, error) {
	j, err := root.toJSONNode()
	if err != nil {
		return nil, err
	}
	return json.Marshal(j)
}

func (j *jsonNode) fromJSONNode() (*Node, error) {
	if j == nil {
		return nil, nil
	}
	var fr bls12381.Fr
	var pri *bls12381.Fr
	if j.Private != "" {
		privBytes, err := hex.DecodeString(j.Private)
		if err != nil {
			return nil, err
		}
		pri = fr.FromBytes(privBytes)
	} else {
		pri = nil
	}
	pubBytes, err := hex.DecodeString(j.Public)
	if err != nil {
		return nil, err
	}
	g1 := bls12381.NewG1()
	pub, err := g1.FromCompressed(pubBytes)
	if err != nil {
		return nil, err
	}
	left, err := j.Left.fromJSONNode()
	if err != nil {
		return nil, err
	}
	right, err := j.Right.fromJSONNode()
	if err != nil {
		return nil, err
	}
	return &Node{
		Private: pri,
		Public:  *pub,
		Left:    left,
		Right:   right,
	}, nil
}

func UnmarshalJSON(data []byte) (*Node, error) {
	var j jsonNode
	if err := json.Unmarshal(data, &j); err != nil {
		return nil, err
	}
	return j.fromJSONNode()
}

func Equal(x *Node, y *Node) bool {
	if x == nil && y == nil {
		return true
	}
	if (x == nil) != (y == nil) {
		return false
	}
	if (x.Private == nil) != (y.Private == nil) {
		return false
	}
	if x.Private != nil && !x.Private.Equal(y.Private) {
		return false
	}
	g1 := bls12381.NewG1()
	if !g1.Equal(&x.Public, &y.Public) {
		return false
	}
	return Equal(x.Left, y.Left) && Equal(x.Right, y.Right)
}

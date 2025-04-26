package ipfs

import (
	"bytes"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func EncodeToDAGCBOR(doc interface{}) ([]byte, error) {
	builder := basicnode.Prototype.Any.NewBuilder()
	err := assembleFromInterface(doc, builder)
	if err != nil {
		return nil, err
	}
	node := builder.Build()

	var buf bytes.Buffer
	err = dagcbor.Encode(node, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func assembleFromInterface(value interface{}, assembler ipld.NodeAssembler) error {
	switch v := value.(type) {
	case map[string]interface{}:
		ma, _ := assembler.BeginMap(int64(len(v)))
		for k, val := range v {
			ma.AssembleKey().AssignString(k)
			err := assembleFromInterface(val, ma.AssembleValue())
			if err != nil {
				return err
			}
		}
		return ma.Finish()
	case []interface{}:
		la, _ := assembler.BeginList(int64(len(v)))
		for _, item := range v {
			err := assembleFromInterface(item, la.AssembleValue())
			if err != nil {
				return err
			}
		}
		return la.Finish()
	case string:
		return assembler.AssignString(v)
	case float64:
		return assembler.AssignFloat(v)
	case bool:
		return assembler.AssignBool(v)
	case nil:
		return assembler.AssignNull()
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func CidV0(bytes []byte) (*cid.Cid, error) {
	format := cid.V0Builder{}
	id, err := format.Sum(bytes)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func DecodeFromDAGCBOR(data []byte) (interface{}, error) {
	builder := basicnode.Prototype.Any.NewBuilder()
	err := dagcbor.Decode(builder, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	node := builder.Build()
	return nodeToInterface(node)
}

func nodeToInterface(n ipld.Node) (interface{}, error) {
	switch n.Kind() {
	case ipld.Kind_Map:
		m := make(map[string]interface{})
		it := n.MapIterator()
		for !it.Done() {
			k, v, _ := it.Next()
			kString, err := k.AsString()
			if err != nil {
				return nil, err
			}
			i, err := nodeToInterface(v)
			if err != nil {
				return nil, err
			}
			m[kString] = i
		}
		return m, nil
	case ipld.Kind_List:
		var l []interface{}
		it := n.ListIterator()
		for !it.Done() {
			_, v, _ := it.Next()
			i, err := nodeToInterface(v)
			if err != nil {
				return nil, err
			}
			l = append(l, i)
		}
		return l, nil
	case ipld.Kind_String:
		return n.AsString()
	case ipld.Kind_Int:
		return n.AsInt()
	case ipld.Kind_Float:
		return n.AsFloat()
	case ipld.Kind_Bool:
		return n.AsBool()
	case ipld.Kind_Null:
		return nil, nil
	default:
		return nil, nil
	}
}

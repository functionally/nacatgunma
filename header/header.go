package header

import (
	"bytes"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	"github.com/functionally/nacatgunma/key"
)

type Header struct {
	Payload   Payload
	Issuer    string
	Signature []byte
}

func (header *Header) MakeNode() datamodel.Node {
	return fluent.MustBuildMap(basicnode.Prototype__Any{}, 3,
		func(assembler fluent.MapAssembler) {
			assembler.AssembleEntry("Payload").AssignNode(header.Payload.MakeNode())
			assembler.AssembleEntry("Issuer").AssignString(header.Issuer)
			assembler.AssembleEntry("Signature").AssignBytes(header.Signature)
		})
}

func (header *Header) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	err := dagcbor.Encode(header.MakeNode(), &buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decodeHeader(node ipld.Node) (*Header, error) {
	header := &Header{}
	payloadNode, err := node.LookupByString("Payload")
	if err != nil {
		return nil, err
	}
	payload, err := decodePayload(payloadNode)
	if err != nil {
		return nil, err
	}
	header.Payload = *payload
	if v, err := node.LookupByString("Issuer"); err == nil {
		s, err := v.AsString()
		if err != nil {
			return nil, err
		}
		header.Issuer = s
	}
	if v, err := node.LookupByString("Signature"); err == nil {
		b, err := v.AsBytes()
		if err != nil {
			return nil, err
		}
		header.Signature = b
	}
	return header, nil
}

func UnmarshalHeader(data []byte) (*Header, error) {
	nb := basicnode.Prototype__Any{}.NewBuilder()
	if err := dagcbor.Decode(nb, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	node := nb.Build()
	return decodeHeader(node)
}

func (header *Header) Verify() (bool, error) {
	bytes, err := header.Payload.Marshal()
	if err != nil {
		return false, err
	}
	err = key.Verify(header.Issuer, bytes, header.Signature, header.Issuer)
	return err == nil, err
}

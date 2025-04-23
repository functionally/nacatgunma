package achain

import (
	"crypto/ed25519"

	"github.com/fido-device-onboard/go-fdo/cbor"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

type Payload struct {
	Version   int64
	Accept    []cid.Cid
	Reject    []cid.Cid
	Body      cid.Cid
	Schema    string
	MediaType string
}

func (payload *Payload) MakeNode() datamodel.Node {
	return fluent.MustBuildMap(basicnode.Prototype__Any{}, 6,
		func(assembler fluent.MapAssembler) {
			assembler.AssembleEntry("Version").AssignInt(payload.Version)
			assembler.AssembleEntry("Accept").CreateList(2, func(la fluent.ListAssembler) {
				for _, accept := range payload.Accept {
					la.AssembleValue().AssignLink(cidlink.Link{Cid: accept})
				}
			})
			assembler.AssembleEntry("Reject").CreateList(1, func(la fluent.ListAssembler) {
				for _, reject := range payload.Reject {
					la.AssembleValue().AssignLink(cidlink.Link{Cid: reject})
				}
			})
			assembler.AssembleEntry("Body").AssignLink(cidlink.Link{Cid: payload.Body})
			assembler.AssembleEntry("Schema").AssignString(payload.Schema)
			assembler.AssembleEntry("MediaType").AssignString(payload.MediaType)
		})
}

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

func (header *Header) Marshal() (*cid.Cid, []byte, error) {
	bytes, err := cbor.Marshal(header)
	if err != nil {
		return nil, nil, err
	}
	format := cid.V0Builder{}
	id, err := format.Sum(bytes)
	if err != nil {
		return nil, nil, err
	}
	return &id, bytes, nil
}

func (payload *Payload) Sign(key *Key) (*Header, error) {
	bytes, err := cbor.Marshal(payload)
	if err != nil {
		return nil, err
	}
	s, err := key.Private.Sign(nil, bytes, &ed25519.Options{
		Context: key.Did,
	})
	if err != nil {
		return nil, err
	}
	return &Header{
		Payload:   *payload,
		Issuer:    key.Did,
		Signature: s,
	}, nil
}

func (header *Header) Verify() (bool, error) {
	bytes, err := cbor.Marshal(header.Payload)
	if err != nil {
		return false, err
	}
	pub, err := PublicKeyFromDid(header.Issuer)
	if err != nil {
		return false, err
	}
	err = ed25519.VerifyWithOptions(pub, bytes, header.Signature, &ed25519.Options{
		Context: header.Issuer,
	})
	return err == nil, err
}

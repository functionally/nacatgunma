package header

import (
	"bytes"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	"github.com/functionally/nacatgunma/key"
)

type Payload struct {
	Version   int64
	Accept    []cid.Cid
	Reject    []cid.Cid
	Body      cid.Cid
	SchemaUri string
	MediaType string
	Comment   string
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
			assembler.AssembleEntry("Schema").AssignString(payload.SchemaUri)
			assembler.AssembleEntry("MediaType").AssignString(payload.MediaType)
			assembler.AssembleEntry("Comment").AssignString(payload.Comment)
		})
}

func (payload *Payload) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	err := dagcbor.Encode(payload.MakeNode(), &buffer)
	return buffer.Bytes(), err
}

func (payload *Payload) Sign(k key.Key) (*Header, error) {
	bytes, err := payload.Marshal()
	if err != nil {
		return nil, err
	}
	did := key.Did(k)
	s, err := k.Sign(bytes, did)
	if err != nil {
		return nil, err
	}
	return &Header{
		Payload:   *payload,
		Issuer:    did,
		Signature: s,
	}, nil
}

func decodePayload(node ipld.Node) (*Payload, error) {
	payload := &Payload{}
	if v, err := node.LookupByString("Version"); err == nil {
		if ver, err := v.AsInt(); err == nil {
			payload.Version = ver
		} else {
			return nil, err
		}
	}
	if v, err := node.LookupByString("Accept"); err == nil {
		iter := v.ListIterator()
		for !iter.Done() {
			_, val, _ := iter.Next()
			lnk, err := val.AsLink()
			if err != nil {
				return nil, err
			}
			if cl, ok := lnk.(cidlink.Link); ok {
				payload.Accept = append(payload.Accept, cl.Cid)
			}
		}
	}
	if v, err := node.LookupByString("Reject"); err == nil {
		iter := v.ListIterator()
		for !iter.Done() {
			_, val, _ := iter.Next()
			lnk, err := val.AsLink()
			if err != nil {
				return nil, err
			}
			if cl, ok := lnk.(cidlink.Link); ok {
				payload.Reject = append(payload.Reject, cl.Cid)
			}
		}
	}
	if v, err := node.LookupByString("Body"); err == nil {
		lnk, err := v.AsLink()
		if err != nil {
			return nil, err
		}
		if cl, ok := lnk.(cidlink.Link); ok {
			payload.Body = cl.Cid
		}
	}
	if v, err := node.LookupByString("Schema"); err == nil {
		s, err := v.AsString()
		if err != nil {
			return nil, err
		}
		payload.SchemaUri = s
	}
	if v, err := node.LookupByString("MediaType"); err == nil {
		s, err := v.AsString()
		if err != nil {
			return nil, err
		}
		payload.MediaType = s
	}
	if v, err := node.LookupByString("Comment"); err == nil {
		s, err := v.AsString()
		if err != nil {
			return nil, err
		}
		payload.Comment = s
	}
	return payload, nil
}

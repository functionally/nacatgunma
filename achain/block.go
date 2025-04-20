package achain

import (
	"crypto"
	"crypto/x509"

	"github.com/fido-device-onboard/go-fdo/cbor"
	"github.com/fido-device-onboard/go-fdo/cose"
	"github.com/ipfs/go-cid"
)

type Header struct {
	Version   int16
	Accept    []cid.Cid
	Reject    []cid.Cid
	Body      cid.Cid
	Schema    string
	MediaType string
}

type SignedHeader struct {
	IssuerDER []byte
	Signature cose.Sign1[Header, []byte]
}

func (sh *SignedHeader) Marshal() (*cid.Cid, []byte, error) {
	bytes, err := cbor.Marshal(sh)
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

func Sign(key crypto.Signer, header *Header) (*SignedHeader, error) {
	pubDER, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return nil, err
	}
	s := cose.Sign1[Header, []byte]{
		Payload: cbor.NewByteWrap(*header),
	}
	s.Sign(key, nil, pubDER, nil)
	sh := SignedHeader{
		IssuerDER: pubDER,
		Signature: s,
	}
	return &sh, nil
}

func (sh *SignedHeader) Verify() (bool, error) {
	pubDER := sh.IssuerDER
	pub, err := x509.ParsePKIXPublicKey(pubDER)
	if err != nil {
		return false, err
	}
	return sh.Signature.Verify(pub, nil, pubDER)
}

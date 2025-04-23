package main

import (
	"fmt"
	"os"

	"github.com/functionally/achain/achain"
	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
)

func main() {
	k, err := achain.GenerateKey()
	if err != nil {
		panic(err)
	}
	err = achain.WritePrivateKey("private.pem", k)
	if err != nil {
		panic(err)
	}
	k1, err := achain.ReadPrivateKey("private.pem")
	if err != nil {
		panic(err)
	}
	fmt.Println(k.Equal(k1))
	body, err := cid.Parse("QmYGc9ncJhbejE4TLbP3NMX5fjHvZioCTFknD6HbnjBvpm")
	if err != nil {
		panic(err)
	}
	plain := achain.Header{
		Version:   1,
		Accept:    []cid.Cid{},
		Reject:    []cid.Cid{},
		Body:      body,
		Schema:    "",
		MediaType: "application/octet-stream",
	}
	sig, err := achain.Sign(k, &plain)
	if err != nil {
		panic(err)
	}
	zid, zbytes, err := sig.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Println(zid)
	fmt.Println(len(zbytes))
	sigBytes, err := cbor.Marshal(sig)
	if err != nil {
		panic(err)
	}
	os.WriteFile("tmp.cbor", sigBytes, 0644)
	sigBytes1, err := os.ReadFile("tmp.cbor")
	if err != nil {
		panic(err)
	}
	sig1 := achain.SignedHeader{}
	err = cbor.Unmarshal(sigBytes1, &sig1)
	if err != nil {
		panic(err)
	}
	okay, err := sig1.Verify()
	if err != nil {
		panic(err)
	}
	fmt.Println(okay)
}

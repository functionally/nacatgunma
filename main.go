package main

import (
	"fmt"
	"os"

	"github.com/functionally/achain/achain"
	"github.com/ipfs/go-cid"
)

func main() {
	k, err := achain.GenerateKey()
	if err != nil {
		panic(err)
	}
	err = k.WritePrivateKey("private.pem")
	if err != nil {
		panic(err)
	}
	_, err = achain.ReadPrivateKey("private.pem")
	if err != nil {
		panic(err)
	}
	cidBuilder := cid.V1Builder{Codec: cid.Raw, MhType: 0x12}
	cid1, _ := cidBuilder.Sum([]byte("accept1"))
	cid2, _ := cidBuilder.Sum([]byte("accept2"))
	cid3, _ := cidBuilder.Sum([]byte("reject1"))
	bodyCid, _ := cidBuilder.Sum([]byte("body content"))
	if err != nil {
		panic(err)
	}
	plain := achain.Payload{
		Version:   1,
		Accept:    []cid.Cid{cid1, cid2},
		Reject:    []cid.Cid{cid3},
		Body:      bodyCid,
		Schema:    "",
		MediaType: "application/octet-stream",
	}
	sig, err := plain.Sign(k)
	if err != nil {
		panic(err)
	}
	zid, zbytes, err := sig.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Println(zid)
	os.WriteFile("tmp.cbor", zbytes, 0644)
	sigBytes1, err := os.ReadFile("tmp.cbor")
	if err != nil {
		panic(err)
	}
	sig1, err := achain.UnmarshalHeader(sigBytes1)
	if err != nil {
		panic(err)
	}
	okay, err := sig1.Verify()
	if err != nil {
		panic(err)
	}
	fmt.Println(okay)
}

package ledger

import (
	"testing"

	"github.com/functionally/nacatgunma/header"
	"github.com/functionally/nacatgunma/ipfs"
	"github.com/functionally/nacatgunma/key"
	"github.com/ipfs/go-cid"
)

func makeHeader(accept []cid.Cid, reject []cid.Cid) (cid.Cid, *header.Header) {
	body, _ := cid.Parse("bafyreih3lbpdqibixvdr3twiqwqrx3tgxbcwuooaq6ieyxzjzkw5zoxb3m")
	payload := header.Payload{
		Version:   1,
		Body:      body,
		Accept:    accept,
		Reject:    reject,
		SchemaUri: "https://w3c.github.io/json-ld-cbor/",
		MediaType: "application/vnd.ipld.dag-cbor",
		Comment:   "",
	}
	ky, _ := key.GenerateKey()
	hdr, _ := payload.Sign(ky)
	hdrBytes, _ := hdr.Marshal()
	hdrCid, _ := ipfs.CidV1(hdrBytes)
	return *hdrCid, hdr
}

func empty() map[cid.Cid]header.Header {
	return make(map[cid.Cid]header.Header)
}

var c0, h0 = makeHeader([]cid.Cid{}, []cid.Cid{})

func makeL0() Ledger {
	hs := empty()
	hs[c0] = *h0
	return Ledger{
		Tip:     c0,
		Headers: hs,
	}
}

func TestPrune(t *testing.T) {

	t.Run("Genesis", func(t *testing.T) {
		l0 := makeL0()
		rejects := l0.Prune()
		if len(rejects) != 0 {
			t.Error("Incorrect pruning")
		}
	})

}

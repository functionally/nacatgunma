package ledger

import (
	"log"
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

func assertEqual(xs []cid.Cid, ys []cid.Cid) bool {
	if len(xs) != len(ys) {
		return false
	}
	xx := make(map[cid.Cid]bool)
	for _, x := range xs {
		xx[x] = true
	}
	for _, y := range ys {
		_, ok := xx[y]
		if !ok {
			return false
		}
		delete(xx, y)
	}
	return true
}

var c0, h0 = makeHeader([]cid.Cid{}, []cid.Cid{})

var c1, h1 = makeHeader([]cid.Cid{c0}, []cid.Cid{})

var c2, h2 = makeHeader([]cid.Cid{c1}, []cid.Cid{})

var c3, h3 = makeHeader([]cid.Cid{c2}, []cid.Cid{c0})

var c4, h4 = makeHeader([]cid.Cid{c2}, []cid.Cid{c1})

var c5, h5 = makeHeader([]cid.Cid{c2}, []cid.Cid{c2})

var c6, h6 = makeHeader([]cid.Cid{c1}, []cid.Cid{})

var c7, h7 = makeHeader([]cid.Cid{c2, c6}, []cid.Cid{})

var c8, h8 = makeHeader([]cid.Cid{c7}, []cid.Cid{c6})

var c9, h9 = makeHeader([]cid.Cid{c1}, []cid.Cid{})

var c10, h10 = makeHeader([]cid.Cid{c9}, []cid.Cid{})

var c11, h11 = makeHeader([]cid.Cid{c8, c10}, []cid.Cid{})

var c12, h12 = makeHeader([]cid.Cid{c11}, []cid.Cid{c2})

var c13, h13 = makeHeader([]cid.Cid{c7, c12}, []cid.Cid{})

var c14, h14 = makeHeader([]cid.Cid{c13}, []cid.Cid{c9})

func TestPrune(t *testing.T) {

	t.Run("Genesis", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		le := Ledger{
			Tip:     c0,
			Headers: hs,
		}
		rejects := le.Prune()
		if len(rejects) != 0 {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Two blocks linear", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		le := Ledger{
			Tip:     c1,
			Headers: hs,
		}
		rejects := le.Prune()
		if len(rejects) != 0 {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Three blocks linear", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		le := Ledger{
			Tip:     c2,
			Headers: hs,
		}
		rejects := le.Prune()
		if len(rejects) != 0 {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Reject genesis", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c3] = *h3
		le := Ledger{
			Tip:     c3,
			Headers: hs,
		}
		rejects := le.Prune()
		expected := []cid.Cid{c0, c1, c2}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Reject first", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c4] = *h4
		le := Ledger{
			Tip:     c4,
			Headers: hs,
		}
		rejects := le.Prune()
		expected := []cid.Cid{c1, c2}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Reject second", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c5] = *h5
		le := Ledger{
			Tip:     c5,
			Headers: hs,
		}
		rejects := le.Prune()
		expected := []cid.Cid{c2}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Parallel path", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c6] = *h6
		hs[c7] = *h7
		le := Ledger{
			Tip:     c7,
			Headers: hs,
		}
		rejects := le.Prune()
		if len(rejects) != 0 {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Reject whole parallel path", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c6] = *h6
		hs[c7] = *h7
		hs[c8] = *h8
		le := Ledger{
			Tip:     c8,
			Headers: hs,
		}
		rejects := le.Prune()
		expected := []cid.Cid{c6}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Complex 1", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c6] = *h6
		hs[c7] = *h7
		hs[c8] = *h8
		hs[c9] = *h9
		hs[c10] = *h10
		hs[c11] = *h11
		le := Ledger{
			Tip:     c11,
			Headers: hs,
		}
		rejects := le.Prune()
		expected := []cid.Cid{c6}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Complex 2", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c6] = *h6
		hs[c7] = *h7
		hs[c8] = *h8
		hs[c9] = *h9
		hs[c10] = *h10
		hs[c11] = *h11
		hs[c12] = *h12
		le := Ledger{
			Tip:     c12,
			Headers: hs,
		}
		rejects := le.Prune()
		expected := []cid.Cid{c2, c6, c7, c8}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})
	t.Run("Complex 3", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c6] = *h6
		hs[c7] = *h7
		hs[c8] = *h8
		hs[c9] = *h9
		hs[c10] = *h10
		hs[c11] = *h11
		hs[c12] = *h12
		hs[c13] = *h13
		le := Ledger{
			Tip:     c13,
			Headers: hs,
		}
		rejects := le.Prune()
		log.Println(rejects)
		expected := []cid.Cid{c2, c6, c7, c8}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Complex 4", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c6] = *h6
		hs[c7] = *h7
		hs[c8] = *h8
		hs[c9] = *h9
		hs[c10] = *h10
		hs[c11] = *h11
		hs[c12] = *h12
		hs[c13] = *h13
		hs[c14] = *h14
		le := Ledger{
			Tip:     c14,
			Headers: hs,
		}
		rejects := le.Prune()
		log.Println(rejects)
		expected := []cid.Cid{c2, c6, c7, c8, c9, c10, c11, c12, c13}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

}

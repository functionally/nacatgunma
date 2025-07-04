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
		SchemaURI: "https://w3c.github.io/json-ld-cbor/",
		MediaType: "application/vnd.ipld.dag-cbor",
		Comment:   "",
	}
	ky, _ := key.GenerateKey(key.Ed25519)
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

	t.Run("Ex0 Genesis", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		le := Ledger{
			Tip:     c0,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{c0}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex1 Two blocks linear", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		le := Ledger{
			Tip:     c1,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{c0, c1}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex2 Three blocks linear", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		le := Ledger{
			Tip:     c2,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{c0, c1, c2}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex3 Reject genesis", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c3] = *h3
		le := Ledger{
			Tip:     c3,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{c1, c2, c3}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex4 Reject first", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c4] = *h4
		le := Ledger{
			Tip:     c4,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{c2, c4}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex5 Reject second", func(t *testing.T) {
		hs := empty()
		hs[c0] = *h0
		hs[c1] = *h1
		hs[c2] = *h2
		hs[c5] = *h5
		le := Ledger{
			Tip:     c5,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{c5}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex6 Parallel path", func(t *testing.T) {
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
		visible := le.Visible()
		expected := []cid.Cid{c0, c1, c2, c6, c7}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex7 Reject whole parallel path", func(t *testing.T) {
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
		visible := le.Visible()
		expected := []cid.Cid{c0, c1, c2, c7, c8}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 8: Complex 1", func(t *testing.T) {
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
		visible := le.Visible()
		expected := []cid.Cid{c0, c1, c2, c7, c8, c9, c10, c11}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex9 Complex 2", func(t *testing.T) {
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
		visible := le.Visible()
		expected := []cid.Cid{c0, c1, c7, c8, c9, c10, c11, c12}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex10 Complex 3", func(t *testing.T) {
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
		visible := le.Visible()
		expected := []cid.Cid{c0, c1, c2, c6, c7, c8, c9, c10, c11, c12, c13}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex11 Complex 4", func(t *testing.T) {
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
		visible := le.Visible()
		expected := []cid.Cid{c0, c1, c2, c6, c7, c8, c10, c11, c12, c13, c14}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 12", func(t *testing.T) {
		cA, hA := makeHeader([]cid.Cid{}, []cid.Cid{})
		cB, hB := makeHeader([]cid.Cid{cA}, []cid.Cid{})
		cC, hC := makeHeader([]cid.Cid{cA}, []cid.Cid{})
		cD, hD := makeHeader([]cid.Cid{cC}, []cid.Cid{})
		cE, hE := makeHeader([]cid.Cid{cD}, []cid.Cid{})
		cF, hF := makeHeader([]cid.Cid{cB, cE}, []cid.Cid{cD})
		cG, hG := makeHeader([]cid.Cid{cF}, []cid.Cid{})
		hs := empty()
		hs[cA] = *hA
		hs[cB] = *hB
		hs[cC] = *hC
		hs[cD] = *hD
		hs[cE] = *hE
		hs[cF] = *hF
		hs[cG] = *hG
		le := Ledger{
			Tip:     cG,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{cA, cB, cE, cF, cG}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 13", func(t *testing.T) {
		cA, hA := makeHeader([]cid.Cid{}, []cid.Cid{})
		cC, hC := makeHeader([]cid.Cid{cA}, []cid.Cid{})
		cB, hB := makeHeader([]cid.Cid{cA}, []cid.Cid{cC})
		cD, hD := makeHeader([]cid.Cid{cC}, []cid.Cid{})
		cE, hE := makeHeader([]cid.Cid{cD}, []cid.Cid{})
		cF, hF := makeHeader([]cid.Cid{cB, cE}, []cid.Cid{})
		cG, hG := makeHeader([]cid.Cid{cF}, []cid.Cid{})
		hs := empty()
		hs[cA] = *hA
		hs[cB] = *hB
		hs[cC] = *hC
		hs[cD] = *hD
		hs[cE] = *hE
		hs[cF] = *hF
		hs[cG] = *hG
		le := Ledger{
			Tip:     cG,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{cA, cB, cC, cD, cE, cF, cG}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 14", func(t *testing.T) {
		cR, hR := makeHeader([]cid.Cid{}, []cid.Cid{})
		cB, hB := makeHeader([]cid.Cid{cR}, []cid.Cid{cR})
		cA, hA := makeHeader([]cid.Cid{cB}, []cid.Cid{})
		cY, hY := makeHeader([]cid.Cid{cB}, []cid.Cid{})
		cX, hX := makeHeader([]cid.Cid{cY}, []cid.Cid{})
		cT, hT := makeHeader([]cid.Cid{cA, cX}, []cid.Cid{})
		hs := empty()
		hs[cR] = *hR
		hs[cB] = *hB
		hs[cA] = *hA
		hs[cY] = *hY
		hs[cX] = *hX
		hs[cT] = *hT
		le := Ledger{
			Tip:     cT,
			Headers: hs,
		}
		visible := le.Visible()
		expected := []cid.Cid{cA, cB, cT, cX, cY}
		if !assertEqual(visible, expected) {
			t.Error("Incorrect pruning")
		}
	})
}

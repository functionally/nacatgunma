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

	t.Run("Ex 0: Genesis", func(t *testing.T) {
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

	t.Run("Ex 1: Two blocks linear", func(t *testing.T) {
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

	t.Run("Ex 2: Three blocks linear", func(t *testing.T) {
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

	t.Run("Ex 3: Reject genesis", func(t *testing.T) {
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
		expected := []cid.Cid{c0}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 4: Reject first", func(t *testing.T) {
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
		expected := []cid.Cid{c0, c1}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 5: Reject second", func(t *testing.T) {
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
		expected := []cid.Cid{c0, c1, c2}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 6: Parallel path", func(t *testing.T) {
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

	t.Run("Ex 7: Reject whole parallel path", func(t *testing.T) {
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
		rejects := le.Prune()
		expected := []cid.Cid{c6}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 9: Complex 2", func(t *testing.T) {
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
		expected := []cid.Cid{c2, c6}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 10: Complex 3", func(t *testing.T) {
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
		if len(rejects) > 0 {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 11: Complex 4", func(t *testing.T) {
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
		expected := []cid.Cid{c9}
		if !assertEqual(rejects, expected) {
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
		hs[cE] = *hD
		hs[cE] = *hE
		hs[cF] = *hF
		hs[cG] = *hG
		le := Ledger{
			Tip:     cG,
			Headers: hs,
		}
		rejects := le.Prune()
		log.Printf("A = %v\n", cA)
		log.Printf("B = %v\n", cB)
		log.Printf("C = %v\n", cC)
		log.Printf("D = %v\n", cD)
		log.Printf("E = %v\n", cE)
		log.Printf("F = %v\n", cF)
		log.Printf("G = %v\n", cG)
		log.Println(rejects)
		log.Println()
		expected := []cid.Cid{cC, cD}
		if !assertEqual(rejects, expected) {
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
		hs[cE] = *hD
		hs[cE] = *hE
		hs[cF] = *hF
		hs[cG] = *hG
		le := Ledger{
			Tip:     cG,
			Headers: hs,
		}
		rejects := le.Prune()
		log.Printf("A = %v\n", cA)
		log.Printf("B = %v\n", cB)
		log.Printf("C = %v\n", cC)
		log.Printf("D = %v\n", cD)
		log.Printf("E = %v\n", cE)
		log.Printf("F = %v\n", cF)
		log.Printf("G = %v\n", cG)
		log.Println(rejects)
		log.Println()
		if len(rejects) > 0 {
			t.Error("Incorrect pruning")
		}
	})

	t.Run("Ex 14", func(t *testing.T) {
		cR, hR := makeHeader([]cid.Cid{}, []cid.Cid{})
		cB, hB := makeHeader([]cid.Cid{cR}, []cid.Cid{cR})
		cA, hA := makeHeader([]cid.Cid{cB}, []cid.Cid{})
		cY, hY := makeHeader([]cid.Cid{cB}, []cid.Cid{})
		cX, hX := makeHeader([]cid.Cid{cA}, []cid.Cid{})
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
		rejects := le.Prune()
		log.Printf("R = %v\n", cR)
		log.Printf("B = %v\n", cB)
		log.Printf("A = %v\n", cA)
		log.Printf("Y = %v\n", cY)
		log.Printf("X = %v\n", cX)
		log.Printf("T = %v\n", cT)
		log.Println(rejects)
		log.Println()
		expected := []cid.Cid{cR}
		if !assertEqual(rejects, expected) {
			t.Error("Incorrect pruning")
		}
	})
}

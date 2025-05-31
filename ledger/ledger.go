package ledger

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-cid"

	"github.com/functionally/nacatgunma/header"
)

type Ledger struct {
	Tip     cid.Cid
	Headers map[cid.Cid]header.Header
}

func ReadLedger(tip string, headerDir string) (*Ledger, error) {
	tipCid, err := cid.Parse(tip)
	if err != nil {
		return nil, err
	}
	ledger :=
		Ledger{
			Tip:     tipCid,
			Headers: make(map[cid.Cid]header.Header),
		}
	err = ledger.fillLedger(tipCid, headerDir)
	if err != nil {
		return nil, err
	}
	return &ledger, nil
}

func (ledger *Ledger) fillLedger(headerCid cid.Cid, headerDir string) error {
	headerFile := filepath.Join(headerDir, headerCid.String())
	headerBytes, err := os.ReadFile(headerFile)
	if err != nil {
		return err
	}
	hdr, err := header.UnmarshalHeader(headerBytes)
	if err != nil {
		return err
	}
	ledger.Headers[headerCid] = *hdr
	for _, acceptCid := range hdr.Payload.Accept {
		err := ledger.fillLedger(acceptCid, headerDir)
		if err != nil {
			return err
		}
	}
	for _, rejectCid := range hdr.Payload.Reject {
		err := ledger.fillLedger(rejectCid, headerDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ledger *Ledger) Bodies() []cid.Cid {
	var bodies []cid.Cid
	for _, hdr := range ledger.Headers {
		bodies = append(bodies, hdr.Payload.Body)
	}
	return bodies
}

func (ledger *Ledger) Prune() []cid.Cid {
	colors := ledger.Prunable()
	var rejected []cid.Cid
	for cid := range colors {
		rejected = append(rejected, cid)
	}
	for _, cid := range rejected {
		delete(ledger.Headers, cid)
	}
	return rejected
}

func (ledger *Ledger) Prunable() map[cid.Cid]bool {
	colors := ledger.colorRejected()
	found := true
	for found {
		found = ledger.colorAdjacentRejected(colors)
	}
	return colors
}

func (ledger *Ledger) colorRejected() map[cid.Cid]bool {
	colors := make(map[cid.Cid]bool)
	for _, header := range ledger.Headers {
		for _, rejectCid := range header.Payload.Reject {
			colors[rejectCid] = true
		}
	}
	return colors
}

func (ledger *Ledger) colorAdjacentRejected(colors map[cid.Cid]bool) bool {
	found := false
	for headerCid, header := range ledger.Headers {
		all := len(header.Payload.Accept) > 0
		for _, acceptCid := range header.Payload.Accept {
			all = all && colors[acceptCid]
		}
		if all {
			_, ok := colors[headerCid]
			found = found || !ok
			colors[headerCid] = true
		}
	}
	return found
}

func (ledger *Ledger) WriteLedgerTurtle(outputFile string, prunable map[cid.Cid]bool) error {
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	defer f.Close()
	err = writePrefixTurtle(f)
	if err != nil {
		return nil
	}
	for headerCid, hdr := range ledger.Headers {
		_, prune := prunable[headerCid]
		err = writeHeaderTurtle(f, headerCid, &hdr, prune, ledger.Tip)
		if err != nil {
			return nil
		}
	}
	return nil
}

func writePrefixTurtle(f *os.File) error {
	_, err := f.WriteString(`@prefix :         <urn:uuid:e7c8a7a8-eecb-4474-af36-a0ca474a2af5#> .
@prefix dcterms:  <http://purl.org/dc/terms/> .
@prefix rdfs:     <http://www.w3.org/2000/01/rdf-schema#> .
@prefix xsd:      <http://www.w3.org/2001/XMLSchema#> .
@prefix cid:      <ipfs://> .
`)
	if err != nil {
		return err
	}
	return nil
}

func writeHeaderTurtle(f *os.File, hdrCid cid.Cid, hdr *header.Header, prune bool, tipCid cid.Cid) error {
	_, err := f.WriteString(fmt.Sprintf(
		`
cid:%v a :Header
; dcterms:creator <%v>
; :signature "%v"^^xsd:base64Binary
; :payload [
    dcterms:hasVersion "%v"^^xsd:long
  ; dcterms:conformsTo <%v>
  ; dcterms:format "%v"
  ; rdfs:comment "%v"
  ; :body cid:%v
`,
		hdrCid,
		hdr.Issuer,
		base64.StdEncoding.EncodeToString(hdr.Signature),
		hdr.Payload.Version,
		hdr.Payload.SchemaUri,
		hdr.Payload.MediaType,
		strings.ReplaceAll(hdr.Payload.Comment, `"`, `\"`),
		hdr.Payload.Body,
	))
	if err != nil {
		return err
	}
	for _, acceptCid := range hdr.Payload.Accept {
		_, err = f.WriteString(fmt.Sprintf("  ; :accept cid:%v\n", acceptCid.String()))
		if err != nil {
			return err
		}
	}
	for _, rejectCid := range hdr.Payload.Reject {
		_, err = f.WriteString(fmt.Sprintf("  ; :reject cid:%v\n", rejectCid.String()))
		if err != nil {
			return err
		}
	}
	if prune {
		_, err = f.WriteString(fmt.Sprintf("  ]\n; :rejectedBy cid:%v\n.", tipCid))
	} else {
		_, err = f.WriteString("  ]\n.")
	}
	if err != nil {
		return err
	}
	return nil
}

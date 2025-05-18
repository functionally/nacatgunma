package cardano

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/ipfs/go-cid"
)

type Datum struct {
	Script         bool
	CredentialHash common.Blake2b224
	HeaderCid      cid.Cid
}

func NewDatum(script bool, credential string, header string) (*Datum, error) {
	var datum Datum
	datum.Script = script
	credentialBytes, err := hex.DecodeString(credential)
	if err != nil {
		return nil, err
	}
	datum.CredentialHash = common.Blake2b224(common.NewBlake2b224(credentialBytes))
	datum.HeaderCid, err = cid.Parse(header)
	if err != nil {
		return nil, err
	}
	return &datum, nil
}

func (datum *Datum) ToJSON() ([]byte, error) {
	var buf bytes.Buffer
	var script int
	if datum.Script {
		script = 1
	} else {
		script = 0
	}
	_, err := fmt.Fprintf(
		&buf,
		"{\"list\" : [{\"constructor\" : %v, \"fields\" : [{\"bytes\" : \"%x\"}]}, {\"bytes\" : \"%x\"}]}",
		script,
		datum.CredentialHash.Bytes(),
		datum.HeaderCid.Bytes(),
	)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func RedeemerJSON(metadataKey uint) ([]byte, error) {
	var buf bytes.Buffer
	_, err := fmt.Fprintf(
		&buf,
		"{\"int\" : %v}",
		metadataKey,
	)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MetadataJSON(metadataKey uint, blockchain string, header string) ([]byte, error) {
	var buf bytes.Buffer
	_, err := fmt.Fprintf(
		&buf,
		"{\"%v\" : {\"blockchain\" : \"%v\", \"header\" : {\"ipfs\" : \"%v\"}}}",
		metadataKey,
		blockchain,
		header,
	)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

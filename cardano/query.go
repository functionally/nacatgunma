package cardano

import (
	"fmt"

	"github.com/blinklabs-io/gouroboros/cbor"
	"github.com/blinklabs-io/gouroboros/ledger/babbage"
	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/blinklabs-io/gouroboros/protocol/localstatequery"
	"github.com/ipfs/go-cid"
)

func (client *Client) TipsV1(address common.Address) ([]Tip, error) {
	var query *localstatequery.UTxOByAddressResult
	query, err := client.Query.Client.GetUTxOByAddress([]common.Address{address})
	if err != nil {
		return nil, err
	}
	var tips []Tip
	for id, output := range query.Results {
		tip, err := tipV1(id, output)
		if err == nil {
			tips = append(tips, *tip)
		}
	}
	return tips, nil
}

type Tip struct {
	TxID       localstatequery.UtxoId
	TxOut      babbage.BabbageTransactionOutput
	Credential common.Credential
	HeaderCid  cid.Cid
}

type TipRep struct {
	TxID             string
	ScriptCredential bool
	CredentialHash   string
	HeaderCid        string
}

func TipReps(tips []Tip) []TipRep {
	var result []TipRep
	for _, tip := range tips {
		result = append(result, *tip.Rep())
	}
	return result
}

func (tip *Tip) Rep() *TipRep {
	return &TipRep{
		TxID:             fmt.Sprintf("%v#%v", tip.TxID.Hash, tip.TxID.Idx),
		ScriptCredential: tip.Credential.CredType&0x10 != 0,
		CredentialHash:   fmt.Sprintf("%x", tip.Credential.Credential.Bytes()),
		HeaderCid:        tip.HeaderCid.String(),
	}
}

func tipV1(id localstatequery.UtxoId, output babbage.BabbageTransactionOutput) (*Tip, error) {
	if output.DatumOption == nil {
		return nil, fmt.Errorf("no datum")
	}
	var tip Tip
	tip.TxID = id
	tip.TxOut = output
	data, err := output.DatumOption.MarshalCBOR()
	if err != nil {
		return nil, err
	}
	var entry []interface{}
	_, err = cbor.Decode(data, &entry)
	if err != nil {
		return nil, err
	} else if len(entry) < 2 {
		return nil, fmt.Errorf("entry too short")
	}
	entry1, ok := entry[1].(cbor.WrappedCbor)
	if !ok {
		return nil, fmt.Errorf("entry is not wrapped CBOR")
	}
	var pair []interface{}
	_, err = cbor.Decode(entry1, &pair)
	if err != nil {
		return nil, err
	} else if len(pair) != 2 {
		return nil, fmt.Errorf("pair has length %v", len(pair))
	}
	pair0, ok := pair[0].(cbor.Tag)
	if !ok {
		return nil, fmt.Errorf("first entry of pair is not a constructor")
	}
	fields, ok := pair0.Content.([]interface{})
	if !ok || len(fields) != 1 {
		return nil, fmt.Errorf("multiple fields in constructor")
	}
	credential, ok := fields[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("field is not bytes")
	}
	tip.Credential = common.Credential{
		CredType:   uint(pair0.Number),
		Credential: common.CredentialHash(credential),
	}
	pair1, ok := pair[1].([]byte)
	if !ok {
		return nil, fmt.Errorf("second entry of pair is not bytes")
	}
	tip.HeaderCid, err = cid.Cast(pair1)
	if err != nil {
		return nil, err
	}
	return &tip, nil
}

package cardano

import (
	"fmt"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/cbor"
	"github.com/blinklabs-io/gouroboros/ledger/babbage"
	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/blinklabs-io/gouroboros/protocol/localstatequery"
	"github.com/ipfs/go-cid"
)

type Client struct {
	Node  *ouroboros.Connection
	Query *localstatequery.LocalStateQuery
}

func NewClient(nodeSocketPath string, networkMagic uint32) (*Client, error) {
	var err error
	var client Client
	client.Node, err = ouroboros.NewConnection(
		ouroboros.WithNetworkMagic(networkMagic),
		ouroboros.WithNodeToNode(false),
	)
	if err != nil {
		return nil, err
	}
	err = client.Node.Dial("unix", nodeSocketPath)
	if err != nil {
		return nil, err
	}
	client.Query = client.Node.LocalStateQuery()
	return &client, nil
}

type Tip struct {
	TxId       localstatequery.UtxoId
	TxOut      babbage.BabbageTransactionOutput
	Credential common.Credential
	HeaderCid  cid.Cid
}

type tipRep struct {
	TxId             string
	ScriptCredential bool
	CredentialHash   string
	HeaderCid        string
}

func TipReps(tips []Tip) []tipRep {
	var result []tipRep
	for _, tip := range tips {
		result = append(result, *tip.Rep())
	}
	return result
}

func (tip *Tip) Rep() *tipRep {
	return &tipRep{
		TxId:             fmt.Sprintf("%v#%v", tip.TxId.Hash, tip.TxId.Idx),
		ScriptCredential: tip.Credential.CredType&0x10 != 0,
		CredentialHash:   fmt.Sprintf("%x", tip.Credential.Credential.Bytes()),
		HeaderCid:        tip.HeaderCid.String(),
	}
}

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

func tipV1(id localstatequery.UtxoId, output babbage.BabbageTransactionOutput) (*Tip, error) {
	if output.DatumOption == nil {
		return nil, fmt.Errorf("no datum")
	}
	var tip Tip
	tip.TxId = id
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

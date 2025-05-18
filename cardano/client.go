package cardano

import (
	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/protocol/localstatequery"
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

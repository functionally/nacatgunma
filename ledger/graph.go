package ledger

import (
	"github.com/functionally/nacatgunma/header"
	"github.com/ipfs/go-cid"
	"gonum.org/v1/gonum/graph/simple"
)

type HeaderNode struct {
	Index     int64
	HeaderCid cid.Cid
	Header    *header.Header
}

func (node HeaderNode) ID() int64 {
	return node.Index
}

type HeaderTable struct {
	FromIndex map[int64]*HeaderNode
	FromCid   map[cid.Cid]*HeaderNode
}

func (ledger *Ledger) MakeHeaderTable() *HeaderTable {
	var table HeaderTable
	var i int64 = 0
	for hdrCid, hdr := range ledger.Headers {
		node := HeaderNode{
			Index:     i,
			HeaderCid: hdrCid,
			Header:    &hdr,
		}
		table.FromIndex[i] = &node
		table.FromCid[hdrCid] = &node
		i++
	}
	return &table
}

func (table *HeaderTable) MakeDirectedGraph(reverse bool) *simple.DirectedGraph {
	graph := simple.NewDirectedGraph()
	for _, node := range table.FromIndex {
		graph.AddNode(node)
	}
	for _, node := range table.FromCid {
		for _, acceptCid := range node.Header.Payload.Accept {
			if reverse {
				graph.SetEdge(graph.NewEdge(table.FromCid[acceptCid], node))
			} else {
				graph.SetEdge(graph.NewEdge(node, table.FromCid[acceptCid]))
			}
		}
	}
	return graph
}

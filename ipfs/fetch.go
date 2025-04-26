package ipfs

import (
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
)

func FetchNode(sh *shell.Shell, cid string) ([]byte, error) {
	return sh.BlockGet(cid)
}

func StoreNode(sh *shell.Shell, bytes []byte) (*cid.Cid, error) {
	c, err := sh.DagPut(bytes, "dag-cbor", "dag-cbor")
	if err != nil {
		return nil, err
	}
	cid, err := cid.Parse(c)
	if err != nil {
		return nil, err
	}
	return &cid, nil
}

package ipfs

import (
	shell "github.com/ipfs/go-ipfs-api"
)

func FetchNode(sh *shell.Shell, cid string) ([]byte, error) {
	return sh.BlockGet(cid)
}

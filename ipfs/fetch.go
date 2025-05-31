package ipfs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/functionally/nacatgunma/header"
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
)

func FetchChain(sh *shell.Shell, headerCid string, headerDir string, bodyDir *string, force bool, progress bool) error {
	headerFile := filepath.Join(headerDir, headerCid)
	if !force {
		_, err := os.Stat(headerFile)
		if err == nil || !os.IsNotExist(err) {
			if progress {
				log.Printf("Block header previously fetched: %v\n", headerCid)
			}
			return nil
		}
	}
	hdr, err := FetchHeader(sh, headerCid, &headerFile)
	if err != nil {
		return err
	}
	if progress {
		log.Printf("Fetched and verified block header: %v\n", headerCid)
	}
	if bodyDir != nil {
		bodyCid := hdr.Payload.Body.String()
		bodyFile := filepath.Join(*bodyDir, bodyCid)
		if !force {
			_, err := os.Stat(bodyFile)
			if err == nil || !os.IsNotExist(err) {
				log.Printf("Block body previously fetched: %v\n", bodyCid)
			}
		} else {
			bodyBytes, err := FetchNode(sh, bodyCid)
			if err != nil {
				return nil
			}
			err = os.WriteFile(bodyFile, bodyBytes, 0644)
			if err != nil {
				return err
			}
			if progress {
				log.Printf("Fetched block body: %v\n", bodyCid)
			}
		}
	}
	for _, acceptCid := range hdr.Payload.Accept {
		err = FetchChain(sh, acceptCid.String(), headerDir, bodyDir, force, progress)
		if err != nil {
			return err
		}
	}
	return nil
}

func FetchHeader(sh *shell.Shell, cid string, outputFile *string) (*header.Header, error) {
	headerBytes, err := FetchNode(sh, cid)
	if err != nil {
		return nil, err
	}
	hdr, err := header.UnmarshalHeader(headerBytes)
	if err != nil {
		return nil, err
	}
	verified, err := hdr.Verify()
	if err != nil {
		return nil, err
	} else if !verified {
		return nil, fmt.Errorf("header verification failed: %v", cid)
	}
	if outputFile != nil {
		err = os.WriteFile(*outputFile, headerBytes, 0644)
		if err != nil {
			return nil, err
		}
	}
	return hdr, nil
}

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

package cmd

import (
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"

	"github.com/functionally/nacatgunma/header"
	"github.com/functionally/nacatgunma/ipfs"
	"github.com/functionally/nacatgunma/key"
	"github.com/urfave/cli/v2"
)

func IpfsCmds() *cli.Command {
	return &cli.Command{
		Name:  "ipfs",
		Usage: "Interact with IPFS",
		Subcommands: []*cli.Command{
			ipfsChainCmd(),
			ipfsFetchCmd(),
			ipfsStoreCmd(),
		},
	}
}

func ipfsChainCmd() *cli.Command {

	var headerDir string
	var bodyDir string
	var tipCid string
	var ipfsApi string
	var force bool
	var progress bool

	return &cli.Command{
		Name:  "chain",
		Usage: "Fetch the block chain from IPFS.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ipfs-api",
				Value:       "localhost:5001",
				Usage:       "Endpoint for the IPFS API",
				Destination: &ipfsApi,
			},
			&cli.StringFlag{
				Name:        "tip-cid",
				Required:    true,
				Usage:       "The CID for the block header of the tip of the chain",
				Destination: &tipCid,
			},
			&cli.StringFlag{
				Name:        "header-dir",
				Required:    true,
				Usage:       "Output folder for the block headers",
				Destination: &headerDir,
			},
			&cli.StringFlag{
				Name:        "body-dir",
				Required:    false,
				Usage:       "Output folder for the block bodies",
				Destination: &bodyDir,
			},
			&cli.BoolFlag{
				Name:        "force",
				Value:       false,
				Usage:       "Fetch blocks that already exist in output folders",
				Destination: &force,
			},
			&cli.BoolFlag{
				Name:        "progress",
				Value:       false,
				Usage:       "Report progress",
				Destination: &progress,
			},
		},
		Action: func(ctx *cli.Context) error {
			sh := shell.NewShell(ipfsApi)
			var bodyDirPtr *string
			if ctx.IsSet("body-dir") {
				bodyDirPtr = &bodyDir
			} else {
				bodyDirPtr = nil
			}
			return ipfs.FetchChain(sh, tipCid, headerDir, bodyDirPtr, force, progress)
		},
	}
}

func ipfsFetchCmd() *cli.Command {

	var headerFile string
	var bodyFile string
	var headerCid string
	var ipfsApi string

	return &cli.Command{
		Name:  "fetch",
		Usage: "Fetch a block header and body from IPFS.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ipfs-api",
				Value:       "localhost:5001",
				Usage:       "Endpoint for the IPFS API",
				Destination: &ipfsApi,
			},
			&cli.StringFlag{
				Name:        "header-cid",
				Required:    true,
				Usage:       "The CID for the block header",
				Destination: &headerCid,
			},
			&cli.StringFlag{
				Name:        "header-file",
				Required:    true,
				Usage:       "Output file for the header",
				Destination: &headerFile,
			},
			&cli.StringFlag{
				Name:        "body-file",
				Required:    true,
				Usage:       "Output file for the body",
				Destination: &bodyFile,
			},
		},
		Action: func(*cli.Context) error {
			sh := shell.NewShell(ipfsApi)
			hdrBytes, err := ipfs.FetchNode(sh, headerCid)
			if err != nil {
				return err
			}
			hdr, err := header.UnmarshalHeader(hdrBytes)
			if err != nil {
				return err
			}
			bdyBytes, err := ipfs.FetchNode(sh, hdr.Payload.Body.String())
			if err != nil {
				return err
			}
			err = os.WriteFile(headerFile, hdrBytes, 0644)
			if err != nil {
				return err
			}
			err = os.WriteFile(bodyFile, bdyBytes, 0644)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func ipfsStoreCmd() *cli.Command {

	var keyFile string
	var bodyFile string
	var ipfsApi string
	var payload header.Payload
	var accepts cli.StringSlice
	var rejects cli.StringSlice

	return &cli.Command{
		Name:  "store",
		Usage: "Store a block on IPFS.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ipfs-api",
				Value:       "localhost:5001",
				Usage:       "Endpoint for the IPFS API",
				Destination: &ipfsApi,
			},
			&cli.StringFlag{
				Name:        "key-file",
				Required:    true,
				Usage:       "Input file for private key",
				Destination: &keyFile,
			},
			&cli.Int64Flag{
				Name:        "version",
				Value:       1,
				Usage:       "Header version number",
				Destination: &payload.Version,
			},
			&cli.StringFlag{
				Name:        "schema",
				Value:       "DAG-CBOR",
				Usage:       "Schema for the body",
				Destination: &payload.SchemaUri,
			},
			&cli.StringSliceFlag{
				Name:        "accept",
				Usage:       "Accept a CID as a parent",
				Destination: &accepts,
			},
			&cli.StringSliceFlag{
				Name:        "reject",
				Usage:       "Reject a CID as an ancestor",
				Destination: &rejects,
			},
			&cli.StringFlag{
				Name:        "body-file",
				Required:    true,
				Usage:       "Input file for the block body",
				Destination: &bodyFile,
			},
			&cli.StringFlag{
				Name:        "media-type",
				Value:       "application/cbor",
				Usage:       "Media type for body",
				Destination: &payload.MediaType,
			},
		},
		Action: func(*cli.Context) error {
			sh := shell.NewShell(ipfsApi)
			k, err := key.ReadPrivateKey(keyFile)
			if err != nil {
				return err
			}
			bodyBytes, err := os.ReadFile(bodyFile)
			if err != nil {
				return err
			}
			bodyCid, err := ipfs.StoreNode(sh, bodyBytes)
			if err != nil {
				return err
			}
			payload.Body = *bodyCid
			acceptCids, err := parseCIDs(uniqueStrings(accepts.Value()))
			if err != nil {
				return err
			}
			payload.Accept = acceptCids
			rejectCids, err := parseCIDs(uniqueStrings(rejects.Value()))
			if err != nil {
				return err
			}
			payload.Reject = rejectCids
			header, err := payload.Sign(k)
			if err != nil {
				return err
			}
			headerBytes, err := header.Marshal()
			if err != nil {
				return err
			}
			headerCid, err := ipfs.StoreNode(sh, headerBytes)
			if err != nil {
				return err
			}
			fmt.Println(headerCid)
			return nil
		},
	}

}

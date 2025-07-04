package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/functionally/nacatgunma/header"
	"github.com/functionally/nacatgunma/ipfs"
	"github.com/functionally/nacatgunma/key"
	"github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"
)

func HeaderCmds() *cli.Command {
	return &cli.Command{
		Name:  "header",
		Usage: "Header management subcommands",
		Subcommands: []*cli.Command{
			headerBuildCmd(),
			headerExportCmd(),
			headerVerifyCmd(),
		},
	}
}

func headerBuildCmd() *cli.Command {

	var keyFile string
	var headerFile string
	var payload header.Payload
	var body string
	var accepts cli.StringSlice
	var rejects cli.StringSlice

	return &cli.Command{
		Name:  "build",
		Usage: "Build a block header.",
		Flags: []cli.Flag{
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
				Value:       "https://w3c.github.io/json-ld-cbor/",
				Usage:       "Schema for the block body",
				Destination: &payload.SchemaURI,
			},
			&cli.StringSliceFlag{
				Name:        "accept",
				Usage:       "Accept a CID as a parent block",
				Destination: &accepts,
			},
			&cli.StringSliceFlag{
				Name:        "reject",
				Usage:       "Reject a CID as an ancestor block",
				Destination: &rejects,
			},
			&cli.StringFlag{
				Name:        "body",
				Required:    true,
				Usage:       "CID for the block body",
				Destination: &body,
			},
			&cli.StringFlag{
				Name:        "media-type",
				Value:       "application/vnd.ipld.dag-cbor",
				Usage:       "Media type for block body",
				Destination: &payload.MediaType,
			},
			&cli.StringFlag{
				Name:        "comment",
				Value:       "",
				Usage:       "Creator-supplied comment on the block",
				Destination: &payload.Comment,
			},
			&cli.StringFlag{
				Name:        "header-file",
				Required:    true,
				Usage:       "Output file for the block header CBOR",
				Destination: &headerFile,
			},
		},
		Action: func(*cli.Context) error {
			k, err := key.ReadPrivateKey(keyFile)
			if err != nil {
				return err
			}
			bodyCid, err := cid.Parse(body)
			if err != nil {
				return err
			}
			payload.Body = bodyCid
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
			headerCid, err := ipfs.CidV1(headerBytes)
			if err != nil {
				return err
			}
			err = os.WriteFile(headerFile, headerBytes, 0644)
			if err != nil {
				return err
			}
			fmt.Println(headerCid)
			return nil
		},
	}
}

func headerExportCmd() *cli.Command {

	var headerFile string
	var jsonFile string

	return &cli.Command{
		Name:  "export",
		Usage: "Export a block header to JSON.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "header-file",
				Required:    true,
				Usage:       "Input file for the block header CBOR",
				Destination: &headerFile,
			},
			&cli.StringFlag{
				Name:        "output-file",
				Required:    true,
				Usage:       "Output JSON file for the block header",
				Destination: &jsonFile,
			},
		},
		Action: func(*cli.Context) error {
			headerBytes, err := os.ReadFile(headerFile)
			if err != nil {
				return err
			}
			header, err := header.UnmarshalHeader(headerBytes)
			if err != nil {
				return err
			}
			json, err := json.MarshalIndent(header, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal header: %w", err)
			}
			return os.WriteFile(jsonFile, json, 0644)
		},
	}
}

func headerVerifyCmd() *cli.Command {

	var headerFile string

	return &cli.Command{
		Name:  "verify",
		Usage: "Verify a block header.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "header-file",
				Required:    true,
				Usage:       "Input file for the block header CBOR",
				Destination: &headerFile,
			},
		},
		Action: func(*cli.Context) error {
			headerBytes, err := os.ReadFile(headerFile)
			if err != nil {
				return err
			}
			header, err := header.UnmarshalHeader(headerBytes)
			if err != nil {
				return err
			}
			okay, err := header.Verify()
			if err != nil {
				return err
			} else if !okay {
				return fmt.Errorf("signature verification failed")
			}
			fmt.Printf("Verified signature by %s\n", header.Issuer)
			return nil
		},
	}

}

func parseCIDs(strs []string) ([]cid.Cid, error) {
	cids := make([]cid.Cid, 0, len(strs))
	for _, s := range strs {
		c, err := cid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CID %q: %w", s, err)
		}
		cids = append(cids, c)
	}
	return cids, nil
}

func uniqueStrings(input []string) []string {
	set := make(map[string]struct{})

	for _, s := range input {
		set[s] = struct{}{}
	}

	result := make([]string, 0, len(set))
	for s := range set {
		result = append(result, s)
	}

	return result
}

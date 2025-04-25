package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/functionally/achain/achain"
	"github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"
)

func main() {

	var keyFile string
	var headerFile string
	var resolutionFile string
	var payload achain.Payload
	var body string
	var accepts cli.StringSlice
	var rejects cli.StringSlice

	app := &cli.App{
		Name:  "achain",
		Usage: "Manage the AChain.",
		Commands: []*cli.Command{
			{
				Name:  "generate-key",
				Usage: "Generate an Ed25519 key.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "key-file",
						Required:    true,
						Usage:       "Output file for private key",
						Destination: &keyFile,
					},
				},
				Action: func(*cli.Context) error {
					key, err := achain.GenerateKey()
					if err != nil {
						return err
					}
					err = key.WritePrivateKey(keyFile)
					if err != nil {
						return err
					}
					fmt.Println(key.Did)
					return nil
				},
			},
			{
				Name:  "resolve-key",
				Usage: "Resolve an Ed25519 key.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "key-file",
						Required:    true,
						Usage:       "Input file for private key",
						Destination: &keyFile,
					},
					&cli.StringFlag{
						Name:        "resolution-file",
						Required:    true,
						Usage:       "Output file for DID resolution",
						Destination: &resolutionFile,
					},
				},
				Action: func(*cli.Context) error {
					key, err := achain.ReadPrivateKey(keyFile)
					if err != nil {
						return err
					}
					resolution, err := json.MarshalIndent(key.Resolution, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal DocResolution: %w", err)
					}
					err = os.WriteFile(resolutionFile, resolution, 0644)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "build-header",
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
						Value:       "DAG-CBOR",
						Usage:       "Schema for the body",
						Destination: &payload.Schema,
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
						Name:        "body",
						Usage:       "CID for the body",
						Destination: &body,
					},
					&cli.StringFlag{
						Name:        "media-type",
						Value:       "application/cbor",
						Usage:       "Media type for body",
						Destination: &payload.MediaType,
					},
					&cli.StringFlag{
						Name:        "header-file",
						Required:    true,
						Usage:       "Output file for the header CBOR",
						Destination: &headerFile,
					},
				},
				Action: func(*cli.Context) error {
					key, err := achain.ReadPrivateKey(keyFile)
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
					header, err := payload.Sign(key)
					if err != nil {
						return err
					}
					headerCid, headerBytes, err := header.Marshal()
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
			},
			{
				Name:  "verify-header",
				Usage: "Verify a block header.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "header-file",
						Required:    true,
						Usage:       "Input file for the header CBOR",
						Destination: &headerFile,
					},
				},
				Action: func(*cli.Context) error {
					headerBytes, err := os.ReadFile(headerFile)
					if err != nil {
						return err
					}
					header, err := achain.UnmarshalHeader(headerBytes)
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
			},
		},
	}

	if appErr := app.Run(os.Args); appErr != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", appErr)
		os.Exit(1)
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

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli/v2"

	"github.com/functionally/nacatgunma/cardano"
	"github.com/functionally/nacatgunma/header"
	"github.com/functionally/nacatgunma/ipfs"
	"github.com/functionally/nacatgunma/key"
	"github.com/functionally/nacatgunma/rdf"
)

func main() {

	var keyFile string
	var headerFile string
	var jsonFile string
	var rdfFile string
	var bodyFile string
	var baseUri string
	var format string
	var keyDid string
	var headerCid string
	var ipfsApi string
	var payload header.Payload
	var body string
	var accepts cli.StringSlice
	var rejects cli.StringSlice
	var nodeSocket string
	var networkMagic uint
	var address string
	var tipFile string
	var script bool
	var credential string
	var datumFile string
	var metadataKey uint
	var redeemerFile string
	var blockchain string
	var metadataFile string

	app := &cli.App{
		Name:  "nacatgunma",
		Usage: "Manage the Nacatgunma blockchain.",
		Commands: []*cli.Command{

			{
				Name:  "body",
				Usage: "Body management subcommands",
				Subcommands: []*cli.Command{

					{
						Name:  "export",
						Usage: "Export a body as JSON.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "body-file",
								Required:    true,
								Usage:       "Input file for the block body",
								Destination: &bodyFile,
							},
							&cli.StringFlag{
								Name:        "output-file",
								Required:    true,
								Usage:       "Output file for the block body",
								Destination: &jsonFile,
							},
						},
						Action: func(*cli.Context) error {
							bodyBytes, err := os.ReadFile(bodyFile)
							if err != nil {
								return err
							}
							body, err := ipfs.DecodeFromDagCbor(bodyBytes)
							if err != nil {
								return err
							}
							json, err := json.MarshalIndent(body, "", "  ")
							if err != nil {
								return fmt.Errorf("failed to marshal body: %w", err)
							}
							return os.WriteFile(jsonFile, json, 0644)
						},
					},

					{
						Name:  "rdf",
						Usage: "Build a block of RDF N-quads.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "rdf-file",
								Required:    true,
								Usage:       "Input file of RDF N-quads",
								Destination: &rdfFile,
							},
							&cli.StringFlag{
								Name:        "base-uri",
								Value:       "",
								Usage:       "Base URI of the RDF",
								Destination: &baseUri,
							},
							&cli.StringFlag{
								Name:        "format",
								Value:       "application/n-quads",
								Usage:       "MIME type of the RDF format",
								Destination: &format,
							},
							&cli.StringFlag{
								Name:        "body-file",
								Required:    true,
								Usage:       "Output file for the block body",
								Destination: &bodyFile,
							},
						},
						Action: func(*cli.Context) error {
							rdf, err := rdf.ReadRdf(rdfFile, baseUri, format)
							if err != nil {
								return err
							}
							bodyBytes, err := ipfs.EncodeToDagCbor(rdf)
							if err != nil {
								return err
							}
							bodyCid, err := ipfs.CidV1(bodyBytes)
							if err != nil {
								return err
							}
							err = os.WriteFile(bodyFile, bodyBytes, 0644)
							if err != nil {
								return err
							}
							fmt.Println(bodyCid)
							return nil
						},
					},
				},
			},

			{
				Name:  "cardano",
				Usage: "Interact with Cardano",
				Subcommands: []*cli.Command{

					{
						Name:  "datum",
						Usage: "Create datum for a tip.",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:        "script",
								Required:    true,
								Usage:       "Whether the credential is a script instead of a public key",
								Destination: &script,
							},
							&cli.StringFlag{
								Name:        "credential-hash",
								Required:    true,
								Usage:       "Blake2b224 hash of the credential, in hexadecimal",
								Destination: &credential,
							},
							&cli.StringFlag{
								Name:        "header-cid",
								Required:    true,
								Usage:       "The CID of the block header",
								Destination: &headerCid,
							},
							&cli.StringFlag{
								Name:        "datum-file",
								Value:       "/dev/stdout",
								Usage:       "Output file for JSON-formatted datum",
								Destination: &datumFile,
							},
						},
						Action: func(*cli.Context) error {
							datum, err := cardano.NewDatum(script, credential, headerCid)
							if err != nil {
								return err
							}
							datumBytes, err := datum.ToJSON()
							if err != nil {
								return err
							}
							return os.WriteFile(datumFile, datumBytes, 0644)
						},
					},

					{
						Name:  "inputs",
						Usage: "Create inputs for a tip.",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:        "metadata-key",
								Value:       58312,
								Usage:       "Metadata key for the chain",
								Destination: &metadataKey,
							},
							&cli.StringFlag{
								Name:        "blockchain",
								Value:       "https://github.com/functionally/nacatgunma",
								Usage:       "The IRI identifying the blockchain",
								Destination: &blockchain,
							},
							&cli.BoolFlag{
								Name:        "script",
								Required:    true,
								Usage:       "Whether the credential is a script instead of a public key",
								Destination: &script,
							},
							&cli.StringFlag{
								Name:        "credential-hash",
								Required:    true,
								Usage:       "Blake2b224 hash of the credential, in hexadecimal",
								Destination: &credential,
							},
							&cli.StringFlag{
								Name:        "header-cid",
								Required:    true,
								Usage:       "The CID of the block header",
								Destination: &headerCid,
							},
							&cli.StringFlag{
								Name:        "datum-file",
								Required:    true,
								Usage:       "Output file for JSON-formatted datum",
								Destination: &datumFile,
							},
							&cli.StringFlag{
								Name:        "redeemer-file",
								Required:    true,
								Usage:       "Output file for JSON-formatted redeemer",
								Destination: &redeemerFile,
							},
							&cli.StringFlag{
								Name:        "metadata-file",
								Required:    true,
								Usage:       "Output file for JSON-formatted metadata",
								Destination: &metadataFile,
							},
						},
						Action: func(*cli.Context) error {
							datum, err := cardano.NewDatum(script, credential, headerCid)
							if err != nil {
								return err
							}
							datumBytes, err := datum.ToJSON()
							if err != nil {
								return err
							}
							err = os.WriteFile(datumFile, datumBytes, 0644)
							if err != nil {
								return err
							}
							redeemerBytes, err := cardano.RedeemerJSON(metadataKey)
							if err != nil {
								return err
							}
							err = os.WriteFile(redeemerFile, redeemerBytes, 0644)
							if err != nil {
								return err
							}
							metadataBytes, err := cardano.MetadataJSON(metadataKey, blockchain, headerCid)
							if err != nil {
								return err
							}
							return os.WriteFile(metadataFile, metadataBytes, 0644)
						},
					},

					{
						Name:  "metadata",
						Usage: "Create metadata for a tip.",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:        "metadata-key",
								Value:       58312,
								Usage:       "Metadata key for the chain",
								Destination: &metadataKey,
							},
							&cli.StringFlag{
								Name:        "blockchain",
								Value:       "https://github.com/functionally/nacatgunma",
								Usage:       "The IRI identifying the blockchain",
								Destination: &blockchain,
							},
							&cli.StringFlag{
								Name:        "header-cid",
								Required:    true,
								Usage:       "The CID of the block header",
								Destination: &headerCid,
							},
							&cli.StringFlag{
								Name:        "metadata-file",
								Value:       "/dev/stdout",
								Usage:       "Output file for JSON-formatted metadata",
								Destination: &metadataFile,
							},
						},
						Action: func(*cli.Context) error {
							metadataBytes, err := cardano.MetadataJSON(metadataKey, blockchain, headerCid)
							if err != nil {
								return err
							}
							return os.WriteFile(metadataFile, metadataBytes, 0644)
						},
					},

					{
						Name:  "redeemer",
						Usage: "Create redeemer for a tip.",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:        "metadata-key",
								Value:       58312,
								Usage:       "Metadata key for the chain",
								Destination: &metadataKey,
							},
							&cli.StringFlag{
								Name:        "redeemer-file",
								Value:       "/dev/stdout",
								Usage:       "Output file for JSON-formatted redeemer",
								Destination: &redeemerFile,
							},
						},
						Action: func(*cli.Context) error {
							redeemerBytes, err := cardano.RedeemerJSON(metadataKey)
							if err != nil {
								return err
							}
							return os.WriteFile(redeemerFile, redeemerBytes, 0644)
						},
					},

					{
						Name:  "tips",
						Usage: "Find the tips.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "node-socket",
								Required:    true,
								Usage:       "Path to the Cardano node socket",
								Destination: &nodeSocket,
							},
							&cli.UintFlag{
								Name:        "network-magic",
								Value:       764824073,
								Usage:       "Magic number for the Cardano network",
								Destination: &networkMagic,
							},
							&cli.StringFlag{
								Name:        "script-address",
								Required:    true,
								Usage:       "Address of the Plutus script for the tip",
								Destination: &address,
							},
							&cli.StringFlag{
								Name:        "tips-file",
								Value:       "/dev/stdout",
								Usage:       "Output JSON file for tip information",
								Destination: &tipFile,
							},
						},
						Action: func(*cli.Context) error {
							addr, err := common.NewAddress(address)
							if err != nil {
								return err
							}
							client, err := cardano.NewClient(nodeSocket, uint32(networkMagic))
							if err != nil {
								return err
							}
							tips, err := client.TipsV1(addr)
							if err != nil {
								return err
							}
							jsonBytes, _ := json.MarshalIndent(cardano.TipReps(tips), "", "  ")
							return os.WriteFile(tipFile, jsonBytes, 0644)
						},
					},
				},
			},

			{
				Name:  "header",
				Usage: "Header management subcommands",
				Subcommands: []*cli.Command{

					{
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
								Destination: &payload.SchemaUri,
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
							key, err := key.ReadPrivateKey(keyFile)
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
					},

					{
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
					},

					{
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
					},
				},
			},

			{
				Name:  "ipfs",
				Usage: "Interact with IPFS",
				Subcommands: []*cli.Command{

					{
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
					},

					{
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
							key, err := key.ReadPrivateKey(keyFile)
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
							header, err := payload.Sign(key)
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
					},
				},
			},

			{
				Name:  "key",
				Usage: "Key management subcommands",
				Subcommands: []*cli.Command{

					{
						Name:  "generate",
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
							key, err := key.GenerateKey()
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
						Name:  "resolve",
						Usage: "Resolve an Ed25519 key.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "key-did",
								Required:    true,
								Usage:       "The DID for the public key",
								Destination: &keyDid,
							},
							&cli.StringFlag{
								Name:        "output-file",
								Required:    true,
								Usage:       "Output JSON file for DID resolution",
								Destination: &jsonFile,
							},
						},
						Action: func(*cli.Context) error {
							resolution, err := key.ResolveDid(keyDid)
							if err != nil {
								return err
							}
							json, err := json.MarshalIndent(resolution, "", "  ")
							if err != nil {
								return fmt.Errorf("failed to marshal DocResolution: %w", err)
							}
							err = os.WriteFile(jsonFile, json, 0644)
							if err != nil {
								return err
							}
							return nil
						},
					},
				},
			},
		}}

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

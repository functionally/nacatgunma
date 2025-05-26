package cmd

import (
	"encoding/json"
	"os"

	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/urfave/cli/v2"

	"github.com/functionally/nacatgunma/cardano"
)

func CardanoCmds() *cli.Command {
	return &cli.Command{
		Name:  "cardano",
		Usage: "Interact with Cardano",
		Subcommands: []*cli.Command{
			cardanoDatumCmd(),
			cardanoInputsCmd(),
			cardanoMetadataCmd(),
			cardanoRedeemerCmd(),
			cardanoTipsCmd(),
		},
	}
}

func cardanoDatumCmd() *cli.Command {

	var headerCid string
	var script bool
	var credential string
	var datumFile string

	return &cli.Command{
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
	}

}

func cardanoInputsCmd() *cli.Command {

	var headerCid string
	var script bool
	var credential string
	var datumFile string
	var metadataKey uint
	var redeemerFile string
	var blockchain string
	var metadataFile string

	return &cli.Command{
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
	}

}

func cardanoMetadataCmd() *cli.Command {

	var headerCid string
	var metadataKey uint
	var blockchain string
	var metadataFile string

	return &cli.Command{
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
	}

}

func cardanoRedeemerCmd() *cli.Command {

	var metadataKey uint
	var redeemerFile string

	return &cli.Command{
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
	}

}

func cardanoTipsCmd() *cli.Command {

	var nodeSocket string
	var networkMagic uint
	var address string
	var tipFile string

	return &cli.Command{
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
	}

}

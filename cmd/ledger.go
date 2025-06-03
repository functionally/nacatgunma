package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/functionally/nacatgunma/ledger"
)

func LedgerCmds() *cli.Command {
	return &cli.Command{
		Name:  "ledger",
		Usage: "Body management subcommands",
		Subcommands: []*cli.Command{
			ledgerExportCmd(),
			ledgerPruneCmd(),
		},
	}
}

func ledgerExportCmd() *cli.Command {

	var tipCid string
	var headerDir string
	var turtleFile string
	var jsonFile string

	return &cli.Command{
		Name:  "export",
		Usage: "Export headers from the ledger.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tip-cid",
				Required:    true,
				Usage:       "The CID for the block header of the tip of the chain",
				Destination: &tipCid,
			},
			&cli.StringFlag{
				Name:        "header-dir",
				Required:    true,
				Usage:       "Input folder for the block headers",
				Destination: &headerDir,
			},
			&cli.StringFlag{
				Name:        "turtle-file",
				Required:    false,
				Usage:       "Output file for the block headers in Turtle format",
				Destination: &turtleFile,
			},
			&cli.StringFlag{
				Name:        "json-file",
				Required:    false,
				Usage:       "Output file for the block headers in JSON format",
				Destination: &jsonFile,
			},
		},
		Action: func(ctx *cli.Context) error {
			ledger, err := ledger.ReadLedger(tipCid, headerDir)
			if err != nil {
				return err
			}
			prunable := ledger.Prunable()
			if ctx.IsSet("turtle-file") {
				err = ledger.WriteLedgerTurtle(turtleFile, prunable)
				if err != nil {
					return err
				}
			}
			if ctx.IsSet("json-file") {
				json, err := json.MarshalIndent(ledger, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal ledger: %w", err)
				}
				err = os.WriteFile(jsonFile, json, 0644)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func ledgerPruneCmd() *cli.Command {

	var tipCid string
	var headerDir string
	var acceptedFile string
	var rejectedFile string
	var bodyFile string

	return &cli.Command{
		Name:  "prune",
		Usage: "Prune rejected blocks from the ledger.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tip-cid",
				Required:    true,
				Usage:       "The CID for the block header of the tip of the chain",
				Destination: &tipCid,
			},
			&cli.StringFlag{
				Name:        "header-dir",
				Required:    true,
				Usage:       "Input folder for the block headers",
				Destination: &headerDir,
			},
			&cli.StringFlag{
				Name:        "accepted-file",
				Required:    false,
				Usage:       "Output file for the list of accepted block headers",
				Destination: &acceptedFile,
			},
			&cli.StringFlag{
				Name:        "rejected-file",
				Required:    false,
				Usage:       "Output file for the list of rejected block headers",
				Destination: &rejectedFile,
			},
			&cli.StringFlag{
				Name:        "body-file",
				Required:    false,
				Usage:       "Output file for the list of accepted block bodies",
				Destination: &bodyFile,
			},
		},
		Action: func(ctx *cli.Context) error {
			ledger, err := ledger.ReadLedger(tipCid, headerDir)
			if err != nil {
				return err
			}
			rejected := ledger.Prune()
			if ctx.IsSet("accepted-file") {
				f, err := os.OpenFile(acceptedFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
				defer f.Close()
				for hdrCid := range ledger.Headers {
					_, err = f.WriteString(fmt.Sprintf("%v\n", hdrCid.String()))
					if err != nil {
						return err
					}
				}
			}
			if ctx.IsSet("rejected-file") {
				f, err := os.OpenFile(rejectedFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
				defer f.Close()
				for _, hdrCid := range rejected {
					_, err = f.WriteString(fmt.Sprintf("%v\n", hdrCid.String()))
					if err != nil {
						return err
					}
				}
			}
			if ctx.IsSet("body-file") {
				f, err := os.OpenFile(bodyFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
				defer f.Close()
				for _, hdr := range ledger.Headers {
					_, err = f.WriteString(fmt.Sprintf("%v\n", hdr.Payload.Body.String()))
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

}

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/functionally/nacatgunma/ipfs"
)

func BodyCmds() *cli.Command {
	return &cli.Command{
		Name:  "body",
		Usage: "Body management subcommands",
		Subcommands: []*cli.Command{
			bodyExportCmd(),
			rdfCmd(),
			tgdhCmds(),
		},
	}
}

func bodyExportCmd() *cli.Command {

	var jsonFile string
	var bodyFile string

	return &cli.Command{
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
	}

}

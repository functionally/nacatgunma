package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/functionally/nacatgunma/key"
	"github.com/urfave/cli/v2"
)

func KeyCmds() *cli.Command {
	return &cli.Command{
		Name:  "key",
		Usage: "Key management subcommands",
		Subcommands: []*cli.Command{
			keyGenerateCmd(),
			keyResolveCmd(),
		},
	}
}

func keyGenerateCmd() *cli.Command {

	var keyFile string

	return &cli.Command{
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
	}

}

func keyResolveCmd() *cli.Command {

	var jsonFile string
	var keyDid string

	return &cli.Command{
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
	}

}

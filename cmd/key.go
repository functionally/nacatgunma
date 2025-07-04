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
			keyDidCmd(),
			keyGenerateCmd(),
			keyResolveCmd(),
		},
	}
}

func keyDidCmd() *cli.Command {

	var keyFile string

	return &cli.Command{
		Name:  "did",
		Usage: "Print the DID of a cryptographic key.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "key-file",
				Required:    true,
				Usage:       "Input file for private key",
				Destination: &keyFile,
			},
		},
		Action: func(*cli.Context) error {
			k, err := key.ReadPrivateKey(keyFile)
			if err != nil {
				return err
			}
			fmt.Println(key.Did(k))
			return nil
		},
	}

}

func keyGenerateCmd() *cli.Command {

	var keyFile string
	var keyType string

	return &cli.Command{
		Name:  "generate",
		Usage: "Generate a cryptographic key.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "key-file",
				Required:    true,
				Usage:       "Output file for private key",
				Destination: &keyFile,
			},
			&cli.StringFlag{
				Name:        "key-type",
				Value:       "Ed25519",
				Usage:       "The key type, either \"Ed25519\" or \"BLS12-381\"",
				Destination: &keyType,
			},
		},
		Action: func(*cli.Context) error {
			allowed := map[string]key.KeyType{
				"Ed25519":   key.Ed25519,
				"BLS12-381": key.Bls12381,
			}
			kt, okay := allowed[keyType]
			if !okay {
				return fmt.Errorf("unsupported key type: %v", keyType)
			}
			k, err := key.GenerateKey(kt)
			if err != nil {
				return err
			}
			err = key.WritePrivateKey(k, keyFile)
			if err != nil {
				return err
			}
			fmt.Println(key.Did(k))
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

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/functionally/nacatgunma/tgdh"
	"github.com/urfave/cli/v2"
)

func tgdhCmds() *cli.Command {
	return &cli.Command{
		Name:  "tgdh",
		Usage: "Tree-based group DH (BLS12-381) management subcommands",
		Subcommands: []*cli.Command{
			tgdhDecryptCmd(),
			tgdhEncryptCmd(),
			tgdhGenerateCmd(),
			tgdhJoinCmd(),
			tgdhPublicCmd(),
			tgdhPrivateCmd(),
		},
	}
}

func tgdhGenerateCmd() *cli.Command {

	var privateFile string

	return &cli.Command{
		Name:  "generate",
		Usage: "Generate a TGDH private key.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "private-file",
				Required:    true,
				Usage:       "Output file of TGDH private key",
				Destination: &privateFile,
			},
		},
		Action: func(*cli.Context) error {
			leaf, err := tgdh.GenerateLeaf()
			if err != nil {
				return err
			}
			leafBytes, err := leaf.MarshalJSON()
			if err != nil {
				return err
			}
			err = os.WriteFile(privateFile, leafBytes, 0644)
			if err != nil {
				return err
			}
			fmt.Println(leaf.Did())
			return nil
		},
	}

}

func tgdhPublicCmd() *cli.Command {

	var privateFile string
	var publicFile string

	return &cli.Command{
		Name:  "public",
		Usage: "Strip private key information from a TGDH key.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "private-file",
				Required:    true,
				Usage:       "Input file of TGDH private key",
				Destination: &privateFile,
			},
			&cli.StringFlag{
				Name:        "public-file",
				Required:    true,
				Usage:       "Output file of TGDH public key",
				Destination: &publicFile,
			},
		},
		Action: func(*cli.Context) error {
			privateBytes, err := os.ReadFile(privateFile)
			if err != nil {
				return err
			}
			private, err := tgdh.UnmarshalJSON(privateBytes)
			if err != nil {
				return err
			}
			public := private.DeepStrip()
			publicBytes, err := public.MarshalJSON()
			if err != nil {
				return err
			}
			err = os.WriteFile(publicFile, publicBytes, 0644)
			if err != nil {
				return err
			}
			fmt.Println(public.Did())
			return nil
		},
	}

}

func tgdhJoinCmd() *cli.Command {

	var leftFile string
	var rightFile string
	var privateFile string

	return &cli.Command{
		Name:  "join",
		Usage: "Join two TGDH keys into an aggregate TGDH key, where at least one of the keys is private.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "left-file",
				Required:    true,
				Usage:       "Input file of leftmost TGDH key",
				Destination: &leftFile,
			},
			&cli.StringFlag{
				Name:        "right-file",
				Required:    true,
				Usage:       "Input file of rightmost TGDH key",
				Destination: &rightFile,
			},
			&cli.StringFlag{
				Name:        "private-file",
				Required:    true,
				Usage:       "Output file of TGDH private key",
				Destination: &privateFile,
			},
		},
		Action: func(*cli.Context) error {
			leftBytes, err := os.ReadFile(leftFile)
			if err != nil {
				return err
			}
			left, err := tgdh.UnmarshalJSON(leftBytes)
			if err != nil {
				return err
			}
			rightBytes, err := os.ReadFile(rightFile)
			if err != nil {
				return err
			}
			right, err := tgdh.UnmarshalJSON(rightBytes)
			if err != nil {
				return err
			}
			private, err := tgdh.Join(left, right)
			if err != nil {
				return err
			}
			privateBytes, err := private.MarshalJSON()
			if err != nil {
				return err
			}
			err = os.WriteFile(privateFile, privateBytes, 0644)
			if err != nil {
				return err
			}
			fmt.Println(private.Did())
			return nil
		},
	}

}

func tgdhPrivateCmd() *cli.Command {

	var privateFile string
	var publicFile string
	var rootFile string

	return &cli.Command{
		Name:  "private",
		Usage: "Apply a private TGHD key to a public one, deriving the private root.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "private-file",
				Required:    true,
				Usage:       "Input file of private TGDH key",
				Destination: &privateFile,
			},
			&cli.StringFlag{
				Name:        "public-file",
				Required:    true,
				Usage:       "Input file of public TGDH key",
				Destination: &publicFile,
			},
			&cli.StringFlag{
				Name:        "root-file",
				Required:    true,
				Usage:       "Output file of TGDH root private key",
				Destination: &rootFile,
			},
		},
		Action: func(*cli.Context) error {
			privateBytes, err := os.ReadFile(privateFile)
			if err != nil {
				return err
			}
			private, err := tgdh.UnmarshalJSON(privateBytes)
			if err != nil {
				return err
			}
			publicBytes, err := os.ReadFile(publicFile)
			if err != nil {
				return err
			}
			public, err := tgdh.UnmarshalJSON(publicBytes)
			if err != nil {
				return err
			}
			root, err := tgdh.DerivePrivates(private, public)
			if err != nil {
				return err
			}
			rootBytes, err := root.MarshalJSON()
			if err != nil {
				return err
			}
			err = os.WriteFile(rootFile, rootBytes, 0644)
			if err != nil {
				return err
			}
			fmt.Println(root.Did())
			return nil
		},
	}

}

func tgdhEncryptCmd() *cli.Command {

	var privateFile string
	var plaintextFile string
	var contentType string
	var ciphertextFile string

	return &cli.Command{
		Name:  "encrypt",
		Usage: "Encrypt a file using a TGDH private key.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "private-file",
				Required:    true,
				Usage:       "Input file of private TGDH key",
				Destination: &privateFile,
			},
			&cli.StringFlag{
				Name:        "plaintext-file",
				Required:    true,
				Usage:       "Input file of plaintext",
				Destination: &plaintextFile,
			},
			&cli.StringFlag{
				Name:        "content-type",
				Value:       "",
				Usage:       "The content type for the plaintext",
				Destination: &contentType,
			},
			&cli.StringFlag{
				Name:        "ciphertext-file",
				Required:    true,
				Usage:       "Output file for the ciphertext",
				Destination: &ciphertextFile,
			},
		},
		Action: func(*cli.Context) error {
			privateBytes, err := os.ReadFile(privateFile)
			if err != nil {
				return err
			}
			private, err := tgdh.UnmarshalJSON(privateBytes)
			if err != nil {
				return err
			}
			plaintext, err := os.ReadFile(plaintextFile)
			if err != nil {
				return err
			}
			ciphertext, err := private.Encrypt(plaintext, contentType)
			if err != nil {
				return err
			}
			err = os.WriteFile(ciphertextFile, ciphertext, 0644)
			if err != nil {
				return err
			}
			return nil
		},
	}

}

func tgdhDecryptCmd() *cli.Command {

	var privateFile string
	var ciphertextFile string
	var plaintextFile string
	var headersFile string

	return &cli.Command{
		Name:  "decrypt",
		Usage: "Decrypt a file using a TGDH private key.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "private-file",
				Required:    true,
				Usage:       "Input file of private TGDH key",
				Destination: &privateFile,
			},
			&cli.StringFlag{
				Name:        "ciphertext-file",
				Required:    true,
				Usage:       "Input file of the ciphertext",
				Destination: &ciphertextFile,
			},
			&cli.StringFlag{
				Name:        "plaintext-file",
				Required:    true,
				Usage:       "Output file for plaintext",
				Destination: &plaintextFile,
			},
			&cli.StringFlag{
				Name:        "headers-file",
				Required:    false,
				Usage:       "Output file for encryption headers",
				Destination: &headersFile,
			},
		},
		Action: func(ctx *cli.Context) error {
			privateBytes, err := os.ReadFile(privateFile)
			if err != nil {
				return err
			}
			private, err := tgdh.UnmarshalJSON(privateBytes)
			if err != nil {
				return err
			}
			ciphertext, err := os.ReadFile(ciphertextFile)
			if err != nil {
				return err
			}
			headers, plaintext, err := private.Decrypt(ciphertext)
			if err != nil {
				return err
			}
			err = os.WriteFile(plaintextFile, plaintext, 0644)
			if err != nil {
				return err
			}
			if ctx.IsSet("headers-file") {
				headersBytes, err := json.Marshal(headers)
				if err != nil {
					return err
				}
				err = os.WriteFile(headersFile, headersBytes, 0644)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

}

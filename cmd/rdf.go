package cmd

import (
	"fmt"
	"os"

	"github.com/functionally/nacatgunma/ipfs"
	"github.com/functionally/nacatgunma/rdf"
	"github.com/urfave/cli/v2"
)

func rdfCmd() *cli.Command {

	var rdfFile string
	var bodyFile string
	var baseUri string
	var format string

	return &cli.Command{
		Name:  "rdf",
		Usage: "Build a body of RDF N-quads.",
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
	}

}

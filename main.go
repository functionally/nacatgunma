package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/functionally/nacatgunma/cmd"
)

func main() {
	app := &cli.App{
		Name:  "nacatgunma",
		Usage: "Manage the Nacatgunma blockchain.",
		Commands: []*cli.Command{
			cmd.BodyCmds(),
			cmd.CardanoCmds(),
			cmd.HeaderCmds(),
			cmd.IpfsCmds(),
			cmd.KeyCmds(),
			cmd.LedgerCmds(),
		}}
	if appErr := app.Run(os.Args); appErr != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", appErr)
		os.Exit(1)
	}
}

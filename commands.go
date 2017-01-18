package main

import (
	"os"

	"ukcloud-portal-cli/command"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available commands.
var Commands map[string]cli.CommandFactory
var PlumbingCommands map[string]struct{}

// Ui is the cli.Ui used for communicating to the outside world.
var Ui cli.Ui

const (
	ErrorPrefix  = "e:"
	OutputPrefix = "o:"
)

func init() {

	Ui = &cli.PrefixedUi{
		AskPrefix:    OutputPrefix,
		OutputPrefix: OutputPrefix,
		InfoPrefix:   OutputPrefix,
		ErrorPrefix:  ErrorPrefix,
		Ui:           &cli.BasicUi{Writer: os.Stdout},
	}

	meta := command.Meta{
		Color:       true,
		Ui:          Ui,
	}
	
	// The command list is included in the ukc -help
	// output

	PlumbingCommands = map[string]struct{}{
		"state": struct{}{}, // includes all subcommands
		"debug": struct{}{}, // includes all subcommands
	}

	Commands = map[string]cli.CommandFactory{
		"accounts": func() (cli.Command, error) {
			return &command.AccountsCommand{
				Meta:       meta,
			}, nil
		},
	}
}
package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/ukcloud/ukcloud-portal-cli/command"
)

// Commands is the mapping of all the available commands.
var Commands map[string]cli.CommandFactory

// PlumbingCommands is foe the state & Debug
var PlumbingCommands map[string]struct{}

// UI is the cli.Ui used for communicating to the outside world.
var UI cli.Ui

func init() {

	UI = &cli.ColoredUi{
		Ui: &cli.BasicUi{Writer: os.Stdout},
	}

	meta := command.Meta{
		Color: true,
		UI:    UI,
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
				Meta: meta,
			}, nil
		},

		"vdcs": func() (cli.Command, error) {
			return &command.VdcCommand{
				Meta: meta,
			}, nil
		},

		"vorgs": func() (cli.Command, error) {
			return &command.VorgsCommand{
				Meta: meta,
			}, nil
		},

		"vapps": func() (cli.Command, error) {
			return &command.VappCommand{
				Meta: meta,
			}, nil
		},

		"vapps shutdown": func() (cli.Command, error) {
			return &command.VappShutdownCommand{
				Meta: meta,
			}, nil
		},

		"vapps boot": func() (cli.Command, error) {
			return &command.VappBootCommand{
				Meta: meta,
			}, nil
		},

		"vdcbuild": func() (cli.Command, error) {
			return &command.VdcBuildCommand{
				Meta: meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:              meta,
				Revision:          GitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
			}, nil
		},
	}
}

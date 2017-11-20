package command

import (
	"strings"
)

// VorgsShutdownCommand is the command to list UKC Vorgs
type VappShutdownCommand struct {
	Meta
}

// Run is will be executed when `ukc vapps` is called
func (c *VappShutdownCommand) Run(args []string) int {

	cmd := VappCommand{
		Meta: c.Meta,
	}

	cmd.ChangeState(args, "shutdown")

	return 0
}

// Help is called when ukc accounts -help | --help | -h
func (c *VappShutdownCommand) Help() string {
	helpText := `
Usage: ukc vapps shutdown

  ` + c.Synopsis() + `

Required:
  -accountid       Your UKCloud AccountID

Options:
  -email           Email to your UKCloud Portal Account

  -password        Password to your UKCloud Portal Acccount
`
	return strings.TrimSpace(helpText)
}

// Synopsis is used when listing all commands or in the help
func (c *VappShutdownCommand) Synopsis() string {
	return "Control the Vapps in your UKCloud Account"
}

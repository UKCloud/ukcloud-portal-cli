package command

import (
	"bytes"
	"fmt"
)

// VersionCommand is a Command implementation prints the version.
type VersionCommand struct {
	Meta
	Revision          string
	Version           string
	VersionPrerelease string
}

// Help is for showing the help commands
func (c *VersionCommand) Help() string {
	return ""
}

// Run will show the version information
func (c *VersionCommand) Run(args []string) int {
	var versionString bytes.Buffer
	args = c.Meta.process(args, false)

	fmt.Fprintf(&versionString, "UKC v%s", c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	c.UI.Output(versionString.String())

	return 0
}

// Synopsis is used when listing all commands or in the help
func (c *VersionCommand) Synopsis() string {
	return "Prints the UKCloud Portal CLI tool version"
}

package command

import (
	"fmt"
	"github.com/ukcloud/ukcloud-portal-api/api"
	"os"
	"strings"
	"text/tabwriter"
)

// VdcBuildCommand is the command to list UKC Vdcs
type VdcBuildCommand struct {
	Meta
}

// Run is will be executed when `ukc accounts` is called
func (c *VdcBuildCommand) Run(args []string) int {

	var err error
	papi := new(api.API)
	args = c.Meta.process(args, false)

	// Get the flags for this)
	cmdName := "vdcs"
	cmdFlags := c.Meta.flagSet(cmdName)
	cmdFlags.StringVar(&c.Meta.email, "email", "", "email")
	cmdFlags.StringVar(&c.Meta.password, "password", "", "password")
	cmdFlags.IntVar(&c.Meta.buildID, "buildid", 0, "buildid")

	cmdFlags.Usage = func() { c.UI.Error(c.Help()) }
	err = cmdFlags.Parse(args)
	if c.Meta.buildID == 0 || err != nil {
		cmdFlags.Usage()
		return 1
	}

	var authorised bool
	authorised, err = papi.GetAuth(c.Meta.email, c.Meta.password)
	if authorised != true || err != nil {
		c.UI.Error(
			"Sorry, we have been unable to authenticate your credentials\n",
		)
		return 1
	}

	var vdcb api.VdcBuildData
	vdcb, err = papi.GetVdcBuild(c.Meta.buildID)

	if err != nil {
		c.UI.Error(
			"Sorry, we have been unable to retrieve your VDCs\n",
		)
		if err != nil && len(err.Error()) > 0 {
			c.UI.Error(
				err.Error() + "\n",
			)
		}

		return 1
	}

	flags := tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', flags)

	fmt.Fprintln(w, "ID\tCreated At\tCreated By\tStatus\t")
	fmt.Fprintln(w, vdcb.Data.ID+"\t"+vdcb.Data.Attributes.CreatedAt+"\t"+vdcb.Data.Attributes.CreatedBy+"\t"+vdcb.Data.Attributes.State+"\t")

	w.Flush()

	return 0
}

// Help is called when ukc accounts -help | --help | -h
func (c *VdcBuildCommand) Help() string {
	helpText := `
Usage: ukc vdcbuild [options] [path]

  ` + c.Synopsis() + `

Required:
  -buildid       The BuildID you wish to check

Options:
  -email           Email to your UKCloud Portal Account

  -password        Password to your UKCloud Portal Acccount
`
	return strings.TrimSpace(helpText)
}

// Synopsis is used when listing all commands or in the help
func (c *VdcBuildCommand) Synopsis() string {
	return "Gives the status of a VDC Build"
}

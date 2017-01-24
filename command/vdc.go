package command

import (
	"fmt"
	"github.com/ukcloud/ukcloud-portal-api/api"
	"os"
	"strings"
	"text/tabwriter"
)

// VdcCommand is the command to list UKC Vdcs
type VdcCommand struct {
	Meta
}

// Run is will be executed when `ukc accounts` is called
func (c *VdcCommand) Run(args []string) int {

	var err error
	papi := new(api.API)
	args = c.Meta.process(args, false)

	// Get the flags for this)
	cmdName := "vdcs"
	cmdFlags := c.Meta.flagSet(cmdName)
	cmdFlags.StringVar(&c.Meta.email, "email", "", "email")
	cmdFlags.StringVar(&c.Meta.password, "password", "", "password")
	cmdFlags.IntVar(&c.Meta.accountID, "accountid", 0, "accountid")
	cmdFlags.IntVar(&c.Meta.vOrgID, "vorgid", 0, "vorgid")
	cmdFlags.BoolVar(&c.Meta.create, "create", false, "create")
	cmdFlags.BoolVar(&c.Meta.silent, "silent", true, "silent")
	cmdFlags.StringVar(&c.Meta.name, "name", "", "name")

	cmdFlags.Usage = func() { c.UI.Error(c.Help()) }
	err = cmdFlags.Parse(args)
	if c.Meta.accountID == 0 || c.Meta.vOrgID == 0 || err != nil {
		cmdFlags.Usage()
		return 1
	}

	if c.Meta.create == true && c.Meta.name == "" {
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

	if c.Meta.create == true {
		var loc, err = papi.CreateVdc(c.Meta.accountID, c.Meta.vOrgID, c.Meta.name)

		if err != nil && len(err.Error()) > 0 {
			c.UI.Error(
				err.Error() + "\n",
			)
		}

		if c.Meta.silent == true {
			c.UI.Info(
				"VDC now building. You may check the progress here:" + "\n",
			)
			c.UI.Info(
				loc + "\n",
			)
			return 0
		}

		// @TODO - add non-silent flag to keep polling build status
	}

	var vdcs api.VdcArray
	vdcs, err = papi.GetVdc(c.Meta.accountID, c.Meta.vOrgID)

	if err != nil || len(vdcs.Data) <= 0 {
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

	fmt.Fprintln(w, "ID\tVDC Name\t")
	for _, vdc := range vdcs.Data {
		fmt.Fprintln(w, vdc.ID+"\t"+vdc.Attributes.Name+"\t")
	}
	w.Flush()

	return 0
}

// Help is called when ukc accounts -help | --help | -h
func (c *VdcCommand) Help() string {
	helpText := `
Usage: ukc vdcs [options] [path]

  ` + c.Synopsis() + `

Required:
  -accountid     Your UKCloud AccountID
  -vorgid        Your UKCloud VorgID

Options:
  -email         Email to your UKCloud Portal Account

  -password      Password to your UKCloud Portal Acccount

  -create        Create a VDC

  -name          Required if -create provided. Name of the VDC to create

  -silent        Creates the request for a VDC but does not wait for completion
`
	return strings.TrimSpace(helpText)
}

// Synopsis is used when listing all commands or in the help
func (c *VdcCommand) Synopsis() string {
	return "List the VDCs in your UKCloud Portal"
}

package command

import (
	"encoding/json"
	"fmt"
	"github.com/UKCloud/ukcloud-portal-api/api"
	"github.com/UKCloud/ukcloud-portal-cli/ukc"
	"os"
	"strings"
	"text/tabwriter"
)

// VorgsCommand is the command to list UKC Vorgs
type VorgsCommand struct {
	Meta
}

// Run is will be executed when `ukc accounts` is called
func (c *VorgsCommand) Run(args []string) int {

	var err error
	papi := new(api.API)
	args = c.Meta.process(args, false)

	// Get the flags for this)
	cmdName := "vorgs"
	cmdFlags := c.Meta.flagSet(cmdName)
	cmdFlags.StringVar(&c.Meta.email, "email", "", "email")
	cmdFlags.StringVar(&c.Meta.password, "password", "", "password")
	cmdFlags.IntVar(&c.Meta.accountID, "accountid", 0, "accountid")
	cmdFlags.BoolVar(&c.Meta.Json, "json", false, "json")
	cmdFlags.Usage = func() { c.UI.Error(c.Help()) }
	err = cmdFlags.Parse(args)

	var authorised bool
	authorised, err = papi.GetAuth(c.Meta.email, c.Meta.password)
	if authorised != true || err != nil {
		c.UI.Error(
			"Sorry, we have been unable to authenticate your credentials\n",
		)
		return 1
	}

	if c.Meta.accountID == 0 || err != nil {
		c.UI.Error(
			"Possible accounts are:",
		)

		c.Meta.accountID = ukc.PromptAccount(papi)

	}

	var vorgs api.VorgsArray
	vorgs, err = papi.GetVorgs(c.Meta.accountID)

	if err != nil || len(vorgs.Data) <= 0 {
		c.UI.Error(
			"Sorry, we have been unable to retrieve your Vorgs\n",
		)
		if err != nil && len(err.Error()) > 0 {
			c.UI.Error(
				err.Error() + "\n",
			)
		}

		return 1
	}

	if c.Meta.Json {
		output, _ := json.Marshal(vorgs.Data)
		fmt.Println(string(output))
		return 0
	}

	flags := tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', flags)

	fmt.Fprintln(w, "ID\tVorgs Name\t")
	for _, vorgs := range vorgs.Data {
		fmt.Fprintln(w, vorgs.ID+"\t"+vorgs.Attributes.Name+"\t")
	}
	w.Flush()

	return 0
}

// Help is called when ukc accounts -help | --help | -h
func (c *VorgsCommand) Help() string {
	helpText := `
Usage: ukc vorgs [options] [path]

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
func (c *VorgsCommand) Synopsis() string {
	return "List the Vorgs in your UKCloud Portal"
}

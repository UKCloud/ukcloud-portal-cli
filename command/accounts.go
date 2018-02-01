package command

import (
	"encoding/json"
	"fmt"
	"github.com/UKCloud/ukcloud-portal-api/api"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

// AccountsCommand is the command to list UKC Accounts
type AccountsCommand struct {
	Meta
}

// Run is will be executed when `ukc accounts` is called
func (c *AccountsCommand) Run(args []string) int {

	var err error
	papi := new(api.API)
	args = c.Meta.process(args, false)

	// Get the flags for this)
	cmdName := "accounts"
	cmdFlags := c.Meta.flagSet(cmdName)
	cmdFlags.StringVar(&c.Meta.email, "email", "", "email")
	cmdFlags.StringVar(&c.Meta.password, "password", "", "password")
	cmdFlags.BoolVar(&c.Meta.Json, "json", false, "json")

	cmdFlags.Usage = func() { c.UI.Error(c.Help()) }
	if err = cmdFlags.Parse(args); err != nil {
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

	var accounts []api.Accounts
	accounts, err = papi.GetAccounts()
	if err != nil {
		c.UI.Error(
			"Sorry, there was an error fetching your accounts\n",
		)
		return 1
	}

	if c.Meta.Json {
		output, _ := json.Marshal(accounts)
		fmt.Println(string(output))
		return 0
	}

	flags := tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', flags)

	fmt.Fprintln(w, "ID\tAccount Name\t")
	for _, account := range accounts {
		fmt.Fprintln(w, strconv.Itoa(account.ID)+"\t"+account.Name+"\t")
	}
	w.Flush()

	return 0
}

// Help is called when ukc accounts -help | --help | -h
func (c *AccountsCommand) Help() string {
	helpText := `
Usage: ukc accounts [options] [path]

  ` + c.Synopsis() + `

Options:
  -email           Email to your UKCloud Portal Account

  -password        Password to your UKCloud Portal Acccount
`
	return strings.TrimSpace(helpText)
}

// Synopsis is used when listing all commands or in the help
func (c *AccountsCommand) Synopsis() string {
	return "List the Accounts in your UKCloud Portal"
}

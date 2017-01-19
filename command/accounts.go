package command

import (
	"strconv"
	"strings"
)

// AccountsCommand is the command to list UKC Accounts
type AccountsCommand struct {
	Meta
}

// Accounts holds the JSON response
type Accounts struct {
	Name string
	ID   int
}

// Run is will be executed when `ukc accounts` is called
func (c *AccountsCommand) Run(args []string) int {

	args = c.Meta.process(args, false)

	cmdName := "accounts"

	cmdFlags := c.Meta.flagSet(cmdName)

	cmdFlags.StringVar(&c.Meta.email, "email", "", "email")
	cmdFlags.StringVar(&c.Meta.password, "password", "", "password")

	cmdFlags.Usage = func() { c.UI.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if auth(c.Meta.email, c.Meta.password) != 0 {
		c.UI.Error(
			"Sorry, we have been unable to authenticate your credentials\n",
		)
		return 1
	}

	var accounts []Accounts

	getJSON("https://portal.skyscapecloud.com/api/accounts.json", &accounts)

	c.UI.Output("ID\t| Account Name")
	c.UI.Output("----------------------------------------------------------")
	for _, account := range accounts {
		c.UI.Output(strconv.Itoa(account.ID) + "\t| " + account.Name)
	}

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

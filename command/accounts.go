package command

import (
	"fmt"
	"strings"
)

type AccountsCommand struct {
	Meta
}

type Accounts struct {
	Name string
	Id   int
}

func (c *AccountsCommand) Run(args []string) int {

	args = c.Meta.process(args, false)

	cmdName := "accounts"

	cmdFlags := c.Meta.flagSet(cmdName)

	cmdFlags.StringVar(&c.Meta.email, "email", "", "email")
	cmdFlags.StringVar(&c.Meta.password, "password", "", "password")

	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if auth(c.Meta.email, c.Meta.password) != 0 {
		c.Ui.Error(
			"Sorry, we have been unable to authenticate your credentials\n",
		)
		return 1
	}

	var accounts []Accounts

	getJson("https://portal.skyscapecloud.com/api/accounts.json", &accounts)
	fmt.Println(accounts)

	return 0
}

func (c *AccountsCommand) Help() string {
	helpText := `
Usage: ukc accounts [options] [path]

  Lists the UKCloud Accounts

Options:
  -email           Email to your UKCloud Portal Account

  -password        Password to your UKCloud Portal Acccount
`
	return strings.TrimSpace(helpText)
}

func (c *AccountsCommand) Synopsis() string {
	return "List the Accounts in your UKCloud Portal"
}

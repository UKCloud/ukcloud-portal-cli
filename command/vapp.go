package command

import (
	govcd "github.com/ukcloud/govcloudair"

	"bufio"
	"fmt"
	"github.com/ukcloud/ukcloud-portal-api/api"
	"github.com/ukcloud/ukcloud-portal-cli/ukc"
	"log"
	"os"
	"strconv"
	"strings"
)

type VCDClient struct {
	*govcd.VCDClient
	MaxRetryTimeout int
	InsecureFlag    bool
}

// VorgsCommand is the command to list UKC Vorgs
type VappCommand struct {
	Meta
}

// Run is will be executed when `ukc vapps` is called
func (c *VappCommand) Run(args []string) int {

	args = c.Meta.process(args, false)

	return 0
}

func (c *VappCommand) ChangeState(args []string, action string) int {
	var err error
	papi := new(api.API)
	args = c.Meta.process(args, false)

	// Get the flags for this)
	cmdName := "vapp"
	cmdFlags := c.Meta.flagSet(cmdName)
	cmdFlags.StringVar(&c.Meta.email, "email", "", "email")
	cmdFlags.StringVar(&c.Meta.username, "username", "", "username")
	cmdFlags.StringVar(&c.Meta.password, "password", "", "password")
	cmdFlags.StringVar(&c.Meta.vappName, "vappname", "", "vappname")
	cmdFlags.IntVar(&c.Meta.accountID, "accountid", 0, "accountid")
	cmdFlags.IntVar(&c.Meta.vOrgID, "vorgid", 0, "vorgid")
	cmdFlags.StringVar(&c.Meta.vOrg, "vorg", "", "vorg")
	cmdFlags.StringVar(&c.Meta.vdc, "vdc", "", "vdc")
	cmdFlags.BoolVar(&c.Meta.Json, "json", false, "json")
	cmdFlags.Usage = func() { c.UI.Error(c.Help()) }
	err = cmdFlags.Parse(args)

	fmt.Println(args)

	var authorised bool
	authorised, err = papi.GetAuth(c.Meta.email, c.Meta.password)
	if authorised != true || err != nil {
		c.UI.Error(
			"Sorry, we have been unable to authenticate your credentials\n",
		)
		c.UI.Error(
			err.Error(),
		)
		return 1
	}

	if c.Meta.accountID == 0 {
		c.UI.Error(
			"Possible accounts are:",
		)

		c.Meta.accountID = ukc.PromptAccount(papi)
	}

	creds, err := papi.GetCredentials(c.Meta.accountID)

	if c.Meta.vOrg == "" {
		c.UI.Error(
			"Possible Vorgs are:",
		)

		c.Meta.username, c.Meta.vOrg = ukc.PromptAccountVorg(creds)

	}

	if c.Meta.username == "" {
		c.UI.Error(
			"Enter username",
		)

		reader := bufio.NewReader(os.Stdin)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		c.Meta.username = strings.ToLower(strings.TrimSpace(response))
	}

	if c.Meta.vdc == "" {
		c.UI.Error(
			"Possible VDCs are:",
		)

		vOrgIdSplit := strings.Split(c.Meta.vOrg, "-")

		c.Meta.vOrgID, err = strconv.Atoi(vOrgIdSplit[2])

		if err != nil {
			c.UI.Error(
				"Error extracting vOrgId",
			)

			return 1
		}

		c.Meta.vdc = ukc.PromptVdc(papi, c.Meta.accountID, c.Meta.vOrgID)

		fmt.Println("VDC used: " + c.Meta.vdc)
	}

	fmt.Println(c.Meta.vOrg)

	if c.Meta.vappName == "" || err != nil {
		c.UI.Error(
			"Please enter vApp name:",
		)

		reader := bufio.NewReader(os.Stdin)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		c.Meta.vappName = strings.ToLower(strings.TrimSpace(response))
	}

	config := ukc.Config{
		User:            c.Meta.username,
		Password:        c.Meta.password,
		Org:             c.Meta.vOrg,
		Href:            "https://api.vcd.pod0000b.sys00005.portal.skyscapecloud.com/api",
		VDC:             c.Meta.vdc,
		MaxRetryTimeout: 600,
		InsecureFlag:    true,
	}

	vcdclient, err := config.Client()

	if err != nil {
		c.UI.Error("Error getting client" + err.Error())
		return 1
	}

	vapp, err := vcdclient.OrgVdc.FindVAppByName(c.Meta.vappName)

	if err != nil {
		c.UI.Error("Error finding vapp")
		return 1
	}

	status, err := vapp.GetStatus()

	fmt.Println("Current vApp Status: " + status)

	if err != nil {
		c.UI.Error("Error getting vapp status")
		return 1
	}

	var task govcd.Task
	if action == "shutdown" {
		task, err = vapp.PowerOff()
	}

	if action == "boot" {
		task, err = vapp.PowerOn()
	}

	if action == "reboot" {
		task, err = vapp.Reboot()
	}

	if err != nil {
		c.UI.Error("Error powering off")
		return 1
	}

	if task.Task != nil {
		err = task.WaitTaskCompletion()
		if err != nil {
			c.UI.Error("Error completing tasks")
			return 1
		}
	}

	fmt.Println(task.Task.Status)
	status, err = vapp.GetStatus()

	fmt.Println("Current vApp Status: " + status)

	return 0
}

// Help is called when ukc accounts -help | --help | -h
func (c *VappCommand) Help() string {
	helpText := `
Usage: ukc vapps [shutdown|poweron|reboot]

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
func (c *VappCommand) Synopsis() string {
	return "Control the Vapps in your UKCloud Account"
}

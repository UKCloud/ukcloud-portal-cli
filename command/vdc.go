package command

import (
	"encoding/json"
	"fmt"
	"github.com/ukcloud/ukcloud-portal-api/api"
	"github.com/ukcloud/ukcloud-portal-cli/ukc"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
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
	cmdFlags.BoolVar(&c.Meta.silent, "silent", false, "silent")
	cmdFlags.BoolVar(&c.Meta.Json, "json", false, "json")
	cmdFlags.StringVar(&c.Meta.name, "name", "", "name")
	cmdFlags.Usage = func() { c.UI.Error(c.Help()) }
	err = cmdFlags.Parse(args)

	if c.Meta.create == true && c.Meta.name == "" {
		c.Meta.name = ukc.PromptVdcName(papi)
	}

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

	if c.Meta.vOrgID == 0 {
		c.UI.Error(
			"Possible Vorgs are:",
		)

		c.Meta.vOrgID = ukc.PromptVorg(papi, c.Meta.accountID)
	}

	if c.Meta.create == true {
		var loc, err = papi.CreateVdc(c.Meta.accountID, c.Meta.vOrgID, c.Meta.name)

		if err != nil && len(err.Error()) > 0 {
			c.UI.Error(
				err.Error() + "\n",
			)
		}

		if c.Meta.Json == true {
			buildInfo := map[string]interface{}{
				"BuildId": extractBuildID(loc),
			}
			output, _ := json.Marshal(buildInfo)
			fmt.Println(string(output))
			return 0

		}

		c.UI.Info(
			"VDC now building. Your buildID is: " + strconv.Itoa(extractBuildID(loc)) + "\n",
		)

		if c.Meta.silent == true {
			return 0
		}

		buildID := extractBuildID(loc)
		var vdcb api.VdcBuildData

		status := ""
		prevStatus := ""
		for status != "completed" {

			vdcb, err = papi.GetVdcBuild(buildID)

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

			status = vdcb.Data.Attributes.State
			if status == "completed" {
				break
			}

			if status != prevStatus {

				fmt.Println("Current Status: " + status)
				prevStatus = status
			}

			amt := time.Duration(3000)
			time.Sleep(time.Millisecond * amt)
			fmt.Println("... checking status of build.")

		}
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

	if c.Meta.Json {
		output, _ := json.Marshal(vdcs.Data)
		fmt.Println(string(output))
		return 0
	}

	flags := tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', flags)

	fmt.Fprintln(w, "ID\tVDC Name\tBuild Status\t")
	for _, vdc := range vdcs.Data {
		fmt.Fprintln(w, vdc.ID+"\t"+strings.Trim(vdc.Attributes.Name, "\n")+"\tCompleted\t")
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

func extractBuildID(BuildURL string) int {
	n := strings.SplitN(BuildURL, "/", 4)[3]
	buildID, _ := strconv.ParseInt(n, 10, 64)
	return int(buildID)
}
